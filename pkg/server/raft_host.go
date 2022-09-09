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

	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
)

var _ IHost = (*raftHost)(nil)

func newRaftHost(host *dragonboat.NodeHost, logger zerolog.Logger) *raftHost {
	return &raftHost{
		logger: logger.With().Str("component", "raft-host").Logger(),
		nh:     host,
	}
}

type raftHost struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (r *raftHost) Compact(shardId uint64, replicaId uint64) error {
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

func (r *raftHost) GetHostInfo(opt HostInfoOption) *HostInfo {
	nhi := r.nh.GetNodeHostInfo(dragonboat.NodeHostInfoOption{
		SkipLogInfo: opt.SkipLogInfo,
	})

	cil := make([]ClusterInfo, len(nhi.ClusterInfoList))
	for idx := range nhi.ClusterInfoList {
		cil[idx] = ClusterInfo{
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

	li := make([]NodeInfo, len(nhi.LogInfo))
	for idx := range nhi.LogInfo {
		li[idx] = NodeInfo{
			ShardId:   nhi.LogInfo[idx].ClusterID,
			ReplicaId: nhi.LogInfo[idx].NodeID,
		}
	}

	return &HostInfo{
		HostId:          nhi.NodeHostID,
		RaftAddress:     nhi.RaftAddress,
		Gossip:          GossipInfo{
			Enabled: nhi.Gossip.Enabled,
			AdvertiseAddress: nhi.Gossip.AdvertiseAddress,
			NumOfLiveKnownHosts: nhi.Gossip.NumOfKnownNodeHosts,
		},
		ClusterInfoList: cil,
		LogInfo:         li,
	}
}

func (r *raftHost) HasNodeInfo(shardId uint64, replicaId uint64) bool {
	return r.nh.HasNodeInfo(shardId, replicaId)
}

func (r *raftHost) Id() string {
	return r.nh.ID()
}

func (r *raftHost) HostConfig() HostConfig {
	nhc := r.nh.NodeHostConfig()
	return HostConfig{
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

func (r *raftHost) RaftAddress() string {
	return r.nh.RaftAddress()
}

func (r *raftHost) Snapshot(shardId uint64, opt SnapshotOption, timeout time.Duration) (uint64, error) {
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

func (r *raftHost) Stop() {
	r.logger.Info().Msg("stopping raft host")
	r.nh.Stop()
}

type HostInfo struct {
	HostId          string
	RaftAddress     string
	Gossip          GossipInfo
	ClusterInfoList []ClusterInfo
	LogInfo         []NodeInfo
}

type GossipInfo struct {
	Enabled             bool
	AdvertiseAddress    string
	NumOfLiveKnownHosts int
}

type ClusterInfo struct {
	ShardId           uint64
	ReplicaId         uint64
	Nodes             map[uint64]string
	ConfigChangeIndex uint64
	IsLeader          bool
	IsObserver        bool
	IsWitness         bool
	Pending           bool
}

type NodeInfo struct {
	ShardId   uint64
	ReplicaId uint64
}

type HostConfig struct {
	DeploymentID        uint64
	WALDir              string
	NodeHostDir         string
	RTTMillisecond      uint64
	RaftAddress         string
	AddressByNodeHostID bool
	ListenAddress       string
	MutualTLS           bool
	CAFile              string
	CertFile            string
	KeyFile             string
	EnableMetrics       bool
	NotifyCommit        bool
}

type HostInfoOption struct {
	SkipLogInfo bool
}

type SnapshotOption struct {
	CompactionOverhead         uint64
	ExportPath                 string
	Exported                   bool
	OverrideCompactionOverhead bool
}
