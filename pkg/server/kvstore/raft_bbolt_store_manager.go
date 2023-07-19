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
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dclient "github.com/lni/dragonboat/v3/client"
	aerrs "github.com/mxplusb/pleiades/pkg/errorspb"
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/routing"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	_ runtime.IKVStore = (*BboltStoreManager)(nil)
)

type BboltStoreManagerBuilderParams struct {
	fx.In

	TransactionManager runtime.ITransactionManager
	NodeHost           *dragonboat.NodeHost
	Logger             zerolog.Logger
}

type BboltStoreManagerBuilderResults struct {
	fx.Out

	KVStoreManager runtime.IKVStore
}

func NewBboltStoreManager(params BboltStoreManagerBuilderParams) BboltStoreManagerBuilderResults {
	l := params.Logger.With().Str("component", "store-manager").Logger()
	return BboltStoreManagerBuilderResults{
		KVStoreManager: &BboltStoreManager{
			l,
			params.TransactionManager,
			params.NodeHost,
			&routing.ShardRouter{},
			1000 * time.Millisecond,
		},
	}
}

type BboltStoreManager struct {
	logger         zerolog.Logger
	tm             runtime.ITransactionManager
	nh             *dragonboat.NodeHost
	shardRouter    *routing.ShardRouter
	defaultTimeout time.Duration
}

func (s *BboltStoreManager) CreateAccount(request *kvpb.CreateAccountRequest) (*kvpb.CreateAccountResponse, error) {

	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.CreateAccountResponse{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &kvpb.CreateAccountResponse{}, kv.ErrInvalidOwner
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_CREATE_ACCOUNT_REQUEST,
		Payload: &kvpb.KVStoreWrapper_CreateAccountRequest{
			CreateAccountRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.CreateAccountResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.CreateAccountResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.CreateAccountResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.CreateAccountResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.CreateAccountResponse{}
	if resp.GetCreateAccountReply() == nil {
		response.AccountDescriptor = &kvpb.AccountDescriptor{}
	} else {
		response.AccountDescriptor = resp.GetCreateAccountReply().GetAccountDescriptor()
	}

	if request.Transaction == nil {
		response.Transaction = &kvpb.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *BboltStoreManager) DeleteAccount(request *kvpb.DeleteAccountRequest) (*kvpb.DeleteAccountResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.DeleteAccountResponse{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &kvpb.DeleteAccountResponse{}, kv.ErrInvalidOwner
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_DELETE_ACCOUNT_REQUEST,
		Payload: &kvpb.KVStoreWrapper_DeleteAccountRequest{
			DeleteAccountRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.DeleteAccountResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.DeleteAccountResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.DeleteAccountResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.DeleteAccountResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.DeleteAccountResponse{
		Ok: resp.GetDeleteAccountReply().GetOk(),
	}

	if request.Transaction == nil || cs.SeriesID != dclient.NoOPSeriesID {
		response.Transaction = &kvpb.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *BboltStoreManager) CreateBucket(request *kvpb.CreateBucketRequest) (*kvpb.CreateBucketResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.CreateBucketResponse{}, kv.ErrInvalidAccount
	}

	owner := request.GetOwner()
	if owner == "" {
		s.logger.Trace().Msg("empty owner value")
		return &kvpb.CreateBucketResponse{}, kv.ErrInvalidOwner
	}

	bucketName := request.GetName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &kvpb.CreateBucketResponse{}, kv.ErrInvalidBucketName
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_CREATE_BUCKET_REQUEST,
		Payload: &kvpb.KVStoreWrapper_CreateBucketRequest{
			CreateBucketRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.CreateBucketResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.CreateBucketResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.CreateBucketResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.CreateBucketResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.CreateBucketResponse{}
	if resp.GetCreateBucketReply() == nil {
		response.BucketDescriptor = &kvpb.BucketDescriptor{}
	} else {
		response.BucketDescriptor = resp.GetCreateBucketReply().GetBucketDescriptor()
	}

	if request.Transaction == nil {
		response.Transaction = &kvpb.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *BboltStoreManager) DeleteBucket(request *kvpb.DeleteBucketRequest) (*kvpb.DeleteBucketResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.DeleteBucketResponse{}, kv.ErrInvalidAccount
	}

	name := request.GetName()
	if name == "" {
		s.logger.Trace().Msg("empty name value")
		return &kvpb.DeleteBucketResponse{}, kv.ErrInvalidOwner
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  name,
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_DELETE_BUCKET_REQUEST,
		Payload: &kvpb.KVStoreWrapper_DeleteBucketRequest{
			DeleteBucketRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.DeleteBucketResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.DeleteBucketResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.DeleteBucketResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.DeleteBucketResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.DeleteBucketResponse{
		Ok: resp.GetDeleteBucketReply().GetOk(),
	}

	if request.Transaction == nil || cs.SeriesID != dclient.NoOPSeriesID {
		response.Transaction = &kvpb.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *BboltStoreManager) GetKey(request *kvpb.GetKeyRequest) (*kvpb.GetKeyResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.GetKeyResponse{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &kvpb.GetKeyResponse{}, kv.ErrInvalidOwner
	}

	keyName := request.GetKey()
	if len(keyName) == 0 {
		s.logger.Trace().Msg("empty key name")
		return &kvpb.GetKeyResponse{}, errors.New("empty key name")
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_GET_KEY_REQUEST,
		Payload: &kvpb.KVStoreWrapper_GetKeyRequest{
			GetKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.GetKeyResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	shard := s.shardRouter.AccountToShard(account)

	result, err := s.nh.SyncRead(ctx, shard, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.GetKeyResponse{}, errors.Wrap(err, "can't apply message")
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.([]byte))
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.GetKeyResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.GetKeyResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.GetKeyResponse{}
	kvp := resp.GetGetKeyReply()
	if kvp == nil {
		response.KeyValuePair = &kvpb.KeyValue{}
	} else {
		response.KeyValuePair = kvp.GetKeyValuePair()
	}

	return response, nil
}

func (s *BboltStoreManager) PutKey(request *kvpb.PutKeyRequest) (*kvpb.PutKeyResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.PutKeyResponse{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucketName value")
		return &kvpb.PutKeyResponse{}, kv.ErrInvalidOwner
	}

	keyValuePair := request.GetKeyValuePair()
	if keyValuePair == nil {
		s.logger.Trace().Msg("empty key value pair")
		return &kvpb.PutKeyResponse{}, kv.ErrInvalidBucketName
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_PUT_KEY_REQUEST,
		Payload: &kvpb.KVStoreWrapper_PutKeyRequest{
			PutKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.PutKeyResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.PutKeyResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.PutKeyResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.PutKeyResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.PutKeyResponse{}

	if request.Transaction == nil || cs.SeriesID != dclient.NoOPSeriesID {
		response.Transaction = &kvpb.Transaction{}
	} else {
		response.Transaction = csToTransaction(*cs)
	}

	return response, nil
}

func (s *BboltStoreManager) DeleteKey(request *kvpb.DeleteKeyRequest) (*kvpb.DeleteKeyResponse, error) {
	account := request.GetAccountId()
	if account == 0 {
		s.logger.Trace().Msg("empty account value")
		return &kvpb.DeleteKeyResponse{}, kv.ErrInvalidAccount
	}

	bucketName := request.GetBucketName()
	if bucketName == "" {
		s.logger.Trace().Msg("empty bucket name")
		return &kvpb.DeleteKeyResponse{}, kv.ErrInvalidOwner
	}

	key := request.GetKey()
	if len(key) == 0 {
		s.logger.Trace().Msg("empty key value pair")
		return &kvpb.DeleteKeyResponse{}, kv.ErrInvalidBucketName
	}

	req := &kvpb.KVStoreWrapper{
		Account: request.GetAccountId(),
		Bucket:  bucketName,
		Typ:     kvpb.KVStoreWrapper_REQUEST_TYPE_DELETE_KEY_REQUEST,
		Payload: &kvpb.KVStoreWrapper_DeleteKeyRequest{
			DeleteKeyRequest: request,
		},
	}

	cmd, err := req.MarshalVT()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't marshal request")
		return &kvpb.DeleteKeyResponse{}, errors.Wrap(err, "can't marshal request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(s.defaultTimeout))
	defer cancel()

	var cs *dclient.Session
	if request.Transaction != nil {
		sess, ok := s.tm.SessionFromClientId(request.GetTransaction().GetClientId())
		if !ok {
			// todo (sienna): figure out what to do here
			s.logger.Error().Uint64("dclient-id", request.GetTransaction().GetClientId()).Msg("session not found in cache")
		}
		cs = sess
	} else {
		shard := s.shardRouter.AccountToShard(account)
		cs = s.nh.GetNoOPSession(shard)
	}

	result, err := s.nh.SyncPropose(ctx, cs, cmd)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't apply message")
		return &kvpb.DeleteKeyResponse{}, errors.Wrap(err, "can't apply message")
	}

	if cs.SeriesID != dclient.NoOPSeriesID {
		cs.ProposalCompleted()
	}

	resp := &kvpb.KVStoreWrapper{}
	err = resp.UnmarshalVT(result.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal response")
		return &kvpb.DeleteKeyResponse{}, errors.Wrap(err, "can't unmarshal response")
	}

	if resp.Typ == kvpb.KVStoreWrapper_REQUEST_TYPE_RECOVERABLE_ERROR {
		errMsg := &aerrs.Error{}
		errMsg = resp.GetError()

		return &kvpb.DeleteKeyResponse{}, errors.New(errMsg.Message)
	}

	response := &kvpb.DeleteKeyResponse{}

	if request.Transaction == nil || cs.SeriesID != dclient.NoOPSeriesID {
		response.Transaction = &kvpb.Transaction{}
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

func csToTransaction(cs dclient.Session) *kvpb.Transaction {
	return &kvpb.Transaction{
		ShardId:       cs.ClusterID,
		ClientId:      cs.ClientID,
		TransactionId: cs.SeriesID,
		RespondedTo:   cs.RespondedTo,
	}
}
