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
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var (
	_ ITransactionManager = (*transactionManager)(nil)
)

func newSessionManager(nh *dragonboat.NodeHost, logger zerolog.Logger) *transactionManager {
	l := logger.With().Str("component", "session-manager").Logger()
	return &transactionManager{l, nh}
}

type transactionManager struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (t *transactionManager) CloseTransaction(ctx context.Context, transaction *database.Transaction) error {
	t.logger.Debug().Uint64("shard", transaction.ShardId).Msg("closing transaction")
	return t.nh.SyncCloseSession(ctx, &client.Session{
		ClusterID:   transaction.ShardId,
		ClientID:    transaction.ClientId,
		SeriesID:    transaction.TransactionId,
		RespondedTo: transaction.RespondedTo,
	})
}

func (t *transactionManager) Complete(ctx context.Context, transaction *database.Transaction) *database.Transaction {
	// nb (sienna): I know, I know. stop judging me.
	// is this hacky? yes.
	// does it work? yes.
	// is it the right thing to do now? no.
	// will it help later? yes.

	t.logger.Debug().Uint64("shard", transaction.ShardId).Msg("closing transaction")
	cs := t.transactionToSession(transaction)
	cs.ProposalCompleted()
	return t.csToTransaction(cs)
}

func (t *transactionManager) GetNoOpTransaction(shardId uint64) *database.Transaction {
	t.logger.Debug().Uint64("shard", shardId).Msg("getting noop transaction")
	cs := t.nh.GetNoOPSession(shardId)
	return t.csToTransaction(cs)
}

func (t *transactionManager) GetTransaction(ctx context.Context, shardId uint64) (*database.Transaction, error) {
	t.logger.Debug().Uint64("shard", shardId).Msg("getting transaction")
	cs, err := t.nh.SyncGetSession(ctx, shardId)
	if err != nil {
		t.logger.Error().Err(err).Uint64("shard", shardId).Msg("can't get transaction")
		return nil, err
	}

	return t.csToTransaction(cs), nil
}

func (t *transactionManager) csToTransaction(cs *client.Session) *database.Transaction {
	return &database.Transaction{
		ShardId:       cs.ClusterID,
		ClientId:      cs.ClientID,
		TransactionId: cs.SeriesID,
		RespondedTo:   cs.RespondedTo,
	}
}

func (t *transactionManager) transactionToSession(transaction *database.Transaction) *client.Session {
	return &client.Session{
		ClusterID:   transaction.ShardId,
		ClientID:    transaction.ClientId,
		SeriesID:    transaction.TransactionId,
		RespondedTo: transaction.RespondedTo,
	}
}