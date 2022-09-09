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

	"github.com/mxplusb/pleiades/api/v1/raft"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ RaftHostServer = (*raftHostGrpcAdapter)(nil)

type raftHostGrpcAdapter struct {
	logger zerolog.Logger
	host   IHost
}

func (r *raftHostGrpcAdapter) Compact(ctx context.Context, request *raft.CompactRequest) (*raft.CompactReply, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return &raft.CompactReply{}, err
	}
	err := r.host.Compact(request.GetShardId(), request.GetReplicaId())
	return &raft.CompactReply{}, err
}

func (r *raftHostGrpcAdapter) GetHostConfig(ctx context.Context, request *raft.GetHostConfigRequest) (*raft.GetHostConfigReply, error) {
	hc := r.host.HostConfig()
	return &raft.GetHostConfigReply{
		Config: &raft.HostConfig{
			DeploymentId:                hc.DeploymentID,
			WalDir:                      hc.WALDir,
			HostDir:                     hc.NodeHostDir,
			RoundTripTimeInMilliseconds: hc.RTTMillisecond,
			RaftAddress:                 hc.RaftAddress,
			AddressByHostID:             hc.AddressByNodeHostID,
			ListenAddress:               hc.ListenAddress,
			MutualTls:                   hc.MutualTLS,
			CaFile:                      hc.CAFile,
			CertFile:                    hc.CertFile,
			KeyFile:                     hc.KeyFile,
			EnableMetrics:               hc.EnableMetrics,
			NotifyCommit:                hc.NotifyCommit,
		},
	}, nil
}

func (r *raftHostGrpcAdapter) Snapshot(ctx context.Context, request *raft.SnapshotRequest) (*raft.SnapshotReply, error) {
	if request.GetShardId() <= systemShardStop {
		return &raft.SnapshotReply{}, errors.New("shardId is within system shard range")
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	idx, err := r.host.Snapshot(request.GetShardId(), SnapshotOption{}, timeout)
	return &raft.SnapshotReply{
		SnapshotIndexCaptured: idx,
	}, err
}

func (r *raftHostGrpcAdapter) Stop(ctx context.Context, request *raft.StopRequest) (*raft.StopReply, error) {
	r.host.Stop()
	return &raft.StopReply{}, nil
}

func (r *raftHostGrpcAdapter) mustEmbedUnimplementedRaftHostServer() {
	//TODO implement me
	panic("implement me")
}

func (r *raftHostGrpcAdapter) checkRequestConfig(shardId, replicaId uint64) error {
	if shardId <= systemShardStop {
		return errors.New("shardId is within system shard range")
	}

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}
