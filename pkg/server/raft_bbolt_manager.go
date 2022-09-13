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
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/mxplusb/pleiades/pkg/routing"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var (
	_ IStore = (*bboltStoreManager)(nil)
)

func newBboltStoreManager(tm *raftTransactionManager, nh *dragonboat.NodeHost, logger zerolog.Logger) *bboltStoreManager {
	l := logger.With().Str("component", "store-manager").Logger()
	return &bboltStoreManager{l, tm, nh, &routing.Shard{}, 1000 * time.Millisecond}
}

type bboltStoreManager struct {
	logger         zerolog.Logger
	tm             *raftTransactionManager
	nh             *dragonboat.NodeHost
	shardRouter    *routing.Shard
	defaultTimeout time.Duration
}

func (s *bboltStoreManager) CreateAccount(request *database.CreateAccountRequest) (*database.CreateAccountReply, error) {

	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.CreateAccountReply{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &database.CreateAccountReply{}, kv.ErrInvalidOwner
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Typ:     database.KVStoreWrapper_CREATE_ACCOUNT_REQUEST,
		Payload: &database.KVStoreWrapper_CreateAccountRequest{
			CreateAccountRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.CreateAccountReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.CreateAccountReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.CreateAccountReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.CreateAccountReply{}
	if resp.GetCreateAccountReply() == nil {
		response.AccountDescriptor = &database.AccountDescriptor{}
	} else {
		response.AccountDescriptor = resp.GetCreateAccountReply().GetAccountDescriptor()
	}

	if request.Transaction == nil {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *bboltStoreManager) DeleteAccount(request *database.DeleteAccountRequest) (*database.DeleteAccountReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.DeleteAccountReply{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &database.DeleteAccountReply{}, kv.ErrInvalidOwner
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Typ:     database.KVStoreWrapper_DELETE_ACCOUNT_REQUEST,
		Payload: &database.KVStoreWrapper_DeleteAccountRequest{
			DeleteAccountRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.DeleteAccountReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.DeleteAccountReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.DeleteAccountReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.DeleteAccountReply{
		Ok: resp.GetDeleteAccountReply().GetOk(),
	}

	if request.Transaction == nil || cs.SeriesID != client.NoOPSeriesID {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *bboltStoreManager) CreateBucket(request *database.CreateBucketRequest) (*database.CreateBucketReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.CreateBucketReply{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &database.CreateBucketReply{}, kv.ErrInvalidOwner
	}

	bucketName := request.GetName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &database.CreateBucketReply{}, kv.ErrInvalidBucketName
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     database.KVStoreWrapper_CREATE_BUCKET_REQUEST,
		Payload: &database.KVStoreWrapper_CreateBucketRequest{
			CreateBucketRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.CreateBucketReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.CreateBucketReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.CreateBucketReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.CreateBucketReply{}
	if resp.GetCreateBucketReply() == nil {
		response.BucketDescriptor = &database.BucketDescriptor{}
	} else {
		response.BucketDescriptor = resp.GetCreateBucketReply().GetBucketDescriptor()
	}

	if request.Transaction == nil {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *bboltStoreManager) DeleteBucket(request *database.DeleteBucketRequest) (*database.DeleteBucketReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.DeleteBucketReply{}, kv.ErrInvalidAccount
	}

	name := request.GetName()
	if name == "" {
		s.logger.Trace().Msg("empty name value")
		return &database.DeleteBucketReply{}, kv.ErrInvalidOwner
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  name,
		Typ:     database.KVStoreWrapper_DELETE_BUCKET_REQUEST,
		Payload: &database.KVStoreWrapper_DeleteBucketRequest{
			DeleteBucketRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.DeleteBucketReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.DeleteBucketReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.DeleteBucketReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.DeleteBucketReply{
		Ok: resp.GetDeleteBucketReply().GetOk(),
	}

	if request.Transaction == nil || cs.SeriesID != client.NoOPSeriesID {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *bboltStoreManager) GetKey(request *database.GetKeyRequest) (*database.GetKeyReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.GetKeyReply{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &database.GetKeyReply{}, kv.ErrInvalidOwner
	}

	keyName := request.GetKey()
	if keyName == "" {
		s.logger.Trace().Msg("empty key name")
		return &database.GetKeyReply{}, errors.New("empty key name")
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     database.KVStoreWrapper_GET_KEY_REQUEST,
		Payload: &database.KVStoreWrapper_GetKeyRequest{
			GetKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.GetKeyReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	shard := s.shardRouter.AccountToShard(account)

	result, err := s.nh.SyncRead(ctx, shard, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.GetKeyReply{}, errors.Wrap(err, "can't apply message")
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.([]byte))
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.GetKeyReply{}, errors.Wrap(err, "can't unmarshal response")
	}


	response := &database.GetKeyReply{}
	kvp := resp.GetGetKeyReply()
	if kvp == nil {
		response.KeyValuePair = &database.KeyValue{}
	} else {
		response.KeyValuePair = kvp.GetKeyValuePair()
	}

	return response, nil
}

func (s *bboltStoreManager) PutKey(request *database.PutKeyRequest) (*database.PutKeyReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.PutKeyReply{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucketName value")
		return &database.PutKeyReply{}, kv.ErrInvalidOwner
	}

	keyValuePair := request.GetKeyValuePair()
	if keyValuePair == nil {
		s.logger.Trace().Msg("empty key value pair")
		return &database.PutKeyReply{}, kv.ErrInvalidBucketName
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     database.KVStoreWrapper_PUT_KEY_REQUEST,
		Payload: &database.KVStoreWrapper_PutKeyRequest{
			PutKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.PutKeyReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.PutKeyReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.PutKeyReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.PutKeyReply{}

	if request.Transaction == nil || cs.SeriesID != client.NoOPSeriesID {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *bboltStoreManager) DeleteKey(request *database.DeleteKeyRequest) (*database.DeleteKeyReply, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &database.DeleteKeyReply{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &database.DeleteKeyReply{}, kv.ErrInvalidOwner
	}

	key := request.GetKey()
	if key == "" {
		s.logger.Trace().Msg("empty key value pair")
		return &database.DeleteKeyReply{}, kv.ErrInvalidBucketName
	}

	req := &database.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     database.KVStoreWrapper_DELETE_KEY_REQUEST,
		Payload: &database.KVStoreWrapper_DeleteKeyRequest{
			DeleteKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &database.DeleteKeyReply{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *client.Session
	if request.Transaction != nil {
		cs = s.tm.sessionCache[request.GetTransaction().GetClientId()]
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &database.DeleteKeyReply{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != client.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &database.DeleteKeyReply{}, errors.Wrap(err, "can't unmarshal response")
	}

	response := &database.DeleteKeyReply{}

	if request.Transaction == nil || cs.SeriesID != client.NoOPSeriesID {
		response.Transaction = &database.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	if resp.GetDeleteKeyReply() != nil {
		response.Ok = resp.GetDeleteKeyReply().GetOk()
	} else {
		response.Ok = false
	}

	return response, nil
}
