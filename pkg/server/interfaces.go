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

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
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
	CloseTransaction(ctx context.Context, transaction *database.Transaction) error
	Commit(ctx context.Context, transaction *database.Transaction) *database.Transaction
	GetNoOpTransaction(shardId uint64) *database.Transaction
	GetTransaction(ctx context.Context, shardId uint64) (*database.Transaction, error)
}

type IKVStore interface {
	CreateAccount(request *database.CreateAccountRequest) (*database.CreateAccountReply, error)
	DeleteAccount(request *database.DeleteAccountRequest) (*database.DeleteAccountReply, error)
	CreateBucket(request *database.CreateBucketRequest) (*database.CreateBucketReply, error)
	DeleteBucket(request *database.DeleteBucketRequest) (*database.DeleteBucketReply, error)
	GetKey(request *database.GetKeyRequest) (*database.GetKeyReply, error)
	PutKey(request *database.PutKeyRequest) (*database.PutKeyReply, error)
	DeleteKey(request *database.DeleteKeyRequest) (*database.DeleteKeyReply, error)
}
