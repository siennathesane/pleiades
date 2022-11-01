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
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/server"
	"github.com/mxplusb/pleiades/pkg/server/shard"
	"github.com/mxplusb/pleiades/pkg/utils"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run a development server",
	Long: `runs a development server.

it will boot with 256 predefined shards, configured in 
insecure mode, and will generally be buggy. it will run
the latest and greatest, which means it may or may not 
be usable for consuming applications. there may be unversioned
changes in this command which are not available as part of
the cloud offering. this command is unsupported beyond 
filing bugs against it the team may or may not get to

DO NOT USE THIS IN PRODUCTION`,
	Run: startServer,
}

func init() {
	devCmd.AddCommand(serverCmd)

	serverCmd.LocalFlags().Uint64("deployment-id", 1, "identifier for this deployment")
	config.BindPFlag("server.host.deploymentId", serverCmd.LocalFlags().Lookup("deployment-id"))

	serverCmd.LocalFlags().String("grpc-addr", "0.0.0.0:5050", "grpc listener address")
	config.BindPFlag("server.host.grpcListenAddress", serverCmd.LocalFlags().Lookup("grpc-addr"))

	serverCmd.LocalFlags().String("raft-addr", "0.0.0.0:5051", "raft listener address")
	config.BindPFlag("server.host.listenAddress", serverCmd.LocalFlags().Lookup("raft-addr"))

	serverCmd.LocalFlags().Bool("notify-commit", false, "enable raft commit notifications")
	config.BindPFlag("server.host.notifyCommit", serverCmd.LocalFlags().Lookup("notify-commit"))

	serverCmd.LocalFlags().Uint64("round-trip", 1, "average round trip time, plus processing, in milliseconds to other hosts in the data centre")
	config.BindPFlag("server.host.rtt", serverCmd.LocalFlags().Lookup("round-trip"))

	serverCmd.LocalFlags().String("base-path", config.GetString("server.datastore.basePath"), "base directory for data")
	config.BindPFlag("server.datastore.basePath", serverCmd.LocalFlags().Lookup("base-path"))

	serverCmd.LocalFlags().String("log-dir", "logs", "directory for raft logs, relative to base-path")
	config.BindPFlag("server.datastore.logDir", serverCmd.LocalFlags().Lookup("log-dir"))

	serverCmd.LocalFlags().String("data-dir", "data", "directory for data, relative to base-path")
	config.BindPFlag("server.datastore.dataDir", serverCmd.LocalFlags().Lookup("data-dir"))

	serverCmd.LocalFlags().Bool("reset", false, "clean reset the dev server at init")
	config.BindPFlag("server.reset", serverCmd.LocalFlags().Lookup("reset"))
}

func startServer(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	var serverConfig configuration.Configuration
	err = config.Unmarshal(&serverConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't unmarshal configuration")
	}

	nhc := dconfig.NodeHostConfig{
		DeploymentID:   serverConfig.Server.Host.DeploymentId,
		WALDir:         filepath.Join(config.GetString("server.datastore.basePath"), serverConfig.Server.Datastore.LogDir),
		NodeHostDir:    filepath.Join(config.GetString("server.datastore.basePath"), serverConfig.Server.Datastore.DataDir),
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

	if config.GetBool("reset") {
		err := os.RemoveAll(serverConfig.Server.Datastore.LogDir)
		if err != nil {
			logger.Fatal().Err(err).Str("dir", serverConfig.Server.Datastore.LogDir).Msg("can't remove directory")
		}
		err = os.RemoveAll(serverConfig.Server.Datastore.DataDir)
		if err != nil {
			logger.Fatal().Err(err).Str("dir", serverConfig.Server.Datastore.DataDir).Msg("can't remove directory")
		}
	}

	mux := http.NewServeMux()

	s, err := server.New(nhc, mux, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't create pleiades server")
	}

	var wg sync.WaitGroup
	// shardLimit+1
	for i := uint64(1); i < 257; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			err = s.GetRaftShardManager().NewShard(i, i*257, shard.BBoltStateMachineType, 300*time.Millisecond)
		}()
		utils.Wait(100 * time.Millisecond)
	}
	wg.Wait()

	logger.Debug().Msg("state machines finished, starting server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: config.GetString("server.host.grpcListenAddr"),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("cannot run server")
		}
	}()

	<-done
	s.Stop()
}
