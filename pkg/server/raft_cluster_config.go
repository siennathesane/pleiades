/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	dconfig "github.com/lni/dragonboat/v3/config"
)

type IClusterConfig interface {
	Adapt() dconfig.Config
	StateMachineType() StateMachineType
	ReplicaId() uint64
	ShardId() uint64
}

var _ IClusterConfig = (*Config)(nil)

// Config is used to configure Raft nodes.
type Config struct {
	// replicaId is a non-zero value used to identify a node within a Raft cluster.
	replicaId uint64
	// shardId is the unique value used to identify a Raft cluster.
	shardId uint64
	// stateMachine dictates the type of state machine
	stateMachine StateMachineType
}

func (c *Config) ReplicaId() uint64 {
	return c.replicaId
}

func (c *Config) ShardId() uint64 {
	return c.shardId
}

func (c *Config) StateMachineType() StateMachineType {
	return c.stateMachine
}

func (c *Config) Adapt() dconfig.Config {
	return dconfig.Config{
		NodeID:                  c.replicaId,
		ClusterID:               c.shardId,
		CheckQuorum:             true,
		ElectionRTT:             2,
		HeartbeatRTT:            10,
		SnapshotEntries:         1000,
		CompactionOverhead:      500,
		OrderedConfigChange:     true,
		MaxInMemLogSize:         0,
		SnapshotCompressionType: 0,
		EntryCompressionType:    0,
		DisableAutoCompactions:  false,
		IsObserver:              false,
		IsWitness:               false,
		Quiesce:                 false,
	}
}