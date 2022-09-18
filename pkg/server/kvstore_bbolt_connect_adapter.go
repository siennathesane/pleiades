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

	kvstorev1 "github.com/mxplusb/api/kvstore/v1"
	"github.com/mxplusb/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

var _ kvstorev1connect.KvStoreServiceHandler = (*kvstoreBboltConnectAdapter)(nil)

type kvstoreBboltConnectAdapter struct {
	logger       zerolog.Logger
	storeManager IKVStore
}

func (k *kvstoreBboltConnectAdapter) CreateAccount(ctx context.Context, c *connect.Request[kvstorev1.CreateAccountRequest]) (*connect.Response[kvstorev1.CreateAccountResponse], error) {
	resp, err := k.storeManager.CreateAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) DeleteAccount(ctx context.Context, c *connect.Request[kvstorev1.DeleteAccountRequest]) (*connect.Response[kvstorev1.DeleteAccountResponse], error) {
	resp, err := k.storeManager.DeleteAccount(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) CreateBucket(ctx context.Context, c *connect.Request[kvstorev1.CreateBucketRequest]) (*connect.Response[kvstorev1.CreateBucketResponse], error) {
	resp, err := k.storeManager.CreateBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) DeleteBucket(ctx context.Context, c *connect.Request[kvstorev1.DeleteBucketRequest]) (*connect.Response[kvstorev1.DeleteBucketResponse], error) {
	resp, err := k.storeManager.DeleteBucket(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) GetKey(ctx context.Context, c *connect.Request[kvstorev1.GetKeyRequest]) (*connect.Response[kvstorev1.GetKeyResponse], error) {
	resp, err := k.storeManager.GetKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) PutKey(ctx context.Context, c *connect.Request[kvstorev1.PutKeyRequest]) (*connect.Response[kvstorev1.PutKeyResponse], error) {
	resp, err := k.storeManager.PutKey(c.Msg)
	return connect.NewResponse(resp), err
}

func (k *kvstoreBboltConnectAdapter) DeleteKey(ctx context.Context, c *connect.Request[kvstorev1.DeleteKeyRequest]) (*connect.Response[kvstorev1.DeleteKeyResponse], error) {
	resp, err := k.storeManager.DeleteKey(c.Msg)
	return connect.NewResponse(resp), err
}
