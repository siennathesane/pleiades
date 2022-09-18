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
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ raftv1.HostServiceServer = (*raftHostGrpcAdapter)(nil)

type raftHostGrpcAdapter struct {
	logger zerolog.Logger
	host   IHost
}

func (r *raftHostGrpcAdapter) Compact(ctx context.Context, request *raftv1.CompactRequest) (*raftv1.CompactResponse, error) {
	if err := r.checkRequestConfig(request.GetShardId(), request.GetReplicaId()); err != nil {
		return &raftv1.CompactResponse{}, err
	}
	err := r.host.Compact(request.GetShardId(), request.GetReplicaId())
	return &raftv1.CompactResponse{}, err
}

func (r *raftHostGrpcAdapter) GetHostConfig(ctx context.Context, request *raftv1.GetHostConfigRequest) (*raftv1.GetHostConfigResponse, error) {
	hc := r.host.HostConfig()
	return &raftv1.GetHostConfigResponse{
		Config: &raftv1.HostConfig{
			DeploymentId:                hc.DeploymentID,
			WalDir:                      hc.WALDir,
			HostDir:                     hc.NodeHostDir,
			RoundTripTimeInMilliseconds: hc.RTTMillisecond,
			RaftAddress:                 hc.RaftAddress,
			AddressByHostId:             hc.AddressByNodeHostID,
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

func (r *raftHostGrpcAdapter) Snapshot(ctx context.Context, request *raftv1.SnapshotRequest) (*raftv1.SnapshotResponse, error) {
	if request.GetShardId() <= systemShardStop {
		return &raftv1.SnapshotResponse{}, ErrSystemShardRange
	}

	timeout := time.Duration(request.Timeout) * time.Millisecond

	idx, err := r.host.Snapshot(request.GetShardId(), SnapshotOption{}, timeout)
	return &raftv1.SnapshotResponse{
		SnapshotIndexCaptured: idx,
	}, err
}

func (r *raftHostGrpcAdapter) Stop(ctx context.Context, request *raftv1.StopRequest) (*raftv1.StopResponse, error) {
	r.host.Stop()
	return &raftv1.StopResponse{}, nil
}

func (r *raftHostGrpcAdapter) mustEmbedUnimplementedRaftHostServer() {
	//TODO implement me
	panic("implement me")
}

func (r *raftHostGrpcAdapter) checkRequestConfig(shardId, replicaId uint64) error {
	if shardId <= systemShardStop {
		return ErrSystemShardRange
	}

	if replicaId == 0 {
		return errors.New("replicaId cannot be zero")
	}

	return nil
}
