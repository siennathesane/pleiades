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
	_ cli.Command             = (*FabricRemoveDataCommand)(nil)
	_ cli.CommandAutocomplete = (*FabricRemoveDataCommand)(nil)
)

type FabricRemoveDataCommand struct {
	*BaseCommand

	flagShardId   uint64
	flagReplicaId uint64
}

func (f *FabricRemoveDataCommand) Flags() *FlagSets {
	set := f.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging | FlagSetTimeout)
	fs := set.NewFlagSet("Fabric Options")

	fs.Uint64Var(&Uint64Var{
		Name: "shard-id",
		Usage: `The ID of the new shard. This is global to the node constellation as it increases the 
data fabric size.`,
		Target:            &f.flagShardId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.remove-data.shard-id",
	})

	fs.Uint64Var(&Uint64Var{
		Name:              "replica-id",
		Usage:             `The ID of the new replica. This is specific to each shard.`,
		Target:            &f.flagReplicaId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.remove-data.replica-id",
	})

	return set
}

func (f *FabricRemoveDataCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (f *FabricRemoveDataCommand) AutocompleteFlags() complete.Flags {
	return f.Flags().Completions()
}

// nb (sienna): use word wrap in the editor as this will format properly in the terminal
func (f *FabricRemoveDataCommand) Help() string {
	helpText := `Remove a replica's data.

This command tries to remove all data associated with the specified replica. This command should only be used after the replica has been deleted from its shard. Calling this command on a replica that is still a shard member will corrupt the shard. This command returns an error when the specified node has not been fully offloaded from the node.

` + f.Flags().Help()

	return wordwrap.WrapString(helpText, 80)
}

func (f *FabricRemoveDataCommand) Run(args []string) int {
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

	client := raftv1connect.NewShardServiceClient(httpClient, f.BaseCommand.flagHost)

	descriptor, err := client.RemoveData(ctx, connect.NewRequest(&raftv1.RemoveDataRequest{
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

func (f *FabricRemoveDataCommand) Synopsis() string {
	return "Remove a replica's data."
}
