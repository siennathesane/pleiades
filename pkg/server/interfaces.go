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
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/statemachine"
)

type IShardManager interface {
	AddReplica(cfg IClusterConfig, newHost string, timeout time.Duration) error
	AddShardObserver(cfg IClusterConfig, newHost string, timeout time.Duration) error
	AddShardWitness(cfg IClusterConfig, newHost string, timeout time.Duration) error
	DeleteReplica(cfg IClusterConfig, timeout time.Duration) error
	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
	GetShardMembers(shardId uint64) (*MembershipEntry, error)
	NewShard(cfg IClusterConfig) error
	RemoveData(shardId, replicaId uint64) error
	StopReplica(shardId uint64) (*OperationResult, error)
}

//type INodeConfig interface {
//	NodeHostConfig() config.NodeHostConfig
//	HasNodeInfo(clusterID uint64, nodeID uint64) bool
//	GetNodeHostInfo(opt dragonboat.NodeHostInfoOption) *dragonboat.NodeHostInfo
//}

type INodeHost interface {
	Compact(clusterID uint64, nodeID uint64) (*dragonboat.SysOpState, error)
	GetNodeUser(clusterID uint64) (dragonboat.INodeUser, error)
	ID() string
	LeaderTransfer(clusterID uint64, targetNodeID uint64) error
	RaftAddress() string
	Snapshot(clusterID uint64, opt dragonboat.SnapshotOption, timeout time.Duration) (*dragonboat.RequestState, error)
	Stop()
	StopNode(clusterID uint64, nodeID uint64) error
}

type ITransactionManager interface {
	GetNoOpSession(clusterID uint64) *client.Session
	CloseSession(ctx context.Context, cs *client.Session) error
	GetSession(ctx context.Context, clusterID uint64) (*client.Session, error)
}

type IStore interface {
	Propose(session *client.Session, cmd []byte, timeout time.Duration) (*dragonboat.RequestState, error)
	SyncPropose(ctx context.Context, session *client.Session, cmd []byte) (statemachine.Result, error)
	SyncRead(ctx context.Context, clusterID uint64, query interface{}) (interface{}, error)
}
