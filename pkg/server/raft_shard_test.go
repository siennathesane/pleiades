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

func TestShardManager(t *testing.T) {
	suite.Run(t, new(shardManagerTestSuite))
}

type shardManagerTestSuite struct {
	suite.Suite
	logger                 zerolog.Logger
	node                   *Node
	defaultTimeout         time.Duration
	extendedDefaultTimeout time.Duration
}

func (smts *shardManagerTestSuite) SetupSuite() {
	if testing.Short() {
		smts.T().Skipf("skipping")
	}

	smts.logger = utils.NewTestLogger(smts.T())

	smts.defaultTimeout = 3000 * time.Millisecond
	smts.extendedDefaultTimeout = 5000 * time.Millisecond
}

func (smts *shardManagerTestSuite) TestAddReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestRemoveData")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NoError(err, "there must not be an error when starting the second node")
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(testShardId, secondNodeClusterConfig.NodeID, secondNode.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be at two nodes")
}

func (smts *shardManagerTestSuite) TestAddShardObserver() {
	if testing.Short() {
		smts.T().Skipf("skipping TestRemoveData")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NoError(err, "there must not be an error when starting the second node")
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		IsObserver:    true,
	}

	err = shardManager.AddReplicaObserver(testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add an observer")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().NotNil(1, len(membership.Observers), "there must be at least one shard observer")
}

func (smts *shardManagerTestSuite) TestAddShardWitness() {
	if testing.Short() {
		smts.T().Skipf("skipping TestRemoveData")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NoError(err, "there must not be an error when starting the second node")
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		IsWitness:    true,
	}

	err = shardManager.AddReplicaWitness(testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add an observer")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().NotNil(1, len(membership.Witnesses), "there must be at least one witness")
}

func (smts *shardManagerTestSuite) TestDeleteReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestDeleteReplica")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(testShardId, secondNodeClusterConfig.NodeID, secondNode.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be two replicas")

	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when deleting a replica")

	membershipCtx, _ = context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Nodes), "there must be only one replica")
}

func (smts *shardManagerTestSuite) TestGetLeaderId() {
	if testing.Short() {
		smts.T().Skipf("skipping TestGetLeaderId")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(testShardId, secondNodeClusterConfig.NodeID, secondNode.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)
	
	leader, ok, err := shardManager.GetLeaderId(testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the leader id")
	smts.Require().True(ok, "the leader information must be available")
	smts.Require().Equal(firstNodeClusterConfig.NodeID, leader, "the first node must be the leader")
}

func (smts *shardManagerTestSuite) TestGetShardMembers() {
	if testing.Short() {
		smts.T().Skipf("skipping TestGetShardMembers")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting shard membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Nodes), "there must be at two replicas")

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first replica")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(testShardId, secondNodeClusterConfig.NodeID, secondNode.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ = context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting shard membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be at two replicas")
}

func (smts *shardManagerTestSuite) TestNewShard() {
	if testing.Short() {
		smts.T().Skipf("skipping TestGetShardMembers")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId := firstNodeClusterConfig.ClusterID

	err := shardManager.NewShard(testShardId, firstNodeClusterConfig.NodeID, testStateMachineType, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting shard membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Nodes), "there must be one replica")

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first replica")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}
}

func (smts *shardManagerTestSuite) TestRemoveData() {
	if testing.Short() {
		smts.T().Skipf("skipping TestRemoveData")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NoError(err, "there must not be an error when starting the second node")
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		OrderedConfigChange: false,
	}

	ctx, _ = context.WithTimeout(context.Background(), smts.defaultTimeout)

	err = shardManager.nh.SyncRequestAddNode(ctx, testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()),0)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to delete a node")

	// the actually tested API
	err = shardManager.RemoveData(testShardId, secondNodeClusterConfig.NodeID)
	smts.Require().NoError(err, "there must not be an error when requesting to remove a dead node's data")
}

func (smts *shardManagerTestSuite) TestRemoveReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestRemoveData")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := buildTestNodeHost(smts.T())
	smts.Require().NoError(err, "there must not be an error when starting the second node")
	smts.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		OrderedConfigChange: false,
	}

	ctx, _ = context.WithTimeout(context.Background(), smts.defaultTimeout)

	err = shardManager.nh.SyncRequestAddNode(ctx, testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()),0)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")

	err = secondNode.StartCluster(nil, true, newTestStateMachine, secondNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	// the actually tested API
	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to delete a replica")

	membershipCtx, _ = context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Nodes), "there must be only one node")
}

func (smts *shardManagerTestSuite) TestStartReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestStartReplica")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	firstShardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(firstShardManager, "firstShardManager must not be nil")

	testShardId := rand.Uint64()
	firstTestReplicaId := rand.Uint64()
	err := firstShardManager.NewShard(testShardId, firstTestReplicaId, testStateMachineType, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when creating a new shard")
	time.Sleep(smts.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := firstShardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	for i := 0; i < 5; i++ {
		_, err := firstShardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	secondTestHost := buildTestNodeHost(smts.T())
	secondShardManager := newShardManager(secondTestHost, smts.logger)
	smts.Require().NotNil(secondShardManager, "firstShardManager must not be nil")

	secondTestReplicaId := rand.Uint64()

	err = firstShardManager.AddReplica(testShardId, secondTestReplicaId, secondShardManager.nh.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	err = secondShardManager.StartReplica(testShardId, secondTestReplicaId, testStateMachineType)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := firstShardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(2, len(membership.Nodes), "there must be at two nodes")
}

func (smts *shardManagerTestSuite) TestStartObserverReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestStartObserverReplica")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	firstShardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(firstShardManager, "firstShardManager must not be nil")

	testShardId := rand.Uint64()
	firstTestReplicaId := rand.Uint64()
	err := firstShardManager.NewShard(testShardId, firstTestReplicaId, testStateMachineType, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when creating a new shard")
	time.Sleep(smts.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := firstShardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	for i := 0; i < 5; i++ {
		_, err := firstShardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	secondTestHost := buildTestNodeHost(smts.T())
	secondShardManager := newShardManager(secondTestHost, smts.logger)
	smts.Require().NotNil(secondShardManager, "firstShardManager must not be nil")

	secondTestReplicaId := rand.Uint64()

	err = firstShardManager.AddReplicaObserver(testShardId, secondTestReplicaId, secondShardManager.nh.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	err = secondShardManager.StartReplicaObserver(testShardId, secondTestReplicaId, testStateMachineType)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := firstShardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Observers), "there must be at two nodes")
}

func (smts *shardManagerTestSuite) TestStartWitnessReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestStartObserverReplica")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	firstShardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(firstShardManager, "firstShardManager must not be nil")

	testShardId := rand.Uint64()
	firstTestReplicaId := rand.Uint64()
	err := firstShardManager.NewShard(testShardId, firstTestReplicaId, testStateMachineType, smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when creating a new shard")
	time.Sleep(smts.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := firstShardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	for i := 0; i < 5; i++ {
		_, err := firstShardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	secondTestHost := buildTestNodeHost(smts.T())
	secondShardManager := newShardManager(secondTestHost, smts.logger)
	smts.Require().NotNil(secondShardManager, "firstShardManager must not be nil")

	secondTestReplicaId := rand.Uint64()

	err = firstShardManager.AddReplicaWitness(testShardId, secondTestReplicaId, secondShardManager.nh.RaftAddress(), smts.defaultTimeout)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	err = secondShardManager.StartReplicaWitness(testShardId, secondTestReplicaId, testStateMachineType)
	smts.Require().NoError(err, "there must not be an error when requesting to add a node")
	time.Sleep(smts.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
	membership, err := firstShardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	smts.Require().NoError(err, "there must not be an error when getting cluster membership")
	smts.Require().NotNil(membership, "the membership list must not be nil")
	smts.Require().Equal(1, len(membership.Witnesses), "there must be at two nodes")
}


func (smts *shardManagerTestSuite) TestStopReplica() {
	if testing.Short() {
		smts.T().Skipf("skipping TestGet")
	}

	firstTestHost := buildTestNodeHost(smts.T())
	shardManager := newShardManager(firstTestHost, smts.logger)
	smts.Require().NotNil(shardManager, "shardManager must not be nil")

	testShardId := uint64(0)
	firstNodeClusterConfig := buildTestShardConfig(smts.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	smts.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(smts.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), smts.defaultTimeout)
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	smts.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	smts.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), smts.defaultTimeout)
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		smts.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	_, err = shardManager.StopReplica(testShardId)
	smts.Require().NoError(err, "there must not be an error when stopping the replia")
}
