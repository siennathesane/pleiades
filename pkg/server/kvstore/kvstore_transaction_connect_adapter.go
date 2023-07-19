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
	"github.com/cockroachdb/errors"
	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/kvpb/kvpbconnect"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/transactions"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	_ kvpbconnect.TransactionsServiceHandler = (*KvStoreTransactionConnectAdapter)(nil)
	_ runtime.ServiceHandler                      = (*KvStoreTransactionConnectAdapter)(nil)
)

type KvStoreTransactionConnectAdapterBuilderParams struct {
	fx.In
	TransactionManager runtime.ITransactionManager
	Logger             zerolog.Logger
}

type KvStoreTransactionConnectAdapterBuilderResults struct {
	fx.In

	ConnectAdapter *KvStoreTransactionConnectAdapter
}

type KvStoreTransactionConnectAdapter struct {
	http.Handler
	logger             zerolog.Logger
	transactionManager runtime.ITransactionManager
	path               string
}

func NewKvstoreTransactionConnectAdapter(transactionManager runtime.ITransactionManager, logger zerolog.Logger) *KvStoreTransactionConnectAdapter {
	adapter := &KvStoreTransactionConnectAdapter{logger: logger, transactionManager: transactionManager}
	adapter.path, adapter.Handler = kvpbconnect.NewTransactionsServiceHandler(adapter)
	return adapter
}

func (k *KvStoreTransactionConnectAdapter) Path() string {
	return k.path
}

func (k *KvStoreTransactionConnectAdapter) NewTransaction(ctx context.Context, c *connect.Request[kvpb.NewTransactionRequest]) (*connect.Response[kvpb.NewTransactionResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&kvpb.NewTransactionResponse{}), errors.New("shard id must not be 0")
	}

	transaction, err := k.transactionManager.GetTransaction(ctx, c.Msg.GetShardId())
	if err != nil {
		k.logger.Error().Err(err).Uint64("shard", c.Msg.GetShardId()).Msg("cannot create transaction")
	}
	return connect.NewResponse(&kvpb.NewTransactionResponse{Transaction: transaction}), nil
}

func (k *KvStoreTransactionConnectAdapter) CloseTransaction(ctx context.Context, c *connect.Request[kvpb.CloseTransactionRequest]) (*connect.Response[kvpb.CloseTransactionResponse], error) {
	transaction := c.Msg.GetTransaction()
	if err := k.checkTransaction(transaction); err != nil {
		return connect.NewResponse(&kvpb.CloseTransactionResponse{}), err
	}

	err := k.transactionManager.CloseTransaction(ctx, transaction)
	if err != nil {
		k.logger.Error().Err(err).Msg("can't close transaction")
	}
	return connect.NewResponse(&kvpb.CloseTransactionResponse{}), err
}

func (k *KvStoreTransactionConnectAdapter) Commit(ctx context.Context, c *connect.Request[kvpb.CommitRequest]) (*connect.Response[kvpb.CommitResponse], error) {
	transaction := c.Msg.GetTransaction()
	if err := k.checkTransaction(transaction); err != nil {
		return connect.NewResponse(&kvpb.CommitResponse{}), err
	}

	t := k.transactionManager.Commit(ctx, transaction)

	return connect.NewResponse(&kvpb.CommitResponse{Transaction: t}), nil
}

// todo (sienna): replace this with dclient.Session.ValidForProposal later.
func (k *KvStoreTransactionConnectAdapter) checkTransaction(t *kvpb.Transaction) error {
	// I don't think this can happen because it's a pointer, but better to be safe than sorry
	if t == nil {
		k.logger.Error().Err(transactions.ErrNilTransaction).Msg("attempted close of an empty transaction")
		return transactions.ErrNilTransaction
	}

	// check for noop or unset
	if t.GetTransactionId() == dclient.NoOPSeriesID {
		return transactions.ErrUnupportedTransaction
	}

	// check for unregister
	if t.GetTransactionId() == dclient.SeriesIDForUnregister {
		return transactions.ErrUnupportedTransaction
	}

	// check for pending registration
	if t.TransactionId == dclient.SeriesIDForRegister {
		return transactions.ErrUnupportedTransaction
	}

	return nil
}
