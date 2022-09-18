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
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var _ kvstorev1connect.TransactionsServiceHandler = (*kvstoreTransactionConnectAdapter)(nil)

type kvstoreTransactionConnectAdapter struct {
	logger             zerolog.Logger
	transactionManager ITransactionManager
}

func (k *kvstoreTransactionConnectAdapter) NewTransaction(ctx context.Context, c *connect.Request[kvstorev1.NewTransactionRequest]) (*connect.Response[kvstorev1.NewTransactionResponse], error) {
	if c.Msg.GetShardId() == 0 {
		return connect.NewResponse(&kvstorev1.NewTransactionResponse{}), ErrSystemShardRange
	}

	transaction, err := k.transactionManager.GetTransaction(ctx, c.Msg.GetShardId())
	if err != nil {
		k.logger.Error().Err(err).Uint64("shard", c.Msg.GetShardId()).Msg("cannot create transaction")
	}
	return connect.NewResponse(&kvstorev1.NewTransactionResponse{Transaction: transaction}), nil
}

func (k *kvstoreTransactionConnectAdapter) CloseTransaction(ctx context.Context, c *connect.Request[kvstorev1.CloseTransactionRequest]) (*connect.Response[kvstorev1.CloseTransactionResponse], error) {
	transaction := c.Msg.GetTransaction()
	if err := k.checkTransaction(transaction); err != nil {
		return connect.NewResponse(&kvstorev1.CloseTransactionResponse{}), err
	}

	err := k.transactionManager.CloseTransaction(ctx, transaction)
	if err != nil {
		k.logger.Error().Err(err).Msg("can't close transaction")
	}
	return connect.NewResponse(&kvstorev1.CloseTransactionResponse{}), err
}

func (k *kvstoreTransactionConnectAdapter) Commit(ctx context.Context, c *connect.Request[kvstorev1.CommitRequest]) (*connect.Response[kvstorev1.CommitResponse], error) {
	transaction := c.Msg.GetTransaction()
	if err := k.checkTransaction(transaction); err != nil {
		return connect.NewResponse(&kvstorev1.CommitResponse{}), err
	}

	t := k.transactionManager.Commit(ctx, transaction)

	return connect.NewResponse(&kvstorev1.CommitResponse{Transaction: t}), nil
}

// todo (sienna): replace this with dclient.Session.ValidForProposal later.
func (k *kvstoreTransactionConnectAdapter) checkTransaction(t *kvstorev1.Transaction) error {
	// I don't think this can happen because it's a pointer, but better to be safe than sorry
	if t == nil {
		k.logger.Error().Err(errNilTransaction).Msg("attempted close of an empty transaction")
		return errNilTransaction
	}

	// check for noop or unset
	if t.GetTransactionId() == dclient.NoOPSeriesID {
		return errUnupportedTransaction
	}

	// check for unregister
	if t.GetTransactionId() == dclient.SeriesIDForUnregister {
		return errUnupportedTransaction
	}

	// check for pending registration
	if t.TransactionId == dclient.SeriesIDForRegister {
		return errUnupportedTransaction
	}

	return nil
}
