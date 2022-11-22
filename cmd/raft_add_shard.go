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
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/go-wordwrap"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var raftAddShardHelp = `shards are added for the local node only.

this is to prevent communications and bootstrapping issues. the replica id only matters for uniqueness, but has no other effect. state machine types are one-time choices, and you cannot change the type of a state machine after it's been created. the current supported state machine types are:

0. unspecified
1. test
2. kv

currently, unspecified will return nothing, and the test state machine will only store uint64s as a way to safely test functionality. the kv state machine will provide all the features of the kv store, so check the documentation for current supported features.`

// raftAddShardCmd represents the raftAddShard command
var raftAddShardCmd = &cobra.Command{
	Use:   "add-shard",
	Short: "add a shard to a pleiades server",
	Long: wordwrap.WrapString(raftAddShardHelp, 80),
	Run: runAddShard,
}

var (
	raftAddShardFlags *raftv1.NewShardRequest = &raftv1.NewShardRequest{}
	raftAddShardSMType int32 = 0
)

func init() {
	raftCmd.AddCommand(raftAddShardCmd)

	raftAddShardCmd.Flags().Uint64Var(&raftAddShardFlags.ShardId, "shard-id", 0, "id of the shard to create")
	raftAddShardCmd.Flags().Uint64Var(&raftAddShardFlags.ReplicaId, "replica-id", 0, "id of the replica to create")
	raftAddShardCmd.Flags().Int32Var(&raftAddShardSMType, "state-machine-type", 0, "type of state machine")
	raftAddShardCmd.Flags().Int64Var(&raftAddShardFlags.Timeout, "timeout", 3000, "id of the shard to create")
	raftAddShardCmd.MarkFlagRequired("shard-id")
	raftAddShardCmd.MarkFlagRequired("replica-id")
	raftAddShardCmd.MarkFlagRequired("state-machine-type")
}

func runAddShard(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	logger.Debug().Str("host", config.GetString("server.client.grpcAddr")).Msg("creating client")

	host := raftv1connect.NewShardServiceClient(newInsecureClient(),config.GetString("server.client.grpcAddr"))

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(time.Duration(raftAddShardFlags.Timeout) * time.Millisecond))
	defer cancel()

	raftAddShardFlags.Type = raftv1.StateMachineType(raftAddShardSMType)

	logger.Debug().Interface("request", raftAddShardFlags).Msg("request payload")

	resp, err := host.NewShard(ctx, connect.NewRequest[raftv1.NewShardRequest](raftAddShardFlags))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't create new shard")
	}

	print(protojson.Format(resp.Msg))
}