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

	kvstorev1 "github.com/mxplusb/api/kvstore/v1"
)

type IRaft interface {
	IShardManager
	IHost
	ITransactionManager
	IKVStore
}

type IShardManager interface {
	AddReplica(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaObserver(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	AddReplicaWitness(shardId uint64, replicaId uint64, newHost string, timeout time.Duration) error
	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
	GetShardMembers(shardId uint64) (*MembershipEntry, error)
	// LeaderTransfer(shardId uint64, targetReplicaId uint64) error
	NewShard(shardId uint64, replicaId uint64, stateMachineType StateMachineType, timeout time.Duration) error
	RemoveData(shardId, replicaId uint64) error
	RemoveReplica(shardId uint64, replicaId uint64, timeout time.Duration) error
	StartReplica(shardId uint64, replicaId uint64, stateMachineType StateMachineType) error
	StartReplicaObserver(shardId uint64, replicaId uint64, stateMachineType StateMachineType) error
	StopReplica(shardId uint64) (*OperationResult, error)
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
