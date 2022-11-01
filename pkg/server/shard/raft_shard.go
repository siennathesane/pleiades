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

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/server/eventing"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ runtime.IShardManager = (*RaftShardManager)(nil)

	defaultTimeout = 3000 * time.Millisecond

	errNoConfigChangeId = errors.New("no config change id")
)

func NewShardManager(nodeHost *dragonboat.NodeHost, client *messaging.EmbeddedMessagingStreamClient, raftEventHandler *messaging.RaftEventHandler, logger zerolog.Logger) *RaftShardManager {
	l := logger.With().Str("component", "shard-manager").Logger()
	return &RaftShardManager{l, nodeHost, client, raftEventHandler}
}

type RaftShardManager struct {
	logger           zerolog.Logger
	nh               *dragonboat.NodeHost
	client           *messaging.EmbeddedMessagingStreamClient
	raftEventHandler *messaging.RaftEventHandler
}

func (c *RaftShardManager) AddReplica(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("adding replica")

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

	err = c.nh.SyncRequestAddNode(ctx, shardId, replicaId, newHost, members.ConfigChangeId)
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

func (c *RaftShardManager) NewShard(shardId uint64, replicaId uint64, stateMachineType runtime.StateMachineType, _ time.Duration) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("creating new shard")

	clusterConfig := newDConfig(shardId, replicaId)

	members := make(map[uint64]string)
	members[replicaId] = c.nh.RaftAddress()
	l.Debug().Str("raft-address", c.nh.RaftAddress()).Msg("adding self to members")

	var err error
	switch stateMachineType {
	case testStateMachineType:
		l.Info().Msg("creating test state machine")
		err = c.nh.StartCluster(members, false, serverutils.NewTestStateMachine, clusterConfig)
	case BBoltStateMachineType:
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

func (c *RaftShardManager) StartReplica(shardId uint64, replicaId uint64, stateMachineType runtime.StateMachineType) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("starting replica")

	clusterConfig := newDConfig(shardId, replicaId)

	var err error
	switch stateMachineType {
	case testStateMachineType:
		l.Info().Msg("starting test state machine")
		err = c.nh.StartCluster(nil, true, serverutils.NewTestStateMachine, clusterConfig)
	case BBoltStateMachineType:
		l.Info().Msg("starting bbolt state machine")
		err = c.nh.StartOnDiskCluster(nil, true, kv.NewBBoltFSM, clusterConfig)
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

func (c *RaftShardManager) StartReplicaObserver(shardId uint64, replicaId uint64, stateMachineType runtime.StateMachineType) error {
	l := c.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()
	l.Info().Msg("starting replica observer")

	clusterConfig := newDConfig(shardId, replicaId)
	clusterConfig.IsObserver = true

	var err error
	switch stateMachineType {
	case testStateMachineType:
		l.Info().Msg("starting test state machine observer")
		return c.nh.StartCluster(nil, true, serverutils.NewTestStateMachine, clusterConfig)
	case BBoltStateMachineType:
		l.Info().Msg("starting bbolt state machine observer")
		err = c.nh.StartOnDiskCluster(nil, true, kv.NewBBoltFSM, clusterConfig)
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

func (c *RaftShardManager) storeShardState(shardId uint64, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan *raftv1.RaftEvent)
	go c.raftEventHandler.WaitForMembershipChange(shardId, resultChan, timeout)
	<-resultChan

	memberState, err := c.nh.SyncGetClusterMembership(ctx, shardId)
	if err != nil {
		c.logger.Error().Err(err).Msg("can't get shard state")
	}

	// if it's nil, we can't do anything
	if memberState == nil {
		return
	}

	msg := &raftv1.ShardState{
		LastUpdated:    timestamppb.Now(),
		ShardId:        shardId,
		ConfigChangeId: memberState.ConfigChangeID,
		Replicas:       memberState.Nodes,
		Observers:      memberState.Observers,
		Witnesses:      memberState.Witnesses,
		Removed: func() map[uint64]string {
			m := make(map[uint64]string)
			for k := range memberState.Removed {
				m[k] = ""
			}
			return m
		}(),
		Type: 0,
	}

	payload, err := msg.MarshalVT()
	if err != nil {
		c.logger.Error().Err(err).Msg("can't marshal shard state message")
	}

	_, err = c.client.Publish(eventing.ShardConfigStream, payload)
	if err != nil {
		c.logger.Error().Err(err).Msg("can't publish shard state event")
	}
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
