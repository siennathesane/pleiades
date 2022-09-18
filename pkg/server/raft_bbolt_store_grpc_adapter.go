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

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/rs/zerolog"
)

var _ kvstorev1.KvStoreServiceServer = (*raftBBoltStoreManagerGrpcAdapter)(nil)

type raftBBoltStoreManagerGrpcAdapter struct {
	logger       zerolog.Logger
	storeManager IKVStore
}

func (r *raftBBoltStoreManagerGrpcAdapter) CreateAccount(ctx context.Context, request *kvstorev1.CreateAccountRequest) (*kvstorev1.CreateAccountResponse, error) {
	return r.storeManager.CreateAccount(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteAccount(ctx context.Context, request *kvstorev1.DeleteAccountRequest) (*kvstorev1.DeleteAccountResponse, error) {
	return r.storeManager.DeleteAccount(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) CreateBucket(ctx context.Context, request *kvstorev1.CreateBucketRequest) (*kvstorev1.CreateBucketResponse, error) {
	return r.storeManager.CreateBucket(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteBucket(ctx context.Context, request *kvstorev1.DeleteBucketRequest) (*kvstorev1.DeleteBucketResponse, error) {
	return r.storeManager.DeleteBucket(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) GetKey(ctx context.Context, request *kvstorev1.GetKeyRequest) (*kvstorev1.GetKeyResponse, error) {
	return r.storeManager.GetKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) PutKey(ctx context.Context, request *kvstorev1.PutKeyRequest) (*kvstorev1.PutKeyResponse, error) {
	return r.storeManager.PutKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) DeleteKey(ctx context.Context, request *kvstorev1.DeleteKeyRequest) (*kvstorev1.DeleteKeyResponse, error) {
	return r.storeManager.DeleteKey(request)
}

func (r *raftBBoltStoreManagerGrpcAdapter) mustEmbedUnimplementedKVStoreServiceServer() {
	//TODO implement me
	panic("implement me")
}
