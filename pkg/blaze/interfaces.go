/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"context"
	"time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/statemachine"
)

type ICluster interface {
	NAReadLocalNode(rs *dragonboat.RequestState, query []byte) ([]byte, error)
	ReadIndex(clusterID uint64, timeout time.Duration) (*dragonboat.RequestState, error)
	ReadLocalNode(rs *dragonboat.RequestState, query interface{}) (interface{}, error)
	StaleRead(clusterID uint64, query interface{}) (interface{}, error)
	StartCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateStateMachineFunc, cfg config.Config) error
	StartConcurrentCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateConcurrentStateMachineFunc, cfg config.Config) error
	StartOnDiskCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateOnDiskStateMachineFunc, cfg config.Config) error
	StopCluster(clusterID uint64) error
	SyncGetClusterMembership(ctx context.Context, clusterID uint64) (*dragonboat.Membership, error)
}

type INodeConfig interface {
	NodeHostConfig() config.NodeHostConfig
	HasNodeInfo(clusterID uint64, nodeID uint64) bool
	GetNodeHostInfo(opt dragonboat.NodeHostInfoOption) *dragonboat.NodeHostInfo
}

type INodeHost interface {
	GetLeaderID(clusterID uint64) (uint64, bool, error)
	GetNodeUser(clusterID uint64) (dragonboat.INodeUser, error)
	ID() string
	RaftAddress() string
	RemoveData(clusterID uint64, nodeID uint64) error
	RequestAddNode(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error)
	RequestAddObserver(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error)
	RequestAddWitness(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error)
	RequestCompaction(clusterID uint64, nodeID uint64) (*dragonboat.SysOpState, error)
	RequestDeleteNode(clusterID uint64, nodeID uint64, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error)
	RequestLeaderTransfer(clusterID uint64, targetNodeID uint64) error
	RequestSnapshot(clusterID uint64, opt dragonboat.SnapshotOption, timeout time.Duration) (*dragonboat.RequestState, error)
	Stop()
	StopNode(clusterID uint64, nodeID uint64) error
	// todo (sienna): add these back later if we need them
	//SyncRemoveData(ctx context.Context, clusterID uint64, nodeID uint64) error
	//SyncRequestAddNode(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error
	//SyncRequestAddObserver(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error
	//SyncRequestAddWitness(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error
	//SyncRequestDeleteNode(ctx context.Context, clusterID uint64, nodeID uint64, configChangeIndex uint64) error
	//SyncRequestSnapshot(ctx context.Context, clusterID uint64, opt dragonboat.SnapshotOption) (uint64, error)
}

type ISession interface {
	GetNoOPSession(clusterID uint64) *client.Session
	ProposeSession(session *client.Session, timeout time.Duration) (*dragonboat.RequestState, error)
	SyncCloseSession(ctx context.Context, cs *client.Session) error
	SyncGetSession(ctx context.Context, clusterID uint64) (*client.Session, error)
}

type IStore interface {
	Propose(session *client.Session, cmd []byte, timeout time.Duration) (*dragonboat.RequestState, error)
	SyncPropose(ctx context.Context, session *client.Session, cmd []byte) (statemachine.Result, error)
	SyncRead(ctx context.Context, clusterID uint64, query interface{}) (interface{}, error)
}
