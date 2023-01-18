/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"context"
	"net/http"
	"net/url"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/go-wordwrap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var raftAddReplicaHelp = `add a replica to a shard.

while a replica can be added to the same host which has the primary shard, it's recommended this be a different host. the replica id only matters for uniqueness, but has no other effect. state machine types are one-time choices handled by the primary shard, and you cannot change the type of a state machine after it's been created. the current supported state machine types are:

0. unspecified
1. test
2. kv

currently, unspecified will return nothing, and the test state machine will only store uint64s as a way to safely test functionality. the kv state machine will provide all the features of the kv store, so check the documentation for current supported features.`

// raftAddReplicaCmd represents the raftAddReplicaCmd command
var raftAddReplicaCmd = &cobra.Command{
	Use:   "add-replica",
	Short: "add a replica to a shard",
	Long: wordwrap.WrapString(raftAddShardHelp, 80),
	Run: runAddReplica,
}

var (
	raftAddReplicaFlags *raftv1.AddReplicaRequest = &raftv1.AddReplicaRequest{}
)

func init() {
	raftCmd.AddCommand(raftAddReplicaCmd)

	raftAddReplicaCmd.Flags().Uint64Var(&raftAddReplicaFlags.ShardId, "shard-id", 0, "id of the shard to create")
	raftAddReplicaCmd.Flags().Uint64Var(&raftAddReplicaFlags.ReplicaId, "replica-id", 0, "id of the replica to create")
	raftAddReplicaCmd.Flags().StringVar(&raftAddReplicaFlags.Hostname, "shard-host", "", "type of state machine")
	raftAddReplicaCmd.Flags().Int64Var(&raftAddReplicaFlags.Timeout, "timeout", 5000, "timeout length in milliseconds")
	raftAddReplicaCmd.MarkFlagRequired("shard-id")
	raftAddReplicaCmd.MarkFlagRequired("replica-id")
}

func runAddReplica(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	logger.Debug().Str("host", config.GetString("client.grpcAddr")).Msg("creating client")

	targetHost, err := url.Parse(config.GetString("client.grpcAddr"))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't parse remote host")
	}

	var host raftv1connect.ShardServiceClient
	if targetHost.Scheme != "https" {
		host = raftv1connect.NewShardServiceClient(newInsecureClient(),targetHost.String())
	} else {
		host = raftv1connect.NewShardServiceClient(http.DefaultClient, targetHost.String())
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(raftAddReplicaFlags.Timeout) * time.Millisecond)
	//defer cancel()

	logger.Debug().Interface("request", raftAddReplicaFlags).Msg("request payload")

	resp, err := host.AddReplica(context.Background(), connect.NewRequest[raftv1.AddReplicaRequest](raftAddReplicaFlags))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't add replica")
	}

	print(protojson.Format(resp.Msg))
}