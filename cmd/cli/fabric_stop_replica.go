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
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-wordwrap"
	"github.com/mxplusb/pleiades/pkg/raftpb"
	"github.com/mxplusb/pleiades/pkg/raftpb/raftpbconnect"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*FabricStopReplicaCommand)(nil)
	_ cli.CommandAutocomplete = (*FabricStopReplicaCommand)(nil)
)

type FabricStopReplicaCommand struct {
	*BaseCommand

	flagShardId   uint64
	flagReplicaId uint64
}

func (f *FabricStopReplicaCommand) Flags() *FlagSets {
	set := f.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging | FlagSetTimeout)
	fs := set.NewFlagSet("Fabric Options")

	fs.Uint64Var(&Uint64Var{
		Name: "shard-id",
		Usage: `The ID of the new shard. This is global to the node constellation as it increases the 
data fabric size.`,
		Target:            &f.flagShardId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.stop-replica.shard-id",
	})

	fs.Uint64Var(&Uint64Var{
		Name:              "replica-id",
		Usage:             `The ID of the new replica. This is specific to each shard.`,
		Target:            &f.flagReplicaId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.stop-replica.replica-id",
	})

	return set
}

func (f *FabricStopReplicaCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (f *FabricStopReplicaCommand) AutocompleteFlags() complete.Flags {
	return f.Flags().Completions()
}

// nb (sienna): use word wrap in the editor as this will format properly in the terminal
func (f *FabricStopReplicaCommand) Help() string {
	helpText := `Stop a replica.

This command stops a replica on the targeted node, but it does not remove this replica from the shard. Once stopped, a replica can be restarted with "pleiades fabric start-replica". This command applies to standard, observer, and witness replicas.

` + f.Flags().Help()

	return wordwrap.WrapString(helpText, 80)
}

func (f *FabricStopReplicaCommand) Run(args []string) int {
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

	expiry := time.Now().UTC().Add(time.Duration(config.GetInt32("client.timeout")) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), expiry)
	defer cancel()

	client := raftpbconnect.NewShardServiceClient(httpClient, f.BaseCommand.flagHost)

	descriptor, err := client.StopReplica(ctx, connect.NewRequest(&raftpb.StopReplicaRequest{
		ShardId:   f.flagShardId,
		ReplicaId: f.flagReplicaId,
	}))
	if err != nil {
		f.UI.Error(err.Error())
		return exitCodeRemote
	}

	if descriptor != nil {
		OutputData(f.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (f *FabricStopReplicaCommand) Synopsis() string {
	return "Stop a replica."
}
