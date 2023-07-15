/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package kvstore

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	_ kvstorev1connect.KvStoreServiceHandler = (*KVStoreBboltConnectAdapter)(nil)
	_ runtime.ServiceHandler                 = (*KVStoreBboltConnectAdapter)(nil)
)

type KvStoreBboltConnectAdapterParams struct {
	fx.In

	StoreManager runtime.IKVStore
	Logger       zerolog.Logger
}

type KVStoreBboltConnectAdapter struct {
	http.Handler
	logger       zerolog.Logger
	storeManager runtime.IKVStore
	path         string
}

func NewKvStoreBboltConnectAdapter(logger zerolog.Logger, storeManager runtime.IKVStore) *KVStoreBboltConnectAdapter {
	adapter := &KVStoreBboltConnectAdapter{logger: logger, storeManager: storeManager}
	adapter.path, adapter.Handler = kvstorev1connect.NewKvStoreServiceHandler(adapter)

	return adapter
}

func (k *KVStoreBboltConnectAdapter) Path() string {
	return k.path
}

func (k *KVStoreBboltConnectAdapter) CreateAccount(ctx context.Context, c *connect.Request[kvstorev1.CreateAccountRequest]) (*connect.Response[kvstorev1.CreateAccountResponse], error) {
	resp, err := k.storeManager.CreateAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteAccount(ctx context.Context, c *connect.Request[kvstorev1.DeleteAccountRequest]) (*connect.Response[kvstorev1.DeleteAccountResponse], error) {
	resp, err := k.storeManager.DeleteAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) CreateBucket(ctx context.Context, c *connect.Request[kvstorev1.CreateBucketRequest]) (*connect.Response[kvstorev1.CreateBucketResponse], error) {
	resp, err := k.storeManager.CreateBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteBucket(ctx context.Context, c *connect.Request[kvstorev1.DeleteBucketRequest]) (*connect.Response[kvstorev1.DeleteBucketResponse], error) {
	resp, err := k.storeManager.DeleteBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) GetKey(ctx context.Context, c *connect.Request[kvstorev1.GetKeyRequest]) (*connect.Response[kvstorev1.GetKeyResponse], error) {
	resp, err := k.storeManager.GetKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) PutKey(ctx context.Context, c *connect.Request[kvstorev1.PutKeyRequest]) (*connect.Response[kvstorev1.PutKeyResponse], error) {
	resp, err := k.storeManager.PutKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteKey(ctx context.Context, c *connect.Request[kvstorev1.DeleteKeyRequest]) (*connect.Response[kvstorev1.DeleteKeyResponse], error) {
	resp, err := k.storeManager.DeleteKey(c.Msg)
	return connect.NewResponse(resp), err
}
