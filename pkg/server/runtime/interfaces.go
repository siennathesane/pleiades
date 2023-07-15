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
	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
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
	AddReplica(req *raftv1.AddReplicaRequest) error
	AddReplicaObserver(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaWitness(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
	GetShardMembers(shardId uint64) (*MembershipEntry, error)
	// LeaderTransfer(shardId uint64, targetReplicaId uint64) error
	NewShard(req *raftv1.NewShardRequest) error
	RemoveData(shardId, replicaId uint64) error
	RemoveReplica(shardId uint64, replicaId uint64, timeout time.Duration) error
	StartReplica(req *raftv1.StartReplicaRequest) error
	StartReplicaObserver(req *raftv1.StartReplicaObserverRequest) error
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
	CloseTransaction(ctx context.Context, transaction *kvstorev1.Transaction) error
	Commit(ctx context.Context, transaction *kvstorev1.Transaction) *kvstorev1.Transaction
	GetNoOpTransaction(shardId uint64) *kvstorev1.Transaction
	GetTransaction(ctx context.Context, shardId uint64) (*kvstorev1.Transaction, error)
	SessionFromClientId(clientId uint64) (*dclient.Session, bool)
}

type IKVStore interface {
	CreateAccount(request *kvstorev1.CreateAccountRequest) (*kvstorev1.CreateAccountResponse, error)
	DeleteAccount(request *kvstorev1.DeleteAccountRequest) (*kvstorev1.DeleteAccountResponse, error)
	CreateBucket(request *kvstorev1.CreateBucketRequest) (*kvstorev1.CreateBucketResponse, error)
	DeleteBucket(request *kvstorev1.DeleteBucketRequest) (*kvstorev1.DeleteBucketResponse, error)
	GetKey(request *kvstorev1.GetKeyRequest) (*kvstorev1.GetKeyResponse, error)
	PutKey(request *kvstorev1.PutKeyRequest) (*kvstorev1.PutKeyResponse, error)
	DeleteKey(request *kvstorev1.DeleteKeyRequest) (*kvstorev1.DeleteKeyResponse, error)
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
