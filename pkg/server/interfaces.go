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
	AddReplica(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaObserver(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaWitness(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
	GetShardMembers(shardId uint64) (*MembershipEntry, error)
	NewShard(shardId uint64, replicaId uint64, stateMachineType StateMachineType, timeout time.Duration) error
	RemoveData(shardId, replicaId uint64) error
	RemoveReplica(shardId uint64, replicaId uint64, timeout time.Duration) error
	StartReplica(shardId uint64, replicaId uint64, stateMachineType StateMachineType) error
	StartReplicaObserver(shardId uint64, replicaId uint64, stateMachineType StateMachineType) error
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
