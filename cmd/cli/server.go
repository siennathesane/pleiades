/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/server"
	"github.com/mxplusb/pleiades/pkg/utils/serverutils"
	"github.com/posener/complete"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var (
	_ cli.Command = (*ServerCommand)(nil)
)

func init() {
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
}

type ServerCommand struct {
	*BaseCommand

	flagDeploymentId            uint64
	flagBasePath                string
	flagServiceDiscoveryAddress string
	flagListenAddr              string
	flagHttpPort                int
	flagFabricPort              int
	flagConstellationPort       int
	flagNotifyCommit            bool
	flagRoundTrip               uint64
	flagDryRun                  bool
}

func (s *ServerCommand) Help() string {
	helpText := `Runs an instance of the Pleiades Database Constellation.

An instance of Pleiades involves multiple internal processes and several open 
ports. It's recommended to run Pleiades on a large system to get proper 
performance. This command is configured to use recommended defaults, and very 
little of the defaults should be changed. 

` + s.Flags().Help()

	return helpText
}

func (s *ServerCommand) Flags() *FlagSets {
	set := s.flagSet(FlagSetTls | FlagSetLogging)

	f := set.NewFlagSet("Server Options")

	f.Uint64Var(&Uint64Var{
		Name:              "deployment-id",
		Usage:             `Set the deployment ID of this node. It's how nodes determine if they are a part of the same deployment.`,
		Default:           1,
		EnvVar:            EnvPleiadesDeploymentId,
		Target:            &s.flagDeploymentId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.deployment-id",
	})

	f.StringVar(&StringVar{
		Name: "base-path",
		Usage: `The base directory for all of the node data. This can be changed, but changing it between 
runs will reset the node configuration, essentially wiping it clean.`,
		Default:           config.GetString("server.datastore.basePath"),
		Hidden:            false,
		EnvVar:            EnvPleiadesDataDir,
		Target:            &s.flagBasePath,
		Completion:        complete.PredictDirs("*"),
		ConfigurationPath: "server.datastore.basePath",
	})

	hname, err := os.Hostname()
	if err != nil {
		// I have no idea why this would cause an error but meh
		s.UI.Error(err.Error())
		return nil
	}
	f.StringVar(&StringVar{
		Name: "fabric-hostname",
		Usage: `Address the fabric subsystems will use to identify this node to other fabric nodes. This 
cannot be changed between runs, so it's best to set it to the externally addressable hostname.`,
		Default:           strings.ToLower(hname),
		Hidden:            false,
		EnvVar:            EnvPleiadesFabricAddr,
		Target:            &s.flagServiceDiscoveryAddress,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.fabricHostname",
	})

	f.StringVar(&StringVar{
		Name:              "listen-addr",
		Usage:             "The IP address to listen on.",
		Default:           "0.0.0.0",
		Hidden:            false,
		EnvVar:            EnvPleidesListenAddr,
		Target:            &s.flagListenAddr,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.listenAddr",
	})

	f.IntVar(&IntVar{
		Name:              "http-port",
		Usage:             "The HTTP port to listen on.",
		Default:           8080,
		Hidden:            false,
		EnvVar:            EnvPleiadesHttpPort,
		Target:            &s.flagHttpPort,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.httpListenPort",
	})

	f.IntVar(&IntVar{
		Name:              "fabric-port",
		Usage:             "The fabric port to listen on.",
		Default:           8081,
		Hidden:            false,
		EnvVar:            EnvPleiadesFabricPort,
		Target:            &s.flagFabricPort,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.fabricListenPort",
	})

	f.IntVar(&IntVar{
		Name:              "constellation-port",
		Usage:             "The constellation port to listen on.",
		Default:           8082,
		Hidden:            false,
		EnvVar:            EnvPleiadesConstellationPort,
		Target:            &s.flagConstellationPort,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.constellationListenPort",
	})

	f.BoolVar(&BoolVar{
		Name:              "notify-commit",
		Usage:             "Enable commit notifications. This is an alpha feature and is considered unstable.",
		Default:           false,
		Hidden:            false,
		EnvVar:            EnvPleiadesNotifyCommit,
		Target:            &s.flagNotifyCommit,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.notifyCommit",
	})

	f.Uint64Var(&Uint64Var{
		Name: "round-trip",
		Usage: `The length of time it takes to process fabric messages to the nearest nodes in the data 
centre, in milliseconds.`,
		Default:           10,
		Hidden:            false,
		EnvVar:            EnvPleiadesRoundTrip,
		Target:            &s.flagRoundTrip,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.host.rtt",
	})

	f.BoolVar(&BoolVar{
		Name:              "dry-run",
		Usage:             `Dry run an instance of Pleiades. This will print out all of the configurations but not boot 
the server.`,
		Default:           false,
		Hidden:            true,
		Target:            &s.flagDryRun,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "server.dry-run",
	})

	return set
}

func (s *ServerCommand) Synopsis() string {
	return `Run an instance of Pleiades`
}

func (s *ServerCommand) Run(args []string) int {
	f := s.Flags()

	if err := f.Parse(args); err != nil {
		s.UI.Error(err.Error())
		return exitCodeFailureToParseArgs
	}

	ctx := context.Background()
	logger := setupLogger()

	serverutils.SetRootLogger(logger)

	// make the directories
	logDir := filepath.Join(config.GetString("server.datastore.basePath"), configuration.DefaultLogDir)
	err := os.MkdirAll(logDir, 0750)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal().Err(err).Msg("can't create log directory")
		}
	}
	config.Set("server.datastore.logDir", logDir)

	dataDir := filepath.Join(config.GetString("server.datastore.basePath"), configuration.DefaultDataDir)
	err = os.MkdirAll(dataDir, 0750)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal().Err(err).Msg("can't create data directory")
		}
	}
	config.Set("server.datastore.dataDir", dataDir)

	raftAddr := fmt.Sprintf("%s:%d", config.GetString("server.host.fabricHostname"), config.GetUint("server.host.fabricListenPort"))
	config.Set("server.host.fabricAddr", raftAddr)

	nodeAddr := fmt.Sprintf("%s:%d", config.GetString("server.host.listenAddr"), config.GetUint("server.host.fabricListenPort"))
	config.Set("server.host.nodeAddr", nodeAddr)

	if s.flagDryRun {
		OutputData(s.UI, config.AllSettings())
		return exitCodeGood
	}

	logger.Debug().Interface("config", config.AllSettings()).Msg("runtime configuration")

	app := fx.New(
		fx.Provide(func() *viper.Viper {
			return config
		}),
		fx.Provide(func() dconfig.NodeHostConfig {
			nhc := dconfig.NodeHostConfig{
				DeploymentID:   config.GetUint64("server.host.deployment-id"),
				WALDir:         logDir,
				NodeHostDir:    dataDir,
				RTTMillisecond: config.GetUint64("server.host.rtt"),
				ListenAddress:  nodeAddr,
				RaftAddress:    raftAddr,
				EnableMetrics:  true,
				NotifyCommit:   config.GetBool("server.host.notifyCommit"),
			}

			certFile := config.GetString("tls.cert-file")
			keyFile := config.GetString("tls.key-file")
			caFile := config.GetString("tls.ca-cert-file")

			if certFile != "" && keyFile != "" && caFile != "" {
				nhc.MutualTLS = true
				nhc.CAFile = caFile
				nhc.CertFile = certFile
				nhc.KeyFile = keyFile
			}
			return nhc
		}),
		fx.Provide(func() zerolog.Logger { // the generalized logger
			return logger
		}),
		fx.WithLogger(func() fxevent.Logger { // this provides the fx logger, not the general logger
			return zerologAdapter{logger}
		}),
		server.ServerModule,
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

	return exitCodeGood
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
			zl.logger.Debug().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("on start hook executed")
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
			zl.logger.Debug().Msg("started the server")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			zl.logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			zl.logger.Debug().Str("function", e.ConstructorName).Msg("initialized logger")
		}
	}
}
