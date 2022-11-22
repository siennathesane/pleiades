/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package shard

import (
	"context"
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	_ runtime.IShardManager = (*RaftShardManager)(nil)

	defaultTimeout = 3000 * time.Millisecond

	errNoConfigChangeId = errors.New("no config change id")
)

type ShardManagerBuilderParams struct {
	fx.In

	NodeHost *dragonboat.NodeHost
	Logger   zerolog.Logger
}

type ShardManagerBuilderResults struct {
	fx.Out

	RaftShardManager runtime.IShardManager
}

func NewManager(params ShardManagerBuilderParams) ShardManagerBuilderResults {
	l := params.Logger.With().Str("component", "shard-manager").Logger()
	return ShardManagerBuilderResults{
		RaftShardManager: &RaftShardManager{l, params.NodeHost},
	}
}

type RaftShardManager struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (c *RaftShardManager) AddReplica(req *raftv1.AddReplicaRequest) error {
	l := c.logger.With().Uint64("shard", req.GetShardId()).Uint64("replica", req.GetReplicaId()).Logger()
	l.Info().Msg("adding replica")

	if req.GetTimeout() == 0 {
		l.Debug().Msg("using default timeout")
		req.Timeout = int64(defaultTimeout)
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(req.GetTimeout())))
	defer cancel()

	members, err := c.GetShardMembers(req.GetShardId())
	if err != nil {
		l.Error().Err(err).Msg("failed to get shard members")
		return err
	}

	err = c.nh.SyncRequestAddNode(ctx, req.GetShardId(), req.GetReplicaId(), req.GetHostname(), members.ConfigChangeId)
	if err != nil {
		l.Error().Err(err).Msg("failed to add replica")
		return err
	}

	return nil
}

func (c *RaftShardManager) AddReplicaObserver(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("adding shard observer")

	if timeout == 0 {
		l.Debug().Msg("using default timeout")
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	members, err := c.GetShardMembers(shardId)
	if err != nil {
		l.Error().Err(err).Msg("failed to get shard members")
		return err
	}

	err = c.nh.SyncRequestAddObserver(ctx, shardId, replicaId, newHost, members.ConfigChangeId)
	if err != nil {
		l.Error().Err(err).Msg("failed to add shard observer")
		return err
	}

	return nil
}

func (c *RaftShardManager) AddReplicaWitness(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("adding shard observer")

	if timeout == 0 {
		l.Debug().Msg("using default timeout")
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	members, err := c.GetShardMembers(shardId)
	if err != nil {
		l.Error().Err(err).Msg("failed to get shard members")
		return err
	}

	err = c.nh.SyncRequestAddWitness(ctx, shardId, replicaId, newHost, members.ConfigChangeId)
	if err != nil {
		l.Error().Err(err).Msg("failed to add shard observer")
		return err
	}

	return nil
}

func (c *RaftShardManager) GetLeaderId(shardId uint64) (leader uint64, ok bool, err error) {
	c.logger.Info().Uint64("shard", shardId).Msg("getting leader id")
	return c.nh.GetLeaderID(shardId)
}

func (c *RaftShardManager) GetShardMembers(shardId uint64) (*runtime.MembershipEntry, error) {
	l := c.logger.With().Uint64("shard", shardId).Logger()
	l.Debug().Msg("getting shard members")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(defaultTimeout))
	defer cancel()

	membership, err := c.nh.SyncGetClusterMembership(ctx, shardId)
	if err != nil {
		l.Error().Err(err).Msg("failed to get shard membership")
		return nil, err
	}

	return &runtime.MembershipEntry{
		ConfigChangeId: membership.ConfigChangeID,
		Replicas:       membership.Nodes,
		Observers:      membership.Observers,
		Witnesses:      membership.Witnesses,
	}, nil
}

func (c *RaftShardManager) NewShard(req *raftv1.NewShardRequest) error {
	l := c.logger.With().Uint64("shard", req.GetShardId()).Uint64("replica", req.GetReplicaId()).Logger()
	l.Info().Msg("creating new shard")

	clusterConfig := newDConfig(req.GetShardId(), req.GetReplicaId())

	members := make(map[uint64]string)
	members[req.GetReplicaId()] = c.nh.RaftAddress()
	l.Debug().Str("raft-address", c.nh.RaftAddress()).Msg("adding self to members")

	var err error
	switch req.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		l.Info().Msg("creating test state machine")
		err = c.nh.StartCluster(members, false, serverutils.NewTestStateMachine, clusterConfig)
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		l.Info().Msg("creating bbolt state machine")
		err = c.nh.StartOnDiskCluster(members, false, kv.NewBBoltFSM, clusterConfig)
	default:
		l.Error().Msg("unknown state machine type")
		return ErrUnsupportedStateMachine
	}

	if err != nil {
		l.Error().Err(err).Msg("error creating state machine")
		return err
	}

	return nil
}

func (c *RaftShardManager) StartReplica(req *raftv1.StartReplicaRequest) error {
	l := c.logger.With().Uint64("shard", req.GetShardId()).Uint64("replica", req.GetReplicaId()).Logger()
	l.Info().Msg("starting replica")

	clusterConfig := newDConfig(req.GetShardId(), req.GetReplicaId())

	var err error
	switch req.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		if !req.GetRestart() {
			l.Info().Msg("starting test state machine")
			err = c.nh.StartCluster(nil, true, serverutils.NewTestStateMachine, clusterConfig)
		} else {
			l.Info().Msg("restarting test state machine")
			err = c.nh.StartCluster(nil, false, serverutils.NewTestStateMachine, clusterConfig)
		}
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		if !req.GetRestart() {
			l.Info().Msg("starting bbolt state machine")
			err = c.nh.StartOnDiskCluster(nil, true, kv.NewBBoltFSM, clusterConfig)
		} else {
			l.Info().Msg("restarting bbolt state machine")
			err = c.nh.StartOnDiskCluster(nil, false, kv.NewBBoltFSM, clusterConfig)
		}
	default:
		l.Error().Msg("unknown state machine type")
		return ErrUnsupportedStateMachine
	}

	if err != nil {
		l.Error().Err(err).Msg("error starting state machine")
		return err
	}

	return nil
}

func (c *RaftShardManager) StartReplicaObserver(req *raftv1.StartReplicaObserverRequest) error {
	l := c.logger.With().Uint64("shard", req.GetShardId()).Uint64("replica", req.GetReplicaId()).Logger()
	l.Info().Msg("starting replica observer")

	clusterConfig := newDConfig(req.GetShardId(), req.GetReplicaId())
	clusterConfig.IsObserver = true

	var err error
	switch req.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		l.Info().Msg("starting test state machine observer")
		if !req.GetRestart() {
			err = c.nh.StartCluster(nil, true, serverutils.NewTestStateMachine, clusterConfig)
		} else {
			l.Info().Msg("restarting test state machine observer")
			err = c.nh.StartCluster(nil, false, serverutils.NewTestStateMachine, clusterConfig)
		}
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		l.Info().Msg("starting bbolt state machine observer")
		if !req.GetRestart() {
			err = c.nh.StartOnDiskCluster(nil, true, kv.NewBBoltFSM, clusterConfig)
		} else {
			l.Info().Msg("restarting bbolt state machine observer")
			err = c.nh.StartOnDiskCluster(nil, false, kv.NewBBoltFSM, clusterConfig)
		}
	default:
		l.Error().Msg("unknown state machine type")
		return ErrUnsupportedStateMachine
	}

	if err != nil {
		l.Error().Err(err).Msg("error starting state machine observer")
		return err
	}

	return nil
}

// StopReplica will stop a replica
// todo (sienna): this should stop the replica, not the shard...
func (c *RaftShardManager) StopReplica(shardId uint64, replicaId uint64) (*runtime.OperationResult, error) {
	c.logger.Info().Uint64("shard", shardId).Msg("stopping replica")
	err := c.nh.StopNode(shardId, replicaId)
	if err != nil {
		c.logger.Error().Err(err).Uint64("shard", shardId).Msg("failed to stop replica")
		return nil, errors.Wrap(err, "failed to stop replica")
	}

	return &runtime.OperationResult{}, nil
}

func (c *RaftShardManager) RemoveReplica(shardId uint64, replicaId uint64, timeout time.Duration) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("deleting replica")

	if timeout == 0 {
		l.Debug().Msg("using default timeout")
		timeout = defaultTimeout
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	members, err := c.GetShardMembers(shardId)
	if err != nil {
		l.Error().Err(err).Msg("failed to get shard members")
		return err
	}

	if members.ConfigChangeId == 0 {
		l.Error().Err(errNoConfigChangeId).Msg("failed to get config change id from shard members")
		return errNoConfigChangeId
	}

	err = c.nh.SyncRequestDeleteNode(ctx, shardId, replicaId, members.ConfigChangeId)
	if err != nil {
		l.Error().Err(err).Msg("failed to delete replica")
		return err
	}

	return nil
}

func (c *RaftShardManager) RemoveData(shardId, replicaId uint64) error {
	c.logger.Info().Uint64("shard", shardId).Uint64("replica", replicaId).Msg("removing data")

	err := c.nh.RemoveData(shardId, replicaId)
	if err != nil {
		c.logger.Error().Err(err).Uint64("shard", shardId).Uint64("replica", replicaId).Msg("failed to remove data")
	}
	return err
}

func newDConfig(shardId, replicaId uint64) dconfig.Config {
	// nb (sienna): if you change this outside of a major version rollout,
	// it will create inconsistencies across all clusters. don't change this
	// without running it through me. the inconsistencies this will create
	// will affect anything built on top of all state machines, which is a huge
	// business risk that requires my sign-off. it can be changed if there's a
	// need, but ensure I sign off on it first. 🙂
	return dconfig.Config{
		NodeID:                  replicaId,
		ClusterID:               shardId,
		CheckQuorum:             true,
		ElectionRTT:             100,
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
