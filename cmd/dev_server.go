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
	"sync"
	"time"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/server"
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

var reset = true

func init() {
	devCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().Uint64("deployment-id", 1, "identifier for this deployment")
	config.BindPFlag("server.host.deploymentId", devCmd.PersistentFlags().Lookup("deployment-id"))

	serverCmd.PersistentFlags().String("grpc-addr", "0.0.0.0:5050", "grpc listener address")
	config.BindPFlag("server.host.grpcListenAddress", devCmd.PersistentFlags().Lookup("grpc-addr"))

	serverCmd.PersistentFlags().String("raft-addr", "0.0.0.0:5051", "raft listener address")
	config.BindPFlag("server.host.listenAddress", devCmd.PersistentFlags().Lookup("raft-addr"))

	serverCmd.PersistentFlags().Bool("notify-commit", false, "enable raft commit notifications")
	config.BindPFlag("server.host.notifyCommit", devCmd.PersistentFlags().Lookup("notify-commit"))

	serverCmd.PersistentFlags().Uint64("round-trip", 1, "average round trip time, plus processing, in milliseconds to other hosts in the data centre")
	config.BindPFlag("server.host.rtt", devCmd.PersistentFlags().Lookup("round-trip"))

	serverCmd.PersistentFlags().Bool("reset", false, "clean reset the dev server at init")
	config.BindPFlag("server.reset", devCmd.PersistentFlags().Lookup("reset"))
}

func startServer(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	var serverConfig configuration.ServerConfig
	err = config.Unmarshal(&serverConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't unmarshal configuration")
	}

	nhc := dconfig.NodeHostConfig{
		DeploymentID:   serverConfig.Host.DeploymentId,
		WALDir:         serverConfig.Datastore.LogDir,
		NodeHostDir:    serverConfig.Datastore.DataDir,
		RTTMillisecond: serverConfig.Host.Rtt,
		RaftAddress:    serverConfig.Host.ListenAddress,
		EnableMetrics:  true,
		NotifyCommit:   serverConfig.Host.NotifyCommit,
	}

	if serverConfig.Host.MutualTLS {
		nhc.MutualTLS = serverConfig.Host.MutualTLS
		nhc.CAFile = serverConfig.Host.CaFile
		nhc.CertFile = serverConfig.Host.CertFile
		nhc.KeyFile = serverConfig.Host.KeyFile
	}

	if config.GetBool("reset") {
		err := os.RemoveAll(serverConfig.Datastore.LogDir)
		if err != nil {
			logger.Fatal().Err(err).Str("dir", serverConfig.Datastore.LogDir).Msg("can't remove directory")
		}
		err = os.RemoveAll(serverConfig.Datastore.DataDir)
		if err != nil {
			logger.Fatal().Err(err).Str("dir", serverConfig.Datastore.DataDir).Msg("can't remove directory")
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
			err = s.GetRaftShardManager().NewShard(i, i*257, server.BBoltStateMachineType, 300*time.Millisecond)
		}()
		utils.Wait(100 * time.Millisecond)
	}
	wg.Wait()

	logger.Debug().Msg("state machines finished, starting server")

	http.ListenAndServe(
		config.GetString("server.host.listenAddr"),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

	s.Stop()
}
