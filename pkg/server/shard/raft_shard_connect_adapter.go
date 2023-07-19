/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package shard

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/cockroachdb/errors"
	raftpb "github.com/mxplusb/pleiades/pkg/raftpb"
	"github.com/mxplusb/pleiades/pkg/raftpb/raftpbconnect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var _ raftpbconnect.ShardServiceHandler = (*RaftShardConnectAdapter)(nil)
var _ runtime.ServiceHandler = (*RaftShardConnectAdapter)(nil)

type ConnectAdapterBuilderParams struct {
	fx.In

	Logger       zerolog.Logger
	ShardManager runtime.IShardManager
}

type ConnectAdapterBuilderResults struct {
	fx.Out

	ConnectAdapter raftpbconnect.ShardServiceHandler
}

type RaftShardConnectAdapter struct {
	http.Handler
	logger       zerolog.Logger
	path         string
	shardManager runtime.IShardManager
}

func NewRaftShardConnectAdapter(shardManager runtime.IShardManager, logger zerolog.Logger) *RaftShardConnectAdapter {
	adapter := &RaftShardConnectAdapter{logger: logger, shardManager: shardManager}
	adapter.path, adapter.Handler = raftpbconnect.NewShardServiceHandler(adapter)

	return adapter
}

func (r *RaftShardConnectAdapter) Path() string {
	return r.path
}

func (r *RaftShardConnectAdapter) AddReplica(ctx context.Context, c *connect.Request[raftpb.AddReplicaRequest]) (*connect.Response[raftpb.AddReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.AddReplicaResponse{}), err
	}

	err := r.shardManager.AddReplica(c.Msg)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add replica")
		return nil, err
	}
	return connect.NewResponse(&raftpb.AddReplicaResponse{}), err
}

func (r *RaftShardConnectAdapter) AddReplicaObserver(ctx context.Context, c *connect.Request[raftpb.AddReplicaObserverRequest]) (*connect.Response[raftpb.AddReplicaObserverResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.AddReplicaObserverResponse{}), err
	}

	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplicaObserver(c.Msg.GetShardId(), c.Msg.GetReplicaId(), c.Msg.GetHostname(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard observer")
		return nil, err
	}
	return connect.NewResponse(&raftpb.AddReplicaObserverResponse{}), err
}

func (r *RaftShardConnectAdapter) AddReplicaWitness(ctx context.Context, c *connect.Request[raftpb.AddReplicaWitnessRequest]) (*connect.Response[raftpb.AddReplicaWitnessResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.AddReplicaWitnessResponse{}), err
	}

	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	err := r.shardManager.AddReplicaWitness(c.Msg.GetShardId(), c.Msg.GetReplicaId(), c.Msg.GetHostname(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add shard witness")
		return nil, err
	}
	return connect.NewResponse(&raftpb.AddReplicaWitnessResponse{}), err
}

func (r *RaftShardConnectAdapter) GetLeaderId(ctx context.Context, c *connect.Request[raftpb.GetLeaderIdRequest]) (*connect.Response[raftpb.GetLeaderIdResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.GetLeaderIdResponse{}), err
	}

	leader, ok, err := r.shardManager.GetLeaderId(c.Msg.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get leader id")
		return nil, err
	}
	return connect.NewResponse(&raftpb.GetLeaderIdResponse{
		Leader:    leader,
		Available: ok,
	}), err
}

func (r *RaftShardConnectAdapter) GetShardMembers(ctx context.Context, c *connect.Request[raftpb.GetShardMembersRequest]) (*connect.Response[raftpb.GetShardMembersResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&raftpb.GetShardMembersResponse{}), errors.New("invalid shard id, must not be 0")
	}

	membership, err := r.shardManager.GetShardMembers(c.Msg.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't get shard members")
		return nil, err
	}
	return connect.NewResponse(&raftpb.GetShardMembersResponse{
		ConfigChangeId: membership.ConfigChangeId,
		Replicas:       membership.Replicas,
		Witnesses:      membership.Witnesses,
		Observers:      membership.Observers,
		Removed: func() map[uint64]string {
			m := make(map[uint64]string)
			for k := range membership.Removed {
				m[k] = ""
			}
			return m
		}(),
	}), err
}

func (r *RaftShardConnectAdapter) NewShard(ctx context.Context, c *connect.Request[raftpb.NewShardRequest]) (*connect.Response[raftpb.NewShardResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.NewShardResponse{}), err
	}

	r.logger.Trace().Str("state-machine", c.Msg.GetType().String()).Msg("state machine is supported")

	err := r.shardManager.NewShard(c.Msg)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't create new shard")
		return nil, err
	}
	return connect.NewResponse(&raftpb.NewShardResponse{}), err
}

func (r *RaftShardConnectAdapter) RemoveData(ctx context.Context, c *connect.Request[raftpb.RemoveDataRequest]) (*connect.Response[raftpb.RemoveDataResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.RemoveDataResponse{}), err
	}

	err := r.shardManager.RemoveData(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't remove data from the host")
		return nil, err
	}
	return connect.NewResponse(&raftpb.RemoveDataResponse{}), err
}

func (r *RaftShardConnectAdapter) RemoveReplica(ctx context.Context, c *connect.Request[raftpb.RemoveReplicaRequest]) (*connect.Response[raftpb.RemoveReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.RemoveReplicaResponse{}), err
	}

	timeout := time.Duration(c.Msg.Timeout) * time.Millisecond

	err := r.shardManager.RemoveReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId(), timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't delete replica")
		return nil, err
	}
	return connect.NewResponse(&raftpb.RemoveReplicaResponse{}), err
}

func (r *RaftShardConnectAdapter) StartReplica(ctx context.Context, c *connect.Request[raftpb.StartReplicaRequest]) (*connect.Response[raftpb.StartReplicaResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.StartReplicaResponse{}), err
	}

	err := r.shardManager.StartReplica(c.Msg)

	return connect.NewResponse(&raftpb.StartReplicaResponse{}), err
}

func (r *RaftShardConnectAdapter) StartReplicaObserver(ctx context.Context, c *connect.Request[raftpb.StartReplicaObserverRequest]) (*connect.Response[raftpb.StartReplicaObserverResponse], error) {
	if err := r.checkRequestConfig(c.Msg.GetShardId(), c.Msg.GetReplicaId()); err != nil {
		return connect.NewResponse(&raftpb.StartReplicaObserverResponse{}), err
	}

	err := r.shardManager.StartReplicaObserver(c.Msg)

	return connect.NewResponse(&raftpb.StartReplicaObserverResponse{}), err
}

func (r *RaftShardConnectAdapter) StopReplica(ctx context.Context, c *connect.Request[raftpb.StopReplicaRequest]) (*connect.Response[raftpb.StopReplicaResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&raftpb.StopReplicaResponse{}), errors.New("invalid shard id, must not be 0")
	}

	_, err := r.shardManager.StopReplica(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	if err != nil {
		r.logger.Error().Err(err).Msg("can't stop replica")
		return nil, err
	}
	return connect.NewResponse(&raftpb.StopReplicaResponse{}), err
}

func (r *RaftShardConnectAdapter) checkRequestConfig(shardId, replicaId uint64) error {

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}
