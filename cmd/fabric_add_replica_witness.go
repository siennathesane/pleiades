/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"context"
	"fmt"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-wordwrap"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*FabricAddReplicaWitnessCommand)(nil)
	_ cli.CommandAutocomplete = (*FabricAddReplicaWitnessCommand)(nil)
)

type FabricAddReplicaWitnessCommand struct {
	*BaseCommand

	flagShardId    uint64
	flagReplicaId  uint64
	flagHostname   string
	flagFabricPort uint32
}

func (f *FabricAddReplicaWitnessCommand) Flags() *FlagSets {
	set := f.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging | FlagSetTimeout)
	fs := set.NewFlagSet("Fabric Options")

	fs.Uint64Var(&Uint64Var{
		Name:              "shard-id",
		Usage:             `The ID of the target shard. This is global to the node constellation as it increases the data fabric size.`,
		Target:            &f.flagShardId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-replica-witness.shard-id",
	})

	fs.Uint64Var(&Uint64Var{
		Name:              "replica-id",
		Usage:             `The ID of the new replica. This is specific to each shard.`,
		Target:            &f.flagReplicaId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-replica-witness.replica-id",
	})

	fs.StringVar(&StringVar{
		Name:              "fabric-hostname",
		Usage:             `The internally addressable data fabric hostname where the replica will be created. This address must be accessible by other hosts in the data fabric but not necessarily external to the constellation. For example, if the data fabric is externally accessible at kv.example.io, and the internal fabric nodes are addressable at server-[0,1,2).internal.example.io, operators must use the internal addresses.`,
		Target:            &f.flagHostname,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-replica-observer.fabric-hostname",
	})

	fs.Uint32Var(&Uint32Var{
		Name:              "fabric-port",
		Usage:             `The port the internally addressable data fabric node listens on.`,
		Default:           8081,
		Target:            &f.flagFabricPort,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "fabric.add-replica-witness.fabric-port",
	})

	return set
}

func (f *FabricAddReplicaWitnessCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (f *FabricAddReplicaWitnessCommand) AutocompleteFlags() complete.Flags {
	return f.Flags().Completions()
}

func (f *FabricAddReplicaWitnessCommand) Help() string {
	helpText := `Add a replica witness to a shard.

The data fabric is built on top of sharded, replicated, deterministic finite state machines (FSMs). Each FSM consists of one or more replicas, identified by their replica ID. These replicas allow for distributed FSMs, furthering the durability, performance, and scalability of Pleiades. Replica witnesses are replicas that only participate in voting but do not normally participate in log replication and do not have state machines at all.

When a standard replica goes offline, a replica witness steps in to store log entries for the shard until it recovers or is replaced. Replica witnesses allow the cluster to make consensus decisions even when some of the main servers have failed. It is rare for a data fabric to need a witness, so most users will never need to deploy a witness. This command is for advanced use only.

Replicas observers are created through this command, but they are not started. In order to start a replica, you must call "pleiades fabric start-replica" with the shard and replica IDs. Replicas, observers included, cannot be added to the same host which have the primary shard. The replica ID only matters for uniqueness within a shard, but has no other effect.

Pleiades requires manual replica management right now, but future work will automate the shards and replicas.

` + f.Flags().Help()

	return wordwrap.WrapString(helpText, 80)
}

func (f *FabricAddReplicaWitnessCommand) Run(args []string) int {
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

	fabricHost := fmt.Sprintf("%s:%d", config.GetString("fabric.add-replica-witness.fabric-hostname"), config.GetUint32("fabric.add-replica-witness.fabric-port"))

	if trace {
		f.UI.Info(fmt.Sprintf("setting target fabric host to %s", fabricHost))
	}

	descriptor, err := client.AddReplicaWitness(ctx, connect.NewRequest(&raftv1.AddReplicaWitnessRequest{
		ShardId:   f.flagShardId,
		ReplicaId: f.flagReplicaId,
		Hostname:  fabricHost,
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

func (f *FabricAddReplicaWitnessCommand) Synopsis() string {
	return "Add a new replica witness."
}
