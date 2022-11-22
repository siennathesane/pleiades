/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package transactions

import (
	"context"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/server/serverutils"
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
	logger         zerolog.Logger
	shardId        uint64
	nh             *dragonboat.NodeHost
	defaultTimeout time.Duration
}

// we need to ensure that we use a single cluster the entire time to emulate multiple
// sessions in a single cluster. it's a bit... hand-wavey, but like, it works, so fuck it
func (t *TransactionManagerTestSuite) SetupSuite() {

	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 500 * time.Millisecond

	t.nh = serverutils.BuildTestNodeHost(t.T())
	t.Require().NotNil(t.nh, "node must not be nil")

	shardConfig := serverutils.BuildTestShardConfig(t.T())
	t.shardId = shardConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, shardConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.defaultTimeout)
}

func (t *TransactionManagerTestSuite) TestGetNoOpSession() {
	params := &TransactionManagerBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	smRes := NewManager(params)
	sm := smRes.TransactionManager.(*TransactionManager)

	transaction := sm.GetNoOpTransaction(t.shardId)
	t.Require().NotNil(transaction, "the client transaction must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	t.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	_, err := t.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	t.Require().NoError(err, "there must not be an error when proposing a new message")

	t.Require().Panics(func() {
		cs.ProposalCompleted()
	}, "finishing a noop proposal must panic")
}

func (t *TransactionManagerTestSuite) TestGetTransaction() {
	params := &TransactionManagerBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	smRes := NewManager(params)
	sm := smRes.TransactionManager.(*TransactionManager)

	ctx, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, t.shardId)
	t.Require().NoError(err, "there must not be an error when getting the session")
	t.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	t.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	_, err = t.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	t.Require().NoError(err, "there must not be an error when proposing a new message")

	t.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")
}

func (t *TransactionManagerTestSuite) TestCloseTransaction() {
	params := &TransactionManagerBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	smRes := NewManager(params)
	sm := smRes.TransactionManager.(*TransactionManager)

	ctx, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, t.shardId)
	t.Require().NoError(err, "there must not be an error when getting the session")
	t.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	t.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	_, err = t.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	t.Require().NoError(err, "there must not be an error when proposing a new message")

	t.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")

	afterProposal := csToTransaction(*cs)
	err = sm.CloseTransaction(ctx, afterProposal)
	t.Require().NoError(err, "there must not be an error when closing the transaction")
}

func (t *TransactionManagerTestSuite) TestCommit() {
	params := &TransactionManagerBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	smRes := NewManager(params)
	sm := smRes.TransactionManager.(*TransactionManager)

	ctx, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	transaction, err := sm.GetTransaction(ctx, t.shardId)
	t.Require().NoError(err, "there must not be an error when getting the session")
	t.Require().NotNil(transaction, "the client session must not be nil")

	cs, ok := sm.sessionCache[transaction.GetClientId()]
	t.Require().True(ok, "the client session must exist in the cache")

	proposeContext, _ := context.WithTimeout(context.Background(), t.defaultTimeout)
	_, err = t.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	t.Require().NoError(err, "there must not be an error when proposing a new message")

	transactionResult := sm.Commit(proposeContext, transaction)

	t.Require().True(transaction.TransactionId < transactionResult.TransactionId, "the post-proposal transaction id must have incremented")
	t.logger.Printf("increase after a single transaction in a shard: %d", transactionResult.TransactionId-transaction.TransactionId)
}
