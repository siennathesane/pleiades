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
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dclient "github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var (
	_ ITransactionManager = (*raftTransactionManager)(nil)
)

func newSessionManager(nh *dragonboat.NodeHost, logger zerolog.Logger) *raftTransactionManager {
	l := logger.With().Str("component", "session-manager").Logger()
	return &raftTransactionManager{l, nh, make(map[uint64]*dclient.Session)}
}

type raftTransactionManager struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost

	// todo (sienna): there has to be a better/faster version of this
	sessionCache map[uint64]*dclient.Session
}

func (t *raftTransactionManager) CloseTransaction(ctx context.Context, transaction *database.Transaction) error {
	t.logger.Debug().Uint64("shard", transaction.ShardId).Msg("closing transaction")

	cs, ok := t.sessionCache[transaction.GetClientId()]
	if !ok {
		return errors.New("transaction not found")
	}

	err := t.nh.SyncCloseSession(ctx, cs)
	if err != nil {
		t.logger.Error().Err(err).Msg("can't close transaction")
	}
	delete(t.sessionCache, cs.ClientID)

	return err
}

func (t *raftTransactionManager) Commit(ctx context.Context, transaction *database.Transaction) *database.Transaction {
	// nb (sienna): I know, I know. stop judging me.
	// is this hacky? yes.
	// does it work? yes.
	// is it the right thing to do now? no.
	// will it help later? yes.

	t.logger.Debug().Uint64("shard", transaction.ShardId).Msg("closing transaction")

	cs, ok := t.sessionCache[transaction.GetClientId()]
	if !ok {
		return &database.Transaction{}
	}

	cs.ProposalCompleted()

	ta := csToTransaction(*cs)
	return ta
}

func (t *raftTransactionManager) GetNoOpTransaction(shardId uint64) *database.Transaction {
	t.logger.Debug().Uint64("shard", shardId).Msg("getting noop transaction")
	cs := t.nh.GetNoOPSession(shardId)
	t.sessionCache[cs.ClientID] = cs
	return csToTransaction(*cs)
}

func (t *raftTransactionManager) GetTransaction(ctx context.Context, shardId uint64) (*database.Transaction, error) {
	t.logger.Debug().Uint64("shard", shardId).Msg("getting transaction")
	cs, err := t.nh.SyncGetSession(ctx, shardId)
	if err != nil {
		t.logger.Error().Err(err).Uint64("shard", shardId).Msg("can't get transaction")
		return nil, err
	}

	t.sessionCache[cs.ClientID] = cs

	return csToTransaction(*cs), nil
}

func csToTransaction(cs dclient.Session) *database.Transaction {
	return &database.Transaction{
		ShardId:       cs.ClusterID,
		ClientId:      cs.ClientID,
		TransactionId: cs.SeriesID,
		RespondedTo:   cs.RespondedTo,
	}
}
