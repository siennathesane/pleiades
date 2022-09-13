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
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestTransactionManager(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping session manager tests")
	}
	suite.Run(t, new(TransactionManagerTestSuite))
}

type TransactionManagerTestSuite struct {
	suite.Suite
	logger  zerolog.Logger
	shardId uint64
	nh      *dragonboat.NodeHost
	defaultTimeout time.Duration
}

// we need to ensure that we use a single cluster the entire time to emulate multiple
// sessions in a single cluster. it's a bit... hand-wavey, but like, it works, so fuck it
func (smt *TransactionManagerTestSuite) SetupSuite() {

	smt.logger = utils.NewTestLogger(smt.T())
	smt.defaultTimeout = 500 * time.Millisecond

	smt.nh = buildTestNodeHost(smt.T())
	smt.Require().NotNil(smt.nh, "node must not be nil")

	shardConfig := buildTestShardConfig(smt.T())
	smt.shardId = shardConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[shardConfig.NodeID] = smt.nh.RaftAddress()

	err := smt.nh.StartCluster(nodeClusters, false, newTestStateMachine, shardConfig)
	smt.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smt.defaultTimeout)
}

func (smt *TransactionManagerTestSuite) TestGetNoOpSession() {
	sm := newTransactionManager(smt.nh, smt.logger)

	transaction := sm.GetNoOpTransaction(smt.shardId)
	smt.Require().NotNil(transaction, "the client transaction must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	smt.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	_, err := smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().Panics(func() {
		cs.ProposalCompleted()
	}, "finishing a noop proposal must panic")
}

func (smt *TransactionManagerTestSuite) TestGetTransaction() {
	sm := newTransactionManager(smt.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, smt.shardId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	smt.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	_, err = smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")
}

func (smt *TransactionManagerTestSuite) TestCloseTransaction() {
	sm := newTransactionManager(smt.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, smt.shardId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	smt.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	_, err = smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")

	afterProposal := csToTransaction(*cs)
	err = sm.CloseTransaction(ctx, afterProposal)
	smt.Require().NoError(err, "there must not be an error when closing the transaction")
}

func (smt *TransactionManagerTestSuite) TestCommit() {
	sm := newTransactionManager(smt.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, smt.shardId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	smt.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), smt.defaultTimeout)
	_, err = smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	transactionResult := sm.Commit(proposeContext, transaction)

	smt.Require().True(transaction.TransactionId < transactionResult.TransactionId, "the post-proposal transaction id must have incremented")
	smt.logger.Printf("increase after a single transaction in a shard: %d", transactionResult.TransactionId - transaction.TransactionId)
}
