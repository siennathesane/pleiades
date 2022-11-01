/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package raft

import (
	"context"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/bufbuild/connect-go"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ raftv1connect.HostServiceHandler = (*RaftHostConnectAdapter)(nil)

type RaftHostConnectAdapter struct {
	logger zerolog.Logger
	host   runtime.IHost
}

func NewRaftHostConnectAdapter(host runtime.IHost, logger zerolog.Logger) *RaftHostConnectAdapter {
	return &RaftHostConnectAdapter{logger: logger, host: host}
}

func (r *RaftHostConnectAdapter) Compact(ctx context.Context, c *connect.Request[raftv1.CompactRequest]) (*connect.Response[raftv1.CompactResponse], error) {
	if c.Msg.GetShardId() == 0 || c.Msg.GetReplicaId() == 0 {
		return connect.NewResponse(&raftv1.CompactResponse{}), errors.New("invalid shard or replica id")
	}

	err := r.host.Compact(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	return connect.NewResponse(&raftv1.CompactResponse{}), err
}

func (r *RaftHostConnectAdapter) GetHostConfig(ctx context.Context, c *connect.Request[raftv1.GetHostConfigRequest]) (*connect.Response[raftv1.GetHostConfigResponse], error) {
	hc := r.host.HostConfig()
	return connect.NewResponse(&raftv1.GetHostConfigResponse{
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
	}), nil
}

func (r *RaftHostConnectAdapter) Snapshot(ctx context.Context, c *connect.Request[raftv1.SnapshotRequest]) (*connect.Response[raftv1.SnapshotResponse], error) {
	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	idx, err := r.host.Snapshot(c.Msg.GetShardId(), runtime.SnapshotOption{}, timeout)
	return connect.NewResponse(&raftv1.SnapshotResponse{
		SnapshotIndexCaptured: idx,
	}), err
}

func (r *RaftHostConnectAdapter) Stop(ctx context.Context, c *connect.Request[raftv1.StopRequest]) (*connect.Response[raftv1.StopResponse], error) {
	r.host.Stop()
	return connect.NewResponse(&raftv1.StopResponse{}), nil
}
