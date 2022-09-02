/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"context"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestSessionManager(t *testing.T) {
	suite.Run(t, new(SessionManagerTests))
}

type SessionManagerTests struct {
	suite.Suite
	logger zerolog.Logger
	clusterId uint64
	node *Node
}

// we need to ensure that we use a single cluster the entire time to emulate multiple
// sessions in a single cluster. it's a bit... hand-wavey, but like, it works, so fuck it
func (smt *SessionManagerTests) SetupSuite() {
	if testing.Short() {
		smt.T().Skipf("skipping")
	}

	smt.logger = utils.NewTestLogger(smt.T())

	nodeHostConfig := buildTestNodeHostConfig(smt.T())
	node, err := NewRaftControlNode(nodeHostConfig, smt.logger)
	smt.Require().NoError(err, "there must not be an error when starting the node")
	smt.Require().NotNil(node, "node must not be nil")
	smt.node = node

	clusterConfig := buildTestClusterConfig(smt.T())
	smt.clusterId = clusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[clusterConfig.NodeID] = node.RaftAddress()

	err = node.nh.StartCluster(nodeClusters, false, newTestStateMachine, clusterConfig)
	smt.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(5000 * time.Millisecond)
}

func (smt *SessionManagerTests) TestGetNoOpSession() {
	sm := newSessionManager(smt.node.nh, smt.logger)

	cs := sm.GetNoOpSession(smt.clusterId)
	smt.Require().NotNil(cs, "the client session must not be nil")

	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err := smt.node.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().Panics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")
}

func (smt *SessionManagerTests) TestGetSession() {
	sm := newSessionManager(smt.node.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	cs, err := sm.GetSession(ctx, smt.clusterId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(cs, "the client session must not be nil")

	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = smt.node.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")
}

func (smt *SessionManagerTests) TestCloseSession() {
	sm := newSessionManager(smt.node.nh, smt.logger)

	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	cs, err := sm.GetSession(ctx, smt.clusterId)
	smt.Require().NoError(err, "there must not be an error when getting the session")
	smt.Require().NotNil(cs, "the client session must not be nil")

	proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	_, err = smt.node.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	smt.Require().NoError(err, "there must not be an error when proposing a new message")

	smt.Require().NotPanics(func() {
		cs.ProposalCompleted()
	}, "finishing a proposal must not panic")

	closeCtx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	err = smt.node.nh.SyncCloseSession(closeCtx, cs)
	smt.Require().NoError(err, "there must not be an error when closing the session")
	smt.Require().Equal(client.SeriesIDForUnregister, cs.SeriesID, "the series id must be set for unregister")
}
