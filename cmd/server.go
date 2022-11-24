/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/server"
	"github.com/mxplusb/pleiades/pkg/server/eventing"
	"github.com/mxplusb/pleiades/pkg/server/kvstore"
	"github.com/mxplusb/pleiades/pkg/server/raft"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/mxplusb/pleiades/pkg/server/shard"
	"github.com/mxplusb/pleiades/pkg/server/transactions"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run a local instance of pleiades",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	defaultDataBasePath := ""
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}
		defaultDataBasePath = filepath.Join(dir, "Library", "pleiades")
	} else {
		defaultDataBasePath = configuration.DefaultBaseDataPath
	}
	config.Set("server.datastore.basePath", defaultDataBasePath)

	serverCmd.PersistentFlags().Uint64("deployment-id", 0, "identifier for this deployment")
	config.BindPFlag("server.host.deploymentId", serverCmd.PersistentFlags().Lookup("deployment-id"))

	serverCmd.PersistentFlags().String("grpc-addr", "0.0.0.0:8080", "grpc listener address")
	config.BindPFlag("server.host.grpcListenAddress", serverCmd.PersistentFlags().Lookup("grpc-addr"))

	serverCmd.PersistentFlags().String("raft-addr", "0.0.0.0:8081", "raft listener address")
	config.BindPFlag("server.host.listenAddress", serverCmd.PersistentFlags().Lookup("raft-addr"))

	serverCmd.LocalFlags().Bool("notify-commit", false, "enable raft commit notifications")
	config.BindPFlag("server.host.notifyCommit", serverCmd.LocalFlags().Lookup("notify-commit"))

	serverCmd.PersistentFlags().Uint64("round-trip", 0, "average round trip time, plus processing, in milliseconds to other hosts in the data centre")
	config.BindPFlag("server.host.rtt", serverCmd.PersistentFlags().Lookup("round-trip"))

	serverCmd.PersistentFlags().String("base-path", config.GetString("server.datastore.basePath"), "base directory for data")
	config.BindPFlag("server.datastore.basePath", serverCmd.PersistentFlags().Lookup("base-path"))

	serverCmd.PersistentFlags().String("log-dir", "logs", "directory for raft logs, relative to base-path")
	config.BindPFlag("server.datastore.logDir", serverCmd.PersistentFlags().Lookup("log-dir"))

	serverCmd.PersistentFlags().String("data-dir", "data", "directory for data, relative to base-path")
	config.BindPFlag("server.datastore.dataDir", serverCmd.PersistentFlags().Lookup("data-dir"))

	serverCmd.PersistentFlags().String("continent", "north-america", "the continent this server is located in")
	config.BindPFlag("server.gossip.continent", serverCmd.PersistentFlags().Lookup("continent"))

	serverCmd.PersistentFlags().String("region", "us-central1", "the region this server is located in")
	config.BindPFlag("server.gossip.region", serverCmd.PersistentFlags().Lookup("region"))

	serverCmd.PersistentFlags().String("zone", "us-central1-a", "the zone this server is located in")
	config.BindPFlag("server.gossip.zone", serverCmd.PersistentFlags().Lookup("zone"))

	serverCmd.PersistentFlags().Int("gossip-port", 8082, "the port the gossip server runs on")
	config.BindPFlag("server.gossip.port", serverCmd.PersistentFlags().Lookup("gossip-port"))
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	logger := setupLogger(cmd, args)

	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	serverutils.SetRootLogger(logger)

	var serverConfig configuration.Configuration
	err = config.Unmarshal(&serverConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't unmarshal configuration")
	}

	logger.Info().Interface("config", serverConfig).Msg("runtime configuration")

	// make the directories
	logDir := filepath.Join(config.GetString("server.datastore.basePath"), serverConfig.Server.Datastore.LogDir)
	err = os.MkdirAll(logDir, 0750)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal().Err(err).Msg("can't create log directory")
		}
	}
	config.Set("server.datastore.logDir", logDir)

	dataDir := filepath.Join(config.GetString("server.datastore.basePath"), serverConfig.Server.Datastore.DataDir)
	err = os.MkdirAll(dataDir, 0750)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal().Err(err).Msg("can't create log directory")
		}
	}
	config.Set("server.datastore.dataDir", dataDir)

	nhc := dconfig.NodeHostConfig{
		DeploymentID:   serverConfig.Server.Host.DeploymentId,
		WALDir:         logDir,
		NodeHostDir:    dataDir,
		RTTMillisecond: serverConfig.Server.Host.Rtt,
		RaftAddress:    serverConfig.Server.Host.ListenAddress,
		EnableMetrics:  true,
		NotifyCommit:   serverConfig.Server.Host.NotifyCommit,
	}

	if serverConfig.Server.Host.MutualTLS {
		nhc.MutualTLS = serverConfig.Server.Host.MutualTLS
		nhc.CAFile = serverConfig.Server.Host.CaFile
		nhc.CertFile = serverConfig.Server.Host.CertFile
		nhc.KeyFile = serverConfig.Server.Host.KeyFile
	}

	app := fx.New(
		fx.Provide(func() *viper.Viper {
			return config
		}),
		fx.Provide(func() dconfig.NodeHostConfig {
			nhc := dconfig.NodeHostConfig{
				DeploymentID:   serverConfig.Server.Host.DeploymentId,
				WALDir:         logDir,
				NodeHostDir:    dataDir,
				RTTMillisecond: serverConfig.Server.Host.Rtt,
				RaftAddress:    serverConfig.Server.Host.ListenAddress,
				EnableMetrics:  true,
				NotifyCommit:   serverConfig.Server.Host.NotifyCommit,
			}

			if serverConfig.Server.Host.MutualTLS {
				nhc.MutualTLS = serverConfig.Server.Host.MutualTLS
				nhc.CAFile = serverConfig.Server.Host.CaFile
				nhc.CertFile = serverConfig.Server.Host.CertFile
				nhc.KeyFile = serverConfig.Server.Host.KeyFile
			}
			return nhc
		}),
		fx.Provide(func() zerolog.Logger { // the generalized logger
			return logger
		}),
		fx.WithLogger(func() fxevent.Logger { // this provides the fx logger, not the general logger
			return zerologAdapter{logger}
		}),
		fx.Provide(messaging.NewEmbeddedWorkflowServer),
		fx.Provide(eventing.NewServer),
		fx.Provide(server.NewHttpServeMux),
		fx.Provide(server.NewNodeHost),
		fx.Provide(server.NewHttpServer),
		fx.Provide(eventing.NewPubSubClient),
		fx.Provide(eventing.NewStreamClient),
		fx.Provide(eventing.NewLifecycleManager),
		fx.Provide(kvstore.NewBboltStoreManager),
		fx.Provide(server.AsRoute(kvstore.NewKvstoreBboltConnectAdapter)),
		fx.Provide(server.AsRoute(kvstore.NewKvstoreTransactionConnectAdapter)),
		fx.Provide(raft.NewHost),
		fx.Provide(server.AsRoute(raft.NewRaftHostConnectAdapter)),
		fx.Provide(server.AsRoute(shard.NewRaftShardConnectAdapter)),
		fx.Provide(shard.NewManager),
		fx.Provide(transactions.NewManager),
		fx.Invoke(eventing.NewLifecycleManager),
		fx.Invoke(server.NewHttpServer),
	)

	if err := app.Start(ctx); err != nil {
		logger.Fatal().Err(err).Msg("can't start services")
	}

	<-app.Done()

	ctx, cancel := context.WithTimeout(ctx, 3000*time.Millisecond)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		logger.Error().Err(err).Msg("can't safely stop system")
	}

	logger.Info().Msg("done")
}

type zerologAdapter struct {
	logger zerolog.Logger
}

// LogEvent logs the given event to the provided Zap logger.
func (zl zerologAdapter) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		zl.logger.Debug().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on start hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on start hook executed")
		} else {
			zl.logger.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on start hook executed")
		}
	case *fxevent.OnStopExecuting:
		zl.logger.Debug().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on stop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on stop hook failed")
		} else {
			zl.logger.Debug().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on stop hook executed")
		}
	case *fxevent.Supplied:
		zl.logger.Debug().Str("type", e.TypeName).Str("module", e.ModuleName).Msg("supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			zl.logger.Debug().Str("constructor", e.ConstructorName).Str("module", e.ModuleName).Str("type", rtype).Msg("constructor provided")
		}
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Str("module", e.ModuleName).Msg("error encountered while applying options")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			zl.logger.Debug().Str("module", e.ModuleName).Str("type", rtype).Msg("replaced")
		}
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Str("module", e.ModuleName).Msg("error while replacing")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			zl.logger.Debug().Str("module", e.ModuleName).Str("type", rtype).Str("decorated", e.DecoratorName).Msg("decorated")
		}
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Str("module", e.ModuleName).Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		zl.logger.Debug().Str("module", e.ModuleName).Str("function", e.FunctionName).Msg("invoking")
	case *fxevent.Invoked:
		if e.Err != nil {
			zl.logger.Debug().Str("module", e.ModuleName).Err(e.Err).Str("stack", e.Trace).Str("function", e.FunctionName).Msg("invoke failed")
		}
		zl.logger.Debug().Str("module", e.ModuleName).Str("function", e.FunctionName).Msg("invoked")
	case *fxevent.Stopping:
		zl.logger.Debug().Str("signal", strings.ToUpper(e.Signal.String())).Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		zl.logger.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Msg("start failed")
		} else {
			zl.logger.Info().Msg("started the server")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			zl.logger.Debug().Str("function", e.ConstructorName).Msg("initialized logger")
		}
	}
}
