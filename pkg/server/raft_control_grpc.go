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

	"github.com/mxplusb/pleiades/pkg/api/v1/raft"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ ShardManagerServer = (*raftControlGrpcAdapter)(nil)

type raftControlGrpcAdapter struct {
	logger zerolog.Logger
	clusterManager IShardManager
}

func (r *raftControlGrpcAdapter) AddReplica(ctx context.Context, request *raft.AddReplicaRequest) (*raft.AddReplicaReply, error) {
	if err := r.checkRequestConfig(request.ShardId, request.ReplicaId); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.GetTimeout()) * time.Millisecond

	err := r.clusterManager.AddReplica(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
		return nil, err
	}
	return &raft.AddReplicaReply{}, err
}

func (r *raftControlGrpcAdapter) AddShardObserver(ctx context.Context, request *raft.AddShardObserverRequest) (*raft.AddShardObserverReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.AddShardObserver(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard observer")
		return nil, err
	}
	return &raft.AddShardObserverReply{}, err
}

func (r *raftControlGrpcAdapter) AddShardWitness(ctx context.Context, request *raft.AddShardWitnessRequest) (*raft.AddShardWitnessReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.AddShardWitness(request.GetShardId(), request.GetReplicaId(), request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard witness")
		return nil, err
	}
	return &raft.AddShardWitnessReply{}, err
}

func (r *raftControlGrpcAdapter) DeleteReplica(ctx context.Context, request *raft.DeleteReplicaRequest) (*raft.DeleteReplicaReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.DeleteReplica(request.GetShardId(), request.GetReplicaId(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't delete replica")
		return nil, err
	}
	return &raft.DeleteReplicaReply{}, err
}

func (r *raftControlGrpcAdapter) GetLeaderId(ctx context.Context, request *raft.GetLeaderIdRequest) (*raft.GetLeaderIdReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.DeleteReplica(request.GetShardId(), request.GetReplicaId(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get leader id")
		return nil, err
	}
	return &raft.GetLeaderIdReply{}, err
}

func (r *raftControlGrpcAdapter) GetShardMembers(ctx context.Context, request *raft.GetShardMembersRequest) (*raft.GetShardMembersReply, error) {
	if request.GetShardId() <= systemShardStop {
		return nil, errors.New("shardId is within system shard range")
	}

	membership, err := r.clusterManager.GetShardMembers(request.ShardId)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get shard members")
		return nil, err
	}
	return &raft.GetShardMembersReply{
		ConfigChangeId: membership.ConfigChangeId,
		Replicas: membership.Replicas,
		Witnesses: membership.Witnesses,
		Observers: membership.Observers,
		Removed: func() map[uint64]string {
			m := make(map[uint64]string)
			for k, _ := range membership.Removed {
				m[k] = ""
			}
			return m
		}(),
	}, err
}

func (r *raftControlGrpcAdapter) NewShard(ctx context.Context, request *raft.NewShardRequest) (*raft.NewShardReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	var t StateMachineType
	switch request.GetType() {
	case raft.StateMachineType_TEST:
		t = testStateMachineType
	case raft.StateMachineType_KV:
		t = BBoltStateMachineType
	default:
		return nil, ErrUnsupportedStateMachine
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.NewShard(request.GetShardId(), request.GetReplicaId(), t, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't create new shard")
		return nil, err
	}
	return &raft.NewShardReply{}, err
}

func (r *raftControlGrpcAdapter) RemoveData(ctx context.Context, request *raft.RemoveDataRequest) (*raft.RemoveDataReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return nil, err
	}

	err := r.clusterManager.RemoveData(request.ShardId, request.ReplicaId)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't remove data from the host")
		return nil, err
	}
	return &raft.RemoveDataReply{}, err
}

func (r *raftControlGrpcAdapter) StopReplica(ctx context.Context, request *raft.StopReplicaRequest) (*raft.StopReplicaReply, error) {
	if request.GetShardId() <= systemShardStop {
		return nil, errors.New("shardId is within system shard range")
	}

	_, err := r.clusterManager.StopReplica(request.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't stop replica")
		return nil, err
	}
	return &raft.StopReplicaReply{}, err
}

func (r *raftControlGrpcAdapter) mustEmbedUnimplementedShardManagerServer() { }

func (r *raftControlGrpcAdapter) checkRequestConfig(shardId, replicaId uint64) error {
	if shardId <= systemShardStop {
		return errors.New("shardId is within system shard range")
	}

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}