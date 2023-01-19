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
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/api/raft/v1/raftv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-wordwrap"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*FabricAddShardCommand)(nil)
	_ cli.CommandAutocomplete = (*FabricAddShardCommand)(nil)
)

type FabricAddShardCommand struct {
	*BaseCommand

	flagShardId   uint64
	flagReplicaId uint64
	flagType      string
	flagHostname  string
	flagTimeout   int64
}

func (f *FabricAddShardCommand) Flags() *FlagSets {
	set := f.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging | FlagSetTimeout)
	fs := set.NewFlagSet("Fabric Options")

	fs.Uint64Var(&Uint64Var{
		Name: "shard-id",
		Usage: `The ID of the new shard. This is global to the node constellation as it increases the 
data fabric size.`,
		Target:            &f.flagShardId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-shard.shard-id",
	})

	fs.Uint64Var(&Uint64Var{
		Name:              "replica-id",
		Usage:             `The ID of the new replica. This is specific to each shard.`,
		Target:            &f.flagReplicaId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-shard.replica-id",
	})

	fs.StringVar(&StringVar{
		Name: "type",
		Usage: `The type of shard to create. See the greater help message for more information on the 
specific values.`,
		Target:            &f.flagType,
		Completion:        complete.PredictSet("kv"),
		ConfigurationPath: "fabric.add-shard.replica-id",
	})

	return set
}

func (f *FabricAddShardCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (f *FabricAddShardCommand) AutocompleteFlags() complete.Flags {
	return f.Flags().Completions()
}

// nb (sienna): use word wrap in the editor as this will format properly.
func (f *FabricAddShardCommand) Help() string {
	helpText := `Create a new shard for this node.

Currently, shards are added for the individual node called, and do not support multiple hosts being configured. This is to prevent bootstrapping issues common with the underlying Raft implementation. This behaviour is subject to change in later versions of Pleiades. The replica ID is unique to the shard, and so does not require ahead-of-time planning.

The data fabric is built on top of sharded, replicated, deterministic finite state machines (FSMs). With that, when creating new shards, operators must make a one-time choice of which state machine type to use. FSM types are one-time choices, and you cannot change the type of a state machine after it's been created. The list of supported FSMs can be found below.

Key Value FSM

The key value FSM is a generic key value store which can store large amounts of data. Data in the key value store is sharded evenly across all shards based on the key name. Pleiades will self-cluster and self-route all of the shards.

` + f.Flags().Help()

	return wordwrap.WrapString(helpText, 80)
}

func (f *FabricAddShardCommand) Run(args []string) int {
	fs := f.Flags()

	if err := fs.Parse(args); err != nil {
		f.UI.Error(err.Error())
		return exitCodeFailureToParseArgs
	}

	trace := config.GetBool("logging.trace")
	if trace {
		OutputData(f.UI, config.AllSettings())
	}

	httpClient, err := f.Client()
	if err != nil {
		f.UI.Error(err.Error())
		return exitCodeGenericBad
	}

	expiry := time.Now().UTC().Add(time.Duration(config.GetInt32("client.timeout"))*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), expiry)
	defer cancel()

	var smType raftv1.StateMachineType
	switch f.flagType {
	case "kv":
		smType = raftv1.StateMachineType_STATE_MACHINE_TYPE_KV
	default:
		f.UI.Error("unsupported state machine type")
		return exitCodeGenericBad
	}

	client := raftv1connect.NewShardServiceClient(httpClient, f.BaseCommand.flagHost)

	descriptor, err := client.NewShard(ctx, connect.NewRequest(&raftv1.NewShardRequest{
		ShardId:   f.flagShardId,
		ReplicaId: f.flagReplicaId,
		Type:      smType,
		Hostname:  "",
		Timeout:   f.flagTimeout,
	}))
	if err != nil {
		f.UI.Error(err.Error())
		return exitCodeRemote
	}

	if descriptor!= nil {
		OutputData(f.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (f *FabricAddShardCommand) Synopsis() string {
	return "Create a new shard."
}
