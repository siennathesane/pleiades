/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package kvstore

import (
	"context"

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

var _ kvstorev1connect.KvStoreServiceHandler = (*KVstoreBboltConnectAdapter)(nil)

type KVstoreBboltConnectAdapter struct {
	logger       zerolog.Logger
	storeManager runtime.IKVStore
}

func NewKvstoreBboltConnectAdapter(storeManager runtime.IKVStore, logger zerolog.Logger) *KVstoreBboltConnectAdapter {
	return &KVstoreBboltConnectAdapter{logger: logger, storeManager: storeManager}
}

func (k *KVstoreBboltConnectAdapter) CreateAccount(ctx context.Context, c *connect.Request[kvstorev1.CreateAccountRequest]) (*connect.Response[kvstorev1.CreateAccountResponse], error) {
	resp, err := k.storeManager.CreateAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) DeleteAccount(ctx context.Context, c *connect.Request[kvstorev1.DeleteAccountRequest]) (*connect.Response[kvstorev1.DeleteAccountResponse], error) {
	resp, err := k.storeManager.DeleteAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) CreateBucket(ctx context.Context, c *connect.Request[kvstorev1.CreateBucketRequest]) (*connect.Response[kvstorev1.CreateBucketResponse], error) {
	resp, err := k.storeManager.CreateBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) DeleteBucket(ctx context.Context, c *connect.Request[kvstorev1.DeleteBucketRequest]) (*connect.Response[kvstorev1.DeleteBucketResponse], error) {
	resp, err := k.storeManager.DeleteBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) GetKey(ctx context.Context, c *connect.Request[kvstorev1.GetKeyRequest]) (*connect.Response[kvstorev1.GetKeyResponse], error) {
	resp, err := k.storeManager.GetKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) PutKey(ctx context.Context, c *connect.Request[kvstorev1.PutKeyRequest]) (*connect.Response[kvstorev1.PutKeyResponse], error) {
	resp, err := k.storeManager.PutKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVstoreBboltConnectAdapter) DeleteKey(ctx context.Context, c *connect.Request[kvstorev1.DeleteKeyRequest]) (*connect.Response[kvstorev1.DeleteKeyResponse], error) {
	resp, err := k.storeManager.DeleteKey(c.Msg)
	return connect.NewResponse(resp), err
}
