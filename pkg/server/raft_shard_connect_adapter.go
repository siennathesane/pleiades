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
	"context"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/bufbuild/connect-go"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ raftv1connect.ShardServiceHandler = (*raftShardConnectAdapter)(nil)

type raftShardConnectAdapter struct {
	logger       zerolog.Logger
	shardManager IShardManager
}

func (r *raftShardConnectAdapter) AddReplica(ctx context.Context, c *connect.Request[raftv1.AddReplicaRequest]) (*connect.Response[raftv1.AddReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.AddReplicaResponse{}), err
	}

	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId(), c.Msg.GetHostname(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
		return nil, err
	}
	return connect.NewResponse(&raftv1.AddReplicaResponse{}), err
}

func (r *raftShardConnectAdapter) AddReplicaObserver(ctx context.Context, c *connect.Request[raftv1.AddReplicaObserverRequest]) (*connect.Response[raftv1.AddReplicaObserverResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.AddReplicaObserverResponse{}), err
	}

	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplicaObserver(c.Msg.GetShardId(), c.Msg.GetReplicaId(), c.Msg.GetHostname(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard observer")
		return nil, err
	}
	return connect.NewResponse(&raftv1.AddReplicaObserverResponse{}), err
}

func (r *raftShardConnectAdapter) AddReplicaWitness(ctx context.Context, c *connect.Request[raftv1.AddReplicaWitnessRequest]) (*connect.Response[raftv1.AddReplicaWitnessResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.AddReplicaWitnessResponse{}), err
	}

	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplicaWitness(c.Msg.GetShardId(), c.Msg.GetReplicaId(), c.Msg.GetHostname(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard witness")
		return nil, err
	}
	return connect.NewResponse(&raftv1.AddReplicaWitnessResponse{}), err
}

func (r *raftShardConnectAdapter) GetLeaderId(ctx context.Context, c *connect.Request[raftv1.GetLeaderIdRequest]) (*connect.Response[raftv1.GetLeaderIdResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.GetLeaderIdResponse{}), err
	}

	leader, ok, err := r.shardManager.GetLeaderId(c.Msg.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get leader id")
		return nil, err
	}
	return connect.NewResponse(&raftv1.GetLeaderIdResponse{
		Leader:    leader,
		Available: ok,
	}), err
}

func (r *raftShardConnectAdapter) GetShardMembers(ctx context.Context, c *connect.Request[raftv1.GetShardMembersRequest]) (*connect.Response[raftv1.GetShardMembersResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&raftv1.GetShardMembersResponse{}), errors.New("invalid shard id, must not be 0")
	}

	membership, err := r.shardManager.GetShardMembers(c.Msg.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get shard members")
		return nil, err
	}
	return connect.NewResponse(&raftv1.GetShardMembersResponse{
		ConfigChangeId: membership.ConfigChangeId,
		Replicas:       membership.Replicas,
		Witnesses:      membership.Witnesses,
		Observers:      membership.Observers,
		Removed: func() map[uint64]string {
			m := make(map[uint64]string)
			for k, _ := range membership.Removed {
				m[k] = ""
			}
			return m
		}(),
	}), err
}

func (r *raftShardConnectAdapter) NewShard(ctx context.Context, c *connect.Request[raftv1.NewShardRequest]) (*connect.Response[raftv1.NewShardResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.NewShardResponse{}), err
	}

	var t StateMachineType
	switch c.Msg.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	timeout := time.Duration(c.Msg.Timeout) * time.Millisecond

	err := r.shardManager.NewShard(c.Msg.GetShardId(), c.Msg.GetReplicaId(), t, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't create new shard")
		return nil, err
	}
	return connect.NewResponse(&raftv1.NewShardResponse{}), err
}

func (r *raftShardConnectAdapter) RemoveData(ctx context.Context, c *connect.Request[raftv1.RemoveDataRequest]) (*connect.Response[raftv1.RemoveDataResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.RemoveDataResponse{}), err
	}

	err := r.shardManager.RemoveData(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't remove data from the host")
		return nil, err
	}
	return connect.NewResponse(&raftv1.RemoveDataResponse{}), err
}

func (r *raftShardConnectAdapter) RemoveReplica(ctx context.Context, c *connect.Request[raftv1.RemoveReplicaRequest]) (*connect.Response[raftv1.RemoveReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.RemoveReplicaResponse{}), err
	}

	timeout := time.Duration(c.Msg.Timeout) * time.Millisecond

	err := r.shardManager.RemoveReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't delete replica")
		return nil, err
	}
	return connect.NewResponse(&raftv1.RemoveReplicaResponse{}), err
}

func (r *raftShardConnectAdapter) StartReplica(ctx context.Context, c *connect.Request[raftv1.StartReplicaRequest]) (*connect.Response[raftv1.StartReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.StartReplicaResponse{}), err
	}

	var t StateMachineType
	switch c.Msg.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	err := r.shardManager.StartReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId(), t)

	return connect.NewResponse(&raftv1.StartReplicaResponse{}), err
}

func (r *raftShardConnectAdapter) StartReplicaObserver(ctx context.Context, c *connect.Request[raftv1.StartReplicaObserverRequest]) (*connect.Response[raftv1.StartReplicaObserverResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftv1.StartReplicaObserverResponse{}), err
	}

	var t StateMachineType
	switch c.Msg.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	err := r.shardManager.StartReplicaObserver(c.Msg.GetShardId(), c.Msg.GetReplicaId(), t)

	return connect.NewResponse(&raftv1.StartReplicaObserverResponse{}), err
}

func (r *raftShardConnectAdapter) StopReplica(ctx context.Context, c *connect.Request[raftv1.StopReplicaRequest]) (*connect.Response[raftv1.StopReplicaResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&raftv1.StopReplicaResponse{}), errors.New("invalid shard id, must not be 0")
	}

	_, err := r.shardManager.StopReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't stop replica")
		return nil, err
	}
	return connect.NewResponse(&raftv1.StopReplicaResponse{}), err
}

func (r *raftShardConnectAdapter) checkRequestConfig(shardId, replicaId uint64) error {
	if shardId <= systemShardStop {
		return ErrSystemShardRange
	}

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}
