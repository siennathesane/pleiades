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
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/kvpb/kvpbconnect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	_ kvpbconnect.KvStoreServiceHandler = (*KVStoreBboltConnectAdapter)(nil)
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
	adapter.path, adapter.Handler = kvpbconnect.NewKvStoreServiceHandler(adapter)

	return adapter
}

func (k *KVStoreBboltConnectAdapter) Path() string {
	return k.path
}

func (k *KVStoreBboltConnectAdapter) CreateAccount(ctx context.Context, c *connect.Request[kvpb.CreateAccountRequest]) (*connect.Response[kvpb.CreateAccountResponse], error) {
	resp, err := k.storeManager.CreateAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteAccount(ctx context.Context, c *connect.Request[kvpb.DeleteAccountRequest]) (*connect.Response[kvpb.DeleteAccountResponse], error) {
	resp, err := k.storeManager.DeleteAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) CreateBucket(ctx context.Context, c *connect.Request[kvpb.CreateBucketRequest]) (*connect.Response[kvpb.CreateBucketResponse], error) {
	resp, err := k.storeManager.CreateBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteBucket(ctx context.Context, c *connect.Request[kvpb.DeleteBucketRequest]) (*connect.Response[kvpb.DeleteBucketResponse], error) {
	resp, err := k.storeManager.DeleteBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) GetKey(ctx context.Context, c *connect.Request[kvpb.GetKeyRequest]) (*connect.Response[kvpb.GetKeyResponse], error) {
	resp, err := k.storeManager.GetKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) PutKey(ctx context.Context, c *connect.Request[kvpb.PutKeyRequest]) (*connect.Response[kvpb.PutKeyResponse], error) {
	resp, err := k.storeManager.PutKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *KVStoreBboltConnectAdapter) DeleteKey(ctx context.Context, c *connect.Request[kvpb.DeleteKeyRequest]) (*connect.Response[kvpb.DeleteKeyResponse], error) {
	resp, err := k.storeManager.DeleteKey(c.Msg)
	return connect.NewResponse(resp), err
}
