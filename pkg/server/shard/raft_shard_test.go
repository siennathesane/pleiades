/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package shard

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestShardManager(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping shard manager tests")
	}
	suite.Run(t, new(shardManagerTestSuite))
}

type shardManagerTestSuite struct {
	suite.Suite
	logger                 zerolog.Logger
	defaultTimeout         time.Duration
	extendedDefaultTimeout time.Duration
	nats                   *messaging.EmbeddedMessaging
	client                 *messaging.EmbeddedMessagingStreamClient
	eventHandler           *messaging.RaftEventHandler
}

func (t *shardManagerTestSuite) SetupSuite() {
	configuration.Get().SetDefault("server.datastore.basePath", t.T().TempDir())
	t.logger = utils.NewTestLogger(t.T())

	m, err := messaging.NewEmbeddedMessagingWithDefaults(t.logger)
	t.Require().NoError(err, "there must not be an error when creating the embedded nats")
	t.nats = m
	t.nats.Start()

	t.defaultTimeout = 200 * time.Millisecond
	t.extendedDefaultTimeout = 1000 * time.Millisecond

	pubSubClient, err := t.nats.GetPubSubClient()
	t.Require().NoError(err, "there must not be an error when creating the pubsub client")

	client, err := t.nats.GetStreamClient()
	t.Require().NoError(err, "there must not be an error when creating the queue client")
	t.client = client

	t.eventHandler = messaging.NewRaftEventHandler(pubSubClient, client, t.logger)
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestAddReplica() {
	firstTestHost := serverutils.BuildTestNodeHost(t.T())

	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.extendedDefaultTimeout)

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NoError(err, "there must not be an error when starting the second node")
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	// testShardId, secondNodeClusterConfig.NodeID, secondNode.RaftAddress(), utils.Timeout(t.extendedDefaultTimeout)
	err = shardManager.AddReplica(&raftv1.AddReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondNodeClusterConfig.NodeID,
		Hostname:  secondNode.RaftAddress(),
		Timeout:   int64(t.extendedDefaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a node")
	utils.Wait(t.extendedDefaultTimeout)

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")

	t.logger.Debug().Msg("waiting for leader to be elected in two node cluster")
	serverutils.WaitForReadyCluster(t.T(), testShardId, firstTestHost, t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be at two nodes")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestAddShardObserver() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NoError(err, "there must not be an error when starting the second node")
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		IsObserver:   true,
	}

	err = shardManager.AddReplicaObserver(testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when requesting to add an observer")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().NotNil(1, len(membership.Observers), "there must be at least one shard observer")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestAddShardWitness() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NoError(err, "there must not be an error when starting the second node")
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
		IsWitness:    true,
	}

	err = shardManager.AddReplicaWitness(testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when requesting to add an observer")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().NotNil(1, len(membership.Witnesses), "there must be at least one witness")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestDeleteReplica() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(&raftv1.AddReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondNodeClusterConfig.NodeID,
		Hostname:  secondNode.RaftAddress(),
		Timeout:   int64(t.extendedDefaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be two replicas")

	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when deleting a replica")

	membershipCtx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(1, len(membership.Nodes), "there must be only one replica")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestGetLeaderId() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(&raftv1.AddReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondNodeClusterConfig.NodeID,
		Hostname:  secondNode.RaftAddress(),
		Timeout:   int64(t.extendedDefaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	leader, ok, err := shardManager.GetLeaderId(testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the leader id")
	t.Require().True(ok, "the leader information must be available")
	t.Require().Equal(firstNodeClusterConfig.NodeID, leader, "the first node must be the leader")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestGetShardMembers() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting shard membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(1, len(membership.Nodes), "there must be at two replicas")

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first replica")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:       uint64(rand.Intn(10_000)),
		ClusterID:    testShardId,
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}

	err = shardManager.AddReplica(&raftv1.AddReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondNodeClusterConfig.NodeID,
		Hostname:  secondNode.RaftAddress(),
		Timeout:   int64(t.extendedDefaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a replica")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting shard membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be at two replicas")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestNewShard() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId := firstNodeClusterConfig.ClusterID

	// testShardId, firstNodeClusterConfig.NodeID, testStateMachineType, utils.Timeout(t.defaultTimeout)
	err := shardManager.NewShard(&raftv1.NewShardRequest{
		ShardId:   testShardId,
		ReplicaId: firstNodeClusterConfig.NodeID,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Timeout:   utils.Timeout(t.defaultTimeout).Milliseconds(),
	})
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting shard membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(1, len(membership.Nodes), "there must be one replica")

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first replica")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestRemoveData() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NoError(err, "there must not be an error when starting the second node")
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:              uint64(rand.Intn(10_000)),
		ClusterID:           testShardId,
		HeartbeatRTT:        10,
		ElectionRTT:         100,
		OrderedConfigChange: false,
	}

	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))

	err = shardManager.nh.SyncRequestAddNode(ctx, testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0)
	t.Require().NoError(err, "there must not be an error when requesting to add a node")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when requesting to delete a node")

	// the actually tested API
	err = shardManager.RemoveData(testShardId, secondNodeClusterConfig.NodeID)
	t.Require().NoError(err, "there must not be an error when requesting to remove a dead node's data")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestRemoveReplica() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := uint64(0)
	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId = firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	secondNode := serverutils.BuildTestNodeHost(t.T())
	t.Require().NoError(err, "there must not be an error when starting the second node")
	t.Require().NotNil(secondNode, "secondNode must not be nil")

	secondNodeClusterConfig := dconfig.Config{
		NodeID:              uint64(rand.Intn(10_000)),
		ClusterID:           testShardId,
		HeartbeatRTT:        10,
		ElectionRTT:         100,
		OrderedConfigChange: false,
	}

	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))

	err = shardManager.nh.SyncRequestAddNode(ctx, testShardId, secondNodeClusterConfig.NodeID, dragonboat.Target(secondNode.RaftAddress()), 0)
	t.Require().NoError(err, "there must not be an error when requesting to add a node")

	err = secondNode.StartCluster(nil, true, serverutils.NewTestStateMachine, secondNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(t.extendedDefaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be at least one node")

	// the actually tested API
	err = shardManager.RemoveReplica(testShardId, secondNodeClusterConfig.NodeID, utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when requesting to delete a replica")

	membershipCtx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err = shardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(1, len(membership.Nodes), "there must be only one node")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestStartReplica() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	firstParams := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(firstParams)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	firstShardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := rand.Uint64()
	firstTestReplicaId := rand.Uint64()

	// testShardId, firstTestReplicaId, testStateMachineType, utils.Timeout(t.defaultTimeout)
	err := firstShardManager.NewShard(&raftv1.NewShardRequest{
		ShardId:   testShardId,
		ReplicaId: firstTestReplicaId,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Timeout:   utils.Timeout(t.defaultTimeout).Milliseconds(),
	})
	t.Require().NoError(err, "there must not be an error when creating a new shard")
	utils.Wait(t.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := firstShardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	for i := 0; i < 5; i++ {
		_, err := firstShardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	secondTestHost := serverutils.BuildTestNodeHost(t.T())
	secondParams := ShardManagerBuilderParams{
		NodeHost: secondTestHost,
		Logger:   t.logger,
	}

	secondShardManagerResults := NewManager(secondParams)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	secondShardManager := secondShardManagerResults.RaftShardManager.(*RaftShardManager)

	secondTestReplicaId := rand.Uint64()

	// testShardId, secondTestReplicaId, secondShardManager.nh.RaftAddress(), utils.Timeout(t.defaultTimeout)
	err = firstShardManager.AddReplica(&raftv1.AddReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondTestReplicaId,
		Hostname:  secondShardManager.nh.RaftAddress(),
		Timeout:   int64(t.extendedDefaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a node")
	utils.Wait(t.defaultTimeout)

	// testShardId, secondTestReplicaId, testStateMachineType
	err = secondShardManager.StartReplica(&raftv1.StartReplicaRequest{
		ShardId:   testShardId,
		ReplicaId: secondTestReplicaId,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Restart:   false,
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a node")
	utils.Wait(t.defaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := firstShardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(2, len(membership.Nodes), "there must be at two nodes")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestStartObserverReplica() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	firstShardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	testShardId := rand.Uint64()
	firstTestReplicaId := rand.Uint64()

	err := firstShardManager.NewShard(&raftv1.NewShardRequest{
		ShardId:   testShardId,
		ReplicaId: firstTestReplicaId,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Timeout:   utils.Timeout(t.defaultTimeout).Milliseconds(),
	})
	t.Require().NoError(err, "there must not be an error when creating a new shard")
	utils.Wait(t.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := firstShardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	for i := 0; i < 5; i++ {
		_, err := firstShardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	secondTestHost := serverutils.BuildTestNodeHost(t.T())
	params = ShardManagerBuilderParams{
		NodeHost: secondTestHost,
		Logger:   t.logger,
	}

	secondShardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	secondShardManager := secondShardManagerResults.RaftShardManager.(*RaftShardManager)

	secondTestReplicaId := rand.Uint64()

	err = firstShardManager.AddReplicaObserver(testShardId, secondTestReplicaId, secondShardManager.nh.RaftAddress(), utils.Timeout(t.extendedDefaultTimeout))
	t.Require().NoError(err, "there must not be an error when requesting to add a node")
	utils.Wait(t.defaultTimeout)

	err = secondShardManager.StartReplicaObserver(&raftv1.StartReplicaObserverRequest{
		ShardId:   testShardId,
		ReplicaId: secondTestReplicaId,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Restart:   false,
	})
	t.Require().NoError(err, "there must not be an error when requesting to add a node")
	utils.Wait(t.defaultTimeout)

	membershipCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	membership, err := firstShardManager.nh.SyncGetClusterMembership(membershipCtx, testShardId)
	t.Require().NoError(err, "there must not be an error when getting cluster membership")
	t.Require().NotNil(membership, "the membership list must not be nil")
	t.Require().Equal(1, len(membership.Observers), "there must be at two nodes")
}

//goland:noinspection GoVetLostCancel
func (t *shardManagerTestSuite) TestStopReplica() {

	firstTestHost := serverutils.BuildTestNodeHost(t.T())
	params := ShardManagerBuilderParams{
		NodeHost: firstTestHost,
		Logger:   t.logger,
	}

	shardManagerResults := NewManager(params)
	t.Require().NotNil(shardManagerResults.RaftShardManager, "RaftShardManager must not be nil")
	shardManager := shardManagerResults.RaftShardManager.(*RaftShardManager)

	firstNodeClusterConfig := serverutils.BuildTestShardConfig(t.T())
	testShardId := firstNodeClusterConfig.ClusterID
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = shardManager.nh.RaftAddress()

	err := shardManager.nh.StartCluster(nodeClusters, false, serverutils.NewTestStateMachine, firstNodeClusterConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.extendedDefaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := shardManager.nh.SyncGetSession(ctx, testShardId)
	t.Require().NoError(err, "there must not be an error when fetching the client session from the first node")
	t.Require().NotNil(cs, "the first node's client session must not be nil")
	cancel()

	for i := 0; i < 5; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := shardManager.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")

		cs.ProposalCompleted()
	}

	_, err = shardManager.StopReplica(testShardId, firstNodeClusterConfig.NodeID)
	t.Require().NoError(err, "there must not be an error when stopping the replica")
}
