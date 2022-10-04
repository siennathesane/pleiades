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

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ raftv1.ShardServiceServer = (*raftShardGrpcAdapter)(nil)

type raftShardGrpcAdapter struct {
	logger       zerolog.Logger
	shardManager IShardManager
}

func (r *raftShardGrpcAdapter) AddReplica(ctx context.Context, request *raftv1.AddReplicaRequest) (*raftv1.AddReplicaResponse, error) {
	if err := r.checkRequestConfig(request.ShardId, request.ReplicaId); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplica(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
		return nil, err
	}
	return &raftv1.AddReplicaResponse{}, err
}

func (r *raftShardGrpcAdapter) AddReplicaObserver(ctx context.Context, request *raftv1.AddReplicaObserverRequest) (*raftv1.AddReplicaObserverResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.shardManager.AddReplicaObserver(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard observer")
		return nil, err
	}
	return &raftv1.AddReplicaObserverResponse{}, err
}

func (r *raftShardGrpcAdapter) AddReplicaWitness(ctx context.Context, request *raftv1.AddReplicaWitnessRequest) (*raftv1.AddReplicaWitnessResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.shardManager.AddReplicaWitness(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard witness")
		return nil, err
	}
	return &raftv1.AddReplicaWitnessResponse{}, err
}

func (r *raftShardGrpcAdapter) GetLeaderId(ctx context.Context, request *raftv1.GetLeaderIdRequest) (*raftv1.GetLeaderIdResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	leader, ok, err := r.shardManager.GetLeaderId(request.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get leader id")
		return nil, err
	}
	return &raftv1.GetLeaderIdResponse{
		Leader:    leader,
		Available: ok,
	}, err
}

func (r *raftShardGrpcAdapter) GetShardMembers(ctx context.Context, request *raftv1.GetShardMembersRequest) (*raftv1.GetShardMembersResponse, error) {
	if request.GetShardId() <= systemShardStop {
		return nil, ErrSystemShardRange
	}

	membership, err := r.shardManager.GetShardMembers(request.ShardId)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get shard members")
		return nil, err
	}
	return &raftv1.GetShardMembersResponse{
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
	}, err
}

func (r *raftShardGrpcAdapter) NewShard(ctx context.Context, request *raftv1.NewShardRequest) (*raftv1.NewShardResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	var t StateMachineType
	switch request.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.shardManager.NewShard(request.GetShardId(), request.GetReplicaId(), t, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't create new shard")
		return nil, err
	}
	return &raftv1.NewShardResponse{}, err
}

func (r *raftShardGrpcAdapter) RemoveData(ctx context.Context, request *raftv1.RemoveDataRequest) (*raftv1.RemoveDataResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	err := r.shardManager.RemoveData(request.ShardId, request.ReplicaId)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't remove data from the host")
		return nil, err
	}
	return &raftv1.RemoveDataResponse{}, err
}

func (r *raftShardGrpcAdapter) RemoveReplica(ctx context.Context, request *raftv1.RemoveReplicaRequest) (*raftv1.RemoveReplicaResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.shardManager.RemoveReplica(request.GetShardId(), request.GetReplicaId(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't delete replica")
		return nil, err
	}
	return &raftv1.RemoveReplicaResponse{}, err
}

func (r *raftShardGrpcAdapter) StartReplica(ctx context.Context, request *raftv1.StartReplicaRequest) (*raftv1.StartReplicaResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	var t StateMachineType
	switch request.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	err := r.shardManager.StartReplica(request.GetShardId(), request.GetReplicaId(), t)

	return &raftv1.StartReplicaResponse{}, err
}

func (r *raftShardGrpcAdapter) StartReplicaObserver(ctx context.Context, request *raftv1.StartReplicaObserverRequest) (*raftv1.StartReplicaObserverResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	var t StateMachineType
	switch request.GetType() {
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST:
		t = testStateMachineType
	case raftv1.StateMachineType_STATE_MACHINE_TYPE_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	err := r.shardManager.StartReplicaObserver(request.GetShardId(), request.GetReplicaId(), t)

	return &raftv1.StartReplicaObserverResponse{}, err
}

func (r *raftShardGrpcAdapter) StopReplica(ctx context.Context, request *raftv1.StopReplicaRequest) (*raftv1.StopReplicaResponse, error) {
	if request.GetShardId() <= systemShardStop {
		return nil, ErrSystemShardRange
	}

	_, err := r.shardManager.StopReplica(request.GetShardId(), 0)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't stop replica")
		return nil, err
	}
	return &raftv1.StopReplicaResponse{}, err
}

func (r *raftShardGrpcAdapter) mustEmbedUnimplementedShardManagerServer() {}

func (r *raftShardGrpcAdapter) checkRequestConfig(shardId, replicaId uint64) error {
	if shardId <= systemShardStop {
		return ErrSystemShardRange
	}

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}
