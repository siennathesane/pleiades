/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

import (
	"context"
	"net/http"
	"time"

	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/raftpb"
	"go.uber.org/fx"
)

type IRaft interface {
	IShardManager
	IHost
	ITransactionManager
	IKVStore
}

type ServiceHandler interface {
	http.Handler
	Path() string
}

type IShardManager interface {
	AddReplica(req *raftpb.AddReplicaRequest) error
	AddReplicaObserver(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaWitness(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
	GetShardMembers(shardId uint64) (*MembershipEntry, error)
	// LeaderTransfer(shardId uint64, targetReplicaId uint64) error
	NewShard(req *raftpb.NewShardRequest) error
	RemoveData(shardId, replicaId uint64) error
	RemoveReplica(shardId uint64, replicaId uint64, timeout time.Duration) error
	StartReplica(req *raftpb.StartReplicaRequest) error
	StartReplicaObserver(req *raftpb.StartReplicaObserverRequest) error
	StopReplica(shardId uint64, replicaId uint64) (*OperationResult, error)
}

type IHost interface {
	Compact(shardId uint64, replicaId uint64) error
	GetHostInfo(opt HostInfoOption) *HostInfo
	HasNodeInfo(shardId uint64, replicaId uint64) bool
	Id() string
	HostConfig() HostConfig
	RaftAddress() string
	Snapshot(shardId uint64, opt SnapshotOption, timeout time.Duration) (uint64, error)
	Stop()
}

type ITransactionManager interface {
	CloseTransaction(ctx context.Context, transaction *kvpb.Transaction) error
	Commit(ctx context.Context, transaction *kvpb.Transaction) *kvpb.Transaction
	GetNoOpTransaction(shardId uint64) *kvpb.Transaction
	GetTransaction(ctx context.Context, shardId uint64) (*kvpb.Transaction, error)
	SessionFromClientId(clientId uint64) (*dclient.Session, bool)
}

type IKVStore interface {
	CreateAccount(request *kvpb.CreateAccountRequest) (*kvpb.CreateAccountResponse, error)
	DeleteAccount(request *kvpb.DeleteAccountRequest) (*kvpb.DeleteAccountResponse, error)
	CreateBucket(request *kvpb.CreateBucketRequest) (*kvpb.CreateBucketResponse, error)
	DeleteBucket(request *kvpb.DeleteBucketRequest) (*kvpb.DeleteBucketResponse, error)
	GetKey(request *kvpb.GetKeyRequest) (*kvpb.GetKeyResponse, error)
	PutKey(request *kvpb.PutKeyRequest) (*kvpb.PutKeyResponse, error)
	DeleteKey(request *kvpb.DeleteKeyRequest) (*kvpb.DeleteKeyResponse, error)
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ServiceHandler)),
		fx.ResultTags(`group:"routes"`),
	)
}
