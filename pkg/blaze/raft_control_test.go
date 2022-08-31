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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftControl(t *testing.T) {
	suite.Run(t, new(RaftControlTests))
}

type RaftControlTests struct {
	suite.Suite
	logger zerolog.Logger
}

func (rct *RaftControlTests) SetupSuite() {
	rct.logger = utils.NewTestLogger(rct.T())
}

func (rct *RaftControlTests) TestAddNode() {
	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestClusterConfig(rct.T())
	testClusterId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = firstNode.RaftAddress()

	err = firstNode.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(5000 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	cs, err := firstNode.nh.SyncGetSession(ctx, testClusterId)
	rct.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	rct.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		_, err := firstNode.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		rct.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNodeHostConfig := buildTestNodeHostConfig(rct.T())
	secondNode, err := NewRaftControlNode(secondNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the second node")
	rct.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testClusterId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	rs, err := firstNode.RequestAddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
	rct.Require().NoError(err, "there must not be an error when requesting to add a node")

	select {
	case r := <-rs.ResultC():
		rct.Require().True(r.Completed(), "the request must have completed successfully")
	}

	err = secondNode.nh.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(5000 * time.Millisecond)

	membershipCtx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	membership, err := firstNode.nh.SyncGetClusterMembership(membershipCtx,testClusterId)
	rct.Require().NoError(err, "there must not be an error when getting cluster membership")
	rct.Require().NotNil(membership, "the membership list must not be nil")
}

func (rct *RaftControlTests) TestRaftAddress() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	addr := node.RaftAddress()
	rct.Require().NotNil(addr, "the raft address must not be nil")
	rct.T().Logf("raft addr: %s", addr)
}

func (rct *RaftControlTests) TestId() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	id := node.ID()
	rct.Require().NotEmpty(id, "the node id must not be empty")
	rct.T().Logf("node id: %s", id)
}

func (rct *RaftControlTests) TestGetNodeUser() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	nodeUser, err := node.GetNodeUser(clusterConfig.ClusterID)
	rct.Require().NoError(err, "there must not be an error when fetching a node user")
	rct.Require().NotNil(nodeUser, "the node user must not be nil")

	rct.Require().Equal(clusterConfig.ClusterID, nodeUser.ClusterID(), "the nodeUser.ClusterID must be equal to the configured cluster ID")
}

func (rct *RaftControlTests) TestGetId() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	leader, ok, err := node.GetLeaderID(clusterConfig.ClusterID)
	rct.Require().NoError(err, "there must not be an error when fetching the leader")

	// until there's a cluster in place, it's okay to skip the rest of it since we're really just validating it works
	// todo (sienna): implement clustering logic for this test
	if !ok && leader == 0 {
		rct.T().Skipf("skipping due to unimplemented cluster functionality for test")
	}

	rct.Assert().True(ok, "the leader information should be available")
	rct.Assert().NotEmpty(leader, "the leader is should not be 0")
}
