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
	"fmt"
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/api/raft/v1/raftv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-wordwrap"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*FabricGetLeaderIdCommand)(nil)
	_ cli.CommandAutocomplete = (*FabricGetLeaderIdCommand)(nil)
)

type FabricGetLeaderIdCommand struct {
	*BaseCommand

	flagShardId    uint64
}

func (f *FabricGetLeaderIdCommand) Flags() *FlagSets {
	set := f.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging | FlagSetTimeout)
	fs := set.NewFlagSet("Fabric Options")

	fs.Uint64Var(&Uint64Var{
		Name:              "shard-id",
		Usage:             `The ID of the target shard. This is global to the node constellation as it increases the data fabric size.`,
		Target:            &f.flagShardId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.get-leader-id.shard-id",
	})

	return set
}

func (f *FabricGetLeaderIdCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (f *FabricGetLeaderIdCommand) AutocompleteFlags() complete.Flags {
	return f.Flags().Completions()
}

func (f *FabricGetLeaderIdCommand) Help() string {
	helpText := `Get a shard's leader ID.

The data fabric is built on top of sharded, replicated, deterministic finite state machines (FSMs). Each FSM consists of one or more replicas, identified by their replica ID. These replicas allow for distributed FSMs, furthering the durability, performance, and scalability of Pleiades.

This command will query the data fabric node's shard and return the ID of the leader from it's perspective. If the data fabric node doesn't contain an instance of the target shard, the shard is currently electing a new leader, or quorum cannot be guaranteed, the result will be 0.

` + f.Flags().Help()

	return wordwrap.WrapString(helpText, 80)
}

func (f *FabricGetLeaderIdCommand) Run(args []string) int {
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

	if trace {
		f.UI.Info(fmt.Sprintf("operation expires on %s", expiry.Local()))
	}

	ctx, cancel := context.WithDeadline(context.Background(), expiry)
	defer cancel()

	client := raftv1connect.NewShardServiceClient(httpClient, f.BaseCommand.flagHost)

	fabricHost := fmt.Sprintf("%s:%d", config.GetString("fabric.get-leader-id.fabric-hostname"), config.GetUint32("fabric.get-leader-id.fabric-port"))

	if trace {
		f.UI.Info(fmt.Sprintf("setting target fabric host to %s", fabricHost))
	}

	descriptor, err := client.GetLeaderId(ctx, connect.NewRequest(&raftv1.GetLeaderIdRequest{
		ShardId:   f.flagShardId,
		Timeout:   int64(config.GetInt32("client.timeout")),
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

func (f *FabricGetLeaderIdCommand) Synopsis() string {
	return "Get a shard's leader ID."
}
