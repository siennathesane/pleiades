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
	"github.com/rs/zerolog"
)

var _ ShardManagerServer = (*raftControlGrpcAdapter)(nil)

type raftControlGrpcAdapter struct {
	logger zerolog.Logger
	clusterManager IShardManager
}

func (r *raftControlGrpcAdapter) AddReplica(ctx context.Context, request *raft.AddReplicaRequest) (*raft.AddReplicaReply, error) {
	cfg := &Config{
		shardId: request.ShardId,
		replicaId: request.ReplicaId,
		stateMachine: StateMachineType(request.Type),
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.AddReplica(cfg, request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
	}
	return &raft.AddReplicaReply{}, err
}

func (r *raftControlGrpcAdapter) AddShardObserver(ctx context.Context, request *raft.AddShardObserverRequest) (*raft.AddShardObserverRequest, error) {
	cfg := &Config{
		shardId: request.ShardId,
		replicaId: request.ReplicaId,
		stateMachine: StateMachineType(request.Type),
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	err := r.clusterManager.AddReplica(cfg, request.Hostname, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
	}
	return &raft.AddShardObserverRequest{}, err
}

func (r *raftControlGrpcAdapter) AddShardWitness(ctx context.Context, request *raft.AddShardWitnessRequest) (*raft.AddShardWitnessRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) DeleteReplica(ctx context.Context, request *raft.DeleteReplicaRequest) (*raft.DeleteReplicaReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) GetLeaderId(ctx context.Context, request *raft.GetLeaderIdRequest) (*raft.GetLeaderIdReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) GetShardMembers(ctx context.Context, request *raft.GetShardMembersRequest) (*raft.GetShardMembersReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) NewShard(ctx context.Context, request *raft.NewShardRequest) (*raft.NewShardReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) RemoveData(ctx context.Context, request *raft.RemoveDataRequest) (*raft.RemoveDataReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) StopReplica(ctx context.Context, request *raft.StopReplicaRequest) (*raft.StopReplicaReply, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftControlGrpcAdapter) mustEmbedUnimplementedShardManagerServer() {
	//TODO implement me
	panic("implement me")
}

