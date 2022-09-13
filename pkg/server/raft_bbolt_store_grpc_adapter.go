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

	"github.com/mxplusb/pleiades/api/v1/database"
	"github.com/rs/zerolog"
)

var _ KVStoreServiceServer = (*raftBBoltStoreManagerGrpcAdapter)(nil)

type raftBBoltStoreManagerGrpcAdapter struct {
	logger       zerolog.Logger
	storeManager IKVStore
}

func (r *raftBBoltStoreManagerGrpcAdapter) CreateAccount(ctx context.Context, request *database.CreateAccountRequest) (*database.CreateAccountReply, error) {
	return r.storeManager.CreateAccount(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteAccount(ctx context.Context, request *database.DeleteAccountRequest) (*database.DeleteAccountReply, error) {
	return r.storeManager.DeleteAccount(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) CreateBucket(ctx context.Context, request *database.CreateBucketRequest) (*database.CreateBucketReply, error) {
	return r.storeManager.CreateBucket(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteBucket(ctx context.Context, request *database.DeleteBucketRequest) (*database.DeleteBucketReply, error) {
	return r.storeManager.DeleteBucket(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) GetKey(ctx context.Context, request *database.GetKeyRequest) (*database.GetKeyReply, error) {
	return r.storeManager.GetKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) PutKey(ctx context.Context, request *database.PutKeyRequest) (*database.PutKeyReply, error) {
	return r.storeManager.PutKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteKey(ctx context.Context, request *database.DeleteKeyRequest) (*database.DeleteKeyReply, error) {
	return r.storeManager.DeleteKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) mustEmbedUnimplementedKVStoreServiceServer() {
	//TODO implement me
	panic("implement me")
}
