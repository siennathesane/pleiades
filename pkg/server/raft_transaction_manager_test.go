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

func TestSessionManager(t *testing.T) {
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
}

// we need to ensure that we use a single cluster the entire time to emulate multiple
// sessions in a single cluster. it's a bit... hand-wavey, but like, it works, so fuck it
func (smt *TransactionManagerTestSuite) SetupSuite() {

	smt.logger = utils.NewTestLogger(smt.T())

	smt.nh = buildTestNodeHost(smt.T())
	smt.Require().NotNil(smt.nh, "node must not be nil")

	shardConfig := buildTestShardConfig(smt.T())
	smt.shardId = shardConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[shardConfig.NodeID] = smt.nh.RaftAddress()

	err := smt.nh.StartCluster(nodeClusters, false, newTestStateMachine, shardConfig)
	smt.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(5000 * time.Millisecond)
}

func (smt *TransactionManagerTestSuite) TestGetNoOpSession() {
	sm := newSessionManager(smt.nh, smt.logger)

	transaction := sm.GetNoOpTransaction(smt.shardId)
	smt.Require().NotNil(transaction, "the client transaction must not be nil")

	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err := smt.nh.SyncPropose(proposeContext, sm.transactionToSession(transaction), []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().Panics(func() {
		cs := sm.transactionToSession(transaction)
		cs.ProposalCompleted()
	}, "finishing a noop proposal must panic")
}

// this test focuses primarily on ensuring that the transaction can be
// referenced, dereferenced, and reconstituted in order to ensure that
// transactions can be sent across the wire without losing their fidelity
// ultimately this is to ensure that we build out the transaction
// capabilities that we're not so functionally tied to dragonboat that
// we can't implement higher-order logic
func (smt *TransactionManagerTestSuite) TestGetTransaction() {
	sm := newSessionManager(smt.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	transaction, err := sm.GetTransaction(ctx, smt.shardId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(transaction, "the client session must not be nil")

	ogTransaction := *transaction
	cs := sm.transactionToSession(transaction)
	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")

	afterProposal := *sm.csToTransaction(cs)
	smt.Require().True(ogTransaction.TransactionId < afterProposal.TransactionId, "the post-proposal transaction id must have incremented")
	smt.logger.Printf("increase after a single transaction in a shard: %d", afterProposal.TransactionId - ogTransaction.TransactionId)
}

func (smt *TransactionManagerTestSuite) TestCloseTransaction() {
	sm := newSessionManager(smt.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	transaction, err := sm.GetTransaction(ctx, smt.shardId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(transaction, "the client session must not be nil")

	ogTransaction := *transaction
	cs := sm.transactionToSession(transaction)
	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = smt.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")

	afterProposal := *sm.csToTransaction(cs)
	smt.Require().True(ogTransaction.TransactionId < afterProposal.TransactionId, "the post-proposal transaction id must have incremented")
	smt.logger.Printf("increase after a single transaction in a shard: %d", afterProposal.TransactionId - ogTransaction.TransactionId)

	err = sm.CloseTransaction(ctx, &afterProposal)
	smt.Require().NoError(err, "there must not be an error when closing the transaction")
}
