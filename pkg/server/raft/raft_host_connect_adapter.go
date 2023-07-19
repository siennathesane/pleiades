/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package raft

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/cockroachdb/errors"
	"github.com/mxplusb/pleiades/pkg/raftpb"
	"github.com/mxplusb/pleiades/pkg/raftpb/raftpbconnect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	RaftConnectHostModule = fx.Module("raftpb-host-connect-adapter",
		fx.Provide(runtime.AsRoute(NewRaftHostConnectAdapter)),
	)
	_ raftpbconnect.HostServiceHandler = (*RaftHostConnectAdapter)(nil)
	_ runtime.ServiceHandler           = (*RaftHostConnectAdapter)(nil)
)

type RaftHostConnectAdapterBuilderParams struct {
	fx.In

	RaftHost runtime.IHost
	Logger   zerolog.Logger
}

type RaftHostConnectAdapterBuilderResults struct {
	fx.Out

	ConnectAdapter *RaftHostConnectAdapter
}

type RaftHostConnectAdapter struct {
	http.Handler
	logger zerolog.Logger
	host   runtime.IHost
	path   string
}

func NewRaftHostConnectAdapter(raftHost runtime.IHost, logger zerolog.Logger) *RaftHostConnectAdapter {
	if raftHost == nil {
		logger.Fatal().Err(errors.New("raftpb host is nil")).Msg("can't load connect adapter")
	}
	adapter := &RaftHostConnectAdapter{logger: logger, host: raftHost}
	adapter.path, adapter.Handler = raftpbconnect.NewHostServiceHandler(adapter)
	return adapter
}

func (r *RaftHostConnectAdapter) Path() string {
	return r.path
}

func (r *RaftHostConnectAdapter) Compact(ctx context.Context, c *connect.Request[raftpb.CompactRequest]) (*connect.Response[raftpb.CompactResponse], error) {
	if c.Msg.GetShardId() == 0 || c.Msg.GetReplicaId() == 0 {
		return connect.NewResponse(&raftpb.CompactResponse{}), errors.New("invalid shard or replica id")
	}

	err := r.host.Compact(c.Msg.GetShardId(), c.Msg.GetReplicaId())
	return connect.NewResponse(&raftpb.CompactResponse{}), err
}

func (r *RaftHostConnectAdapter) GetHostConfig(ctx context.Context, c *connect.Request[raftpb.GetHostConfigRequest]) (*connect.Response[raftpb.GetHostConfigResponse], error) {
	hc := r.host.HostConfig()
	return connect.NewResponse(&raftpb.GetHostConfigResponse{
		Config: &raftpb.HostConfig{
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

func (r *RaftHostConnectAdapter) Snapshot(ctx context.Context, c *connect.Request[raftpb.SnapshotRequest]) (*connect.Response[raftpb.SnapshotResponse], error) {
	timeout := time.Duration(c.Msg.GetTimeout()) * time.Millisecond

	idx, err := r.host.Snapshot(c.Msg.GetShardId(), runtime.SnapshotOption{}, timeout)
	return connect.NewResponse(&raftpb.SnapshotResponse{
		SnapshotIndexCaptured: idx,
	}), err
}

func (r *RaftHostConnectAdapter) Stop(ctx context.Context, c *connect.Request[raftpb.StopRequest]) (*connect.Response[raftpb.StopResponse], error) {
	r.host.Stop()
	return connect.NewResponse(&raftpb.StopResponse{}), nil
}
