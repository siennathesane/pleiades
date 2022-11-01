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

	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
)

var _ runtime.IHost = (*RaftHost)(nil)

func NewRaftHost(host *dragonboat.NodeHost, logger zerolog.Logger) *RaftHost {
	return &RaftHost{
		logger: logger.With().Str("component", "raft-host").Logger(),
		nh:     host,
	}
}

type RaftHost struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (r *RaftHost) Compact(shardId uint64, replicaId uint64) error {
	l := r.logger.With().Uint64("shard", shardId).Uint64("replica", replicaId).Logger()

	l.Info().Msg("requesting compaction")
	state, err := r.nh.RequestCompaction(shardId, replicaId)
	if err != nil {
		l.Error().Err(err).Msg("cannot request compaction")
		return err
	}

	l.Debug().Msg("awaiting compaction")
	<-state.ResultC()

	l.Debug().Msg("compaction completed")
	return nil
}

func (r *RaftHost) GetHostInfo(opt runtime.HostInfoOption) *runtime.HostInfo {
	nhi := r.nh.GetNodeHostInfo(dragonboat.NodeHostInfoOption{
		SkipLogInfo: opt.SkipLogInfo,
	})

	cil := make([]runtime.ClusterInfo, len(nhi.ClusterInfoList))
	for idx := range nhi.ClusterInfoList {
		cil[idx] = runtime.ClusterInfo{
			ShardId:           nhi.ClusterInfoList[idx].ClusterID,
			ReplicaId:         nhi.ClusterInfoList[idx].NodeID,
			Nodes:             nhi.ClusterInfoList[idx].Nodes,
			ConfigChangeIndex: nhi.ClusterInfoList[idx].ConfigChangeIndex,
			IsLeader:          nhi.ClusterInfoList[idx].IsLeader,
			IsObserver:        nhi.ClusterInfoList[idx].IsObserver,
			IsWitness:         nhi.ClusterInfoList[idx].IsWitness,
			Pending:           nhi.ClusterInfoList[idx].Pending,
		}
	}

	li := make([]runtime.NodeInfo, len(nhi.LogInfo))
	for idx := range nhi.LogInfo {
		li[idx] = runtime.NodeInfo{
			ShardId:   nhi.LogInfo[idx].ClusterID,
			ReplicaId: nhi.LogInfo[idx].NodeID,
		}
	}

	return &runtime.HostInfo{
		HostId:      nhi.NodeHostID,
		RaftAddress: nhi.RaftAddress,
		Gossip: runtime.GossipInfo{
			Enabled:             nhi.Gossip.Enabled,
			AdvertiseAddress:    nhi.Gossip.AdvertiseAddress,
			NumOfLiveKnownHosts: nhi.Gossip.NumOfKnownNodeHosts,
		},
		ClusterInfoList: cil,
		LogInfo:         li,
	}
}

func (r *RaftHost) HasNodeInfo(shardId uint64, replicaId uint64) bool {
	return r.nh.HasNodeInfo(shardId, replicaId)
}

func (r *RaftHost) Id() string {
	return r.nh.ID()
}

func (r *RaftHost) HostConfig() runtime.HostConfig {
	nhc := r.nh.NodeHostConfig()
	return runtime.HostConfig{
		nhc.DeploymentID,
		nhc.WALDir,
		nhc.NodeHostDir,
		nhc.RTTMillisecond,
		nhc.RaftAddress,
		nhc.AddressByNodeHostID,
		nhc.ListenAddress,
		nhc.MutualTLS,
		nhc.CAFile,
		nhc.CertFile,
		nhc.KeyFile,
		nhc.EnableMetrics,
		nhc.NotifyCommit,
	}
}

func (r *RaftHost) RaftAddress() string {
	return r.nh.RaftAddress()
}

func (r *RaftHost) Snapshot(shardId uint64, opt runtime.SnapshotOption, timeout time.Duration) (uint64, error) {
	l := r.logger.With().Uint64("shard", shardId).Logger()
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	l.Info().Msg("requesting snapshot")

	res, err := r.nh.SyncRequestSnapshot(ctx, shardId, dragonboat.SnapshotOption{
		CompactionOverhead:         opt.CompactionOverhead,
		ExportPath:                 opt.ExportPath,
		Exported:                   opt.Exported,
		OverrideCompactionOverhead: opt.OverrideCompactionOverhead,
	})
	if err != nil {
		l.Error().Err(err).Msg("cannot snapshot shard")
		return 0, err
	}

	return res, nil
}

func (r *RaftHost) Stop() {
	r.logger.Info().Msg("stopping raft host")
	r.nh.Stop()
}
