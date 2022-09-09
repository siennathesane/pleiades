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

func (rct *RaftControlTests) TestNewOrGetClusterManager() {
	nhc := buildTestNodeHostConfig(rct.T())
	node, err := NewRaftControlNode(nhc, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(node, "node must not be nil")

	cm, err := node.NewOrGetClusterManager()
	rct.Require().NoError(err, "there must not be an error when getting the cluster manager")
	rct.Require().NotNil(cm, "the cluster manager must not be nil")
}

func (rct *RaftControlTests) TestNewOrGetSessionManager() {
	nhc := buildTestNodeHostConfig(rct.T())
	node, err := NewRaftControlNode(nhc, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(node, "node must not be nil")

	sm, err := node.NewOrGetSessionManager()
	rct.Require().NoError(err, "there must not be an error when getting the session manager")
	rct.Require().NotNil(sm, "the session manager must not be nil")
}

func (rct *RaftControlTests) TestRequestLeaderTransfer() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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

	rs, err := firstNode.AddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
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
	rct.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	err = firstNode.LeaderTransfer(testClusterId, secondNodeClusterConfig.NodeID)
	time.Sleep(3000*time.Millisecond)

	leader, ok, err := secondNode.GetLeaderID(testClusterId)
	rct.Require().NoError(err, "there must not be an error when getting the leader id")
	rct.Require().True(ok, "it must be ok to fetc the leader information")
	rct.Require().Equal(secondNodeClusterConfig.NodeID, leader, "the second node must be the leader")
}

func (rct *RaftControlTests) TestRemoveData() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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

	rs, err := firstNode.AddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
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
	rct.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	rs, err = firstNode.DeleteNode(testClusterId, secondNodeClusterConfig.NodeID, 0, 3000*time.Millisecond)
	rct.Require().NoError(err, "there must not be an error when requesting to delete a node")
	rct.Require().NotNil(rs, "the request state must not be nil")

	select {
	case r := <-rs.ResultC():
		rct.Require().True(r.Completed(), "the delete node request must have completed successfully")
	}

	err = firstNode.RemoveData(testClusterId, secondNodeClusterConfig.NodeID)
	rct.Require().NoError(err, "there must not be an error when requesting to remove a dead node's data")
}

func (rct *RaftControlTests) TestRemoveNode() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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

	rs, err := firstNode.AddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
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
	rct.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	rs, err = firstNode.DeleteNode(testClusterId, secondNodeClusterConfig.NodeID, 0, 3000*time.Millisecond)
	rct.Require().NoError(err, "there must not be an error when requesting to delete a node")
	rct.Require().NotNil(rs, "the request state must not be nil")

	select {
	case r := <-rs.ResultC():
		rct.Require().True(r.Completed(), "the delete node request must have completed successfully")
	}

	membershipCtx, _ = context.WithTimeout(context.Background(), 3000*time.Millisecond)
	membership, err = firstNode.nh.SyncGetClusterMembership(membershipCtx,testClusterId)
	rct.Require().NoError(err, "there must not be an error when getting cluster membership")
	rct.Require().NotNil(membership, "the membership list must not be nil")
	rct.Require().Equal(1, len(membership.Nodes), "there must be only one node")
}

func (rct *RaftControlTests) TestStoreManager() {
	nhc := buildTestNodeHostConfig(rct.T())
	node, err := NewRaftControlNode(nhc, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the node")
	rct.Require().NotNil(node, "node must not be nil")

	cm, err := node.NewOrGetStoreManager()
	rct.Require().NoError(err, "there must not be an error when fetching the store manager")
	rct.Require().NotNil(cm, "the store manager must not be nil")
}

func (rct *RaftControlTests) TestClusterManager() {
	nhc := buildTestNodeHostConfig(rct.T())
	node, err := NewRaftControlNode(nhc, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the node")
	rct.Require().NotNil(node, "node must not be nil")

	cm, err := node.NewOrGetClusterManager()
	rct.Require().NoError(err, "there must not be an error when fetching the cluster manager")
	rct.Require().NotNil(cm, "the cluster manager must not be nil")
}

func (rct *RaftControlTests) TestSessionManager() {
	nhc := buildTestNodeHostConfig(rct.T())
	node, err := NewRaftControlNode(nhc, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the node")
	rct.Require().NotNil(node, "node must not be nil")

	sm, err := node.NewOrGetSessionManager()
	rct.Require().NoError(err, "there must not be an error when fetching the session manager")
	rct.Require().NotNil(sm, "the session manager must not be nil")
}

func (rct *RaftControlTests) TestAddWitness() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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
		IsWitness: true,
	}

	rs, err := firstNode.AddWitness(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
	rct.Require().NoError(err, "there must not be an error when requesting to add an observer")

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
	rct.Require().NotNil(1, len(membership.Witnesses), "there must be at least one witness")
}

func (rct *RaftControlTests) TestAddObserver() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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
		IsObserver: true,
	}

	rs, err := firstNode.AddObserver(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
	rct.Require().NoError(err, "there must not be an error when requesting to add an observer")

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
	rct.Require().NotNil(1, len(membership.Observers), "there must be at least one observer")
}

func (rct *RaftControlTests) TestAddNode() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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

	rs, err := firstNode.AddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
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
	rct.Require().Equal(2, len(membership.Nodes), "there must be at least one node")
}

func (rct *RaftControlTests) TestRaftAddress() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestShardConfig(rct.T())

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

	clusterConfig := buildTestShardConfig(rct.T())

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

	clusterConfig := buildTestShardConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	nodeUser, err := node.GetNodeUser(clusterConfig.ClusterID)
	rct.Require().NoError(err, "there must not be an error when fetching a node user")
	rct.Require().NotNil(nodeUser, "the node user must not be nil")

	rct.Require().Equal(clusterConfig.ClusterID, nodeUser.ClusterID(), "the nodeUser.ShardId must be equal to the configured cluster Id")
}

func (rct *RaftControlTests) TestGetLeaderId() {
	if testing.Short() {
		rct.T().Skipf("skipping")
	}

	firstNodeHostConfig := buildTestNodeHostConfig(rct.T())
	firstNode, err := NewRaftControlNode(firstNodeHostConfig, rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting the first node")
	rct.Require().NotNil(firstNode, "firstNode must not be nil")

	testClusterId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(rct.T())
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

	rs, err := firstNode.AddNode(testClusterId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0, 3000*time.Millisecond)
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
	rct.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	// we know it will be the first node who is the leader because the second node was added. the first node has all the
	// current indices, which means the second node can't lead. I think. idk, we'll see.
	leader, ok, err := firstNode.GetLeaderID(testClusterId)
	rct.Require().NoError(err, "there must not be an error when getting the leader info")
	rct.Require().True(ok, "it must be okay to fetch the leader")
	rct.Require().Equal(firstNodeClusterConfig.NodeID, leader, "the leader must be the first node")
}
