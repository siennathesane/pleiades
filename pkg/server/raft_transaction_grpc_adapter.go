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
	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var _ TransactionsServer = (*raftTransactionGrpcAdapter)(nil)

// todo (sienna): add caching for better session comparisons
type raftTransactionGrpcAdapter struct {
	logger zerolog.Logger
	transactionManager ITransactionManager
}

func (r *raftTransactionGrpcAdapter) NewTransaction(ctx context.Context, request *database.NewTransactionRequest) (*database.NewTransactionReply, error) {
	if request.GetShardId() <= systemShardStop {
		return &database.NewTransactionReply{}, ErrSystemShardRange
	}

	transaction, err := r.transactionManager.GetTransaction(ctx, request.GetShardId())
	if err != nil {
		r.logger.Error().Err(err).Uint64("shard", request.GetShardId()).Msg("cannot create transaction")
	}
	return &database.NewTransactionReply{Transaction: transaction}, nil
}

func (r *raftTransactionGrpcAdapter) CloseTransaction(ctx context.Context, request *database.CloseTransactionRequest) (*database.CloseTransactionReply, error) {
	transaction := request.GetTransaction()
	if err := r.checkTransaction(transaction); err != nil {
		return &database.CloseTransactionReply{}, err
	}

	err := r.transactionManager.CloseTransaction(ctx, transaction)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't close transaction")
	}
	return &database.CloseTransactionReply{}, err
}

func (r *raftTransactionGrpcAdapter) Commit(ctx context.Context, request *database.CommitRequest) (*database.CommitReply, error) {
	transaction := request.GetTransaction()
	if err := r.checkTransaction(transaction); err != nil {
		return &database.CommitReply{}, err
	}

	t := r.transactionManager.Commit(ctx, transaction)

	return &database.CommitReply{Transaction: t}, nil
}

func (r *raftTransactionGrpcAdapter) mustEmbedUnimplementedTransactionsServer() {
	//TODO implement me
	panic("implement me")
}

// todo (sienna): replace this with dclient.Session.ValidForProposal later.
func (r *raftTransactionGrpcAdapter) checkTransaction(t *database.Transaction) error {
	// I don't think this can happen because it's a pointer, but better to be safe than sorry
	if t == nil {
		r.logger.Error().Err(errNilTransaction).Msg("attempted close of an empty transaction")
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