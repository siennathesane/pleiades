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
	"net"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/api/v1/raft"
	"github.com/mxplusb/pleiades/pkg/utils"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestRaftShardGrpcAdapter(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping shard grpc tests")
	}
	suite.Run(t, new(RaftShardGrpcAdapterTestSuite))
}

type RaftShardGrpcAdapterTestSuite struct {
	suite.Suite
	logger                 zerolog.Logger
	adapter                *raftShardGrpcAdapter
	conn                   *grpc.ClientConn
	srv                    *grpc.Server
	testShardId            uint64
	testClusterConfig      dconfig.Config
	testShardManager       *raftShardManager
	defaultTimeout         time.Duration
	extendedDefaultTimeout time.Duration
}

// SetupTest represents a remote Pleiades host
func (t *RaftShardGrpcAdapterTestSuite) SetupTest() {
	t.logger = utils.NewTestLogger(t.T())

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	t.srv = grpc.NewServer()

	t.testShardId = rand.Uint64()
	t.testClusterConfig = buildTestShardConfig(t.T())
	t.defaultTimeout = 300 * time.Millisecond
	t.extendedDefaultTimeout = 500 * time.Millisecond

	t.testShardManager = newShardManager(buildTestNodeHost(t.T()), t.logger)

	t.adapter = &raftShardGrpcAdapter{
		logger:       t.logger,
		shardManager: t.testShardManager,
	}

	err := t.adapter.shardManager.NewShard(t.testShardId, t.testClusterConfig.NodeID, testStateMachineType, utils.Timeout(t.defaultTimeout))
	t.Require().NoError(err, "there must not be an error when starting the test shard")
	utils.Wait(t.defaultTimeout)

	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := t.testShardManager.nh.SyncGetSession(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when starting the setup statemachine")

	for i := 0; i < 5; i++ {
		proposeCtx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := t.testShardManager.nh.SyncPropose(proposeCtx, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a test message during setup")
		cs.ProposalCompleted()
	}

	RegisterShardManagerServer(t.srv, t.adapter)
	go func() {
		if err := t.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	t.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (t *RaftShardGrpcAdapterTestSuite) TearDownTest() {
	// safely close things.
	t.conn.Close()
	t.srv.Stop()

	// clear out the values
	t.srv = nil
	t.adapter = nil
	t.conn = nil
}

func (t *RaftShardGrpcAdapterTestSuite) TestAddReplica() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")
}

func (t *RaftShardGrpcAdapterTestSuite) TestAddReplicaObserver() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId
	clusterConfig.IsObserver = true

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplicaObserver(context.Background(), &raft.AddReplicaObserverRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(1, len(members.Observers), "there must be one observer in the shard")
}

func (t *RaftShardGrpcAdapterTestSuite) TestAddReplicaWitness() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId
	clusterConfig.IsWitness = true

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplicaWitness(context.Background(), &raft.AddReplicaWitnessRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(1, len(members.Witnesses), "there must be one observer in the shard")
}

func (t *RaftShardGrpcAdapterTestSuite) TestGetLeaderId() {
	client := NewShardManagerClient(t.conn)
	resp, err := client.GetLeaderId(context.Background(), &raft.GetLeaderIdRequest{
		ReplicaId: t.testClusterConfig.NodeID,
		ShardId:   t.testShardId,
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when getting the leader id")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().NotEmpty(resp.GetLeader(), "the leader is must not be empty")
	t.Require().True(resp.GetAvailable(), "the leader information must be available")
	t.Require().Equal(t.testClusterConfig.NodeID, resp.GetLeader())
}

// todo (sienna): this should check for both observers and witnesses at a later point
func (t *RaftShardGrpcAdapterTestSuite) TestGetShardMembers() {
	client := NewShardManagerClient(t.conn)
	resp, err := client.GetShardMembers(context.Background(), &raft.GetShardMembersRequest{
		ShardId: t.testShardId,
	})
	t.Require().NoError(err, "there must not be an error when getting the leader id")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().NotEmpty(resp.GetConfigChangeId(), "the config change id must not be empty")
	t.Require().Equal(1, len(resp.GetReplicas()), "there must be at least one replica in the cluster")
	t.Require().Equal(0, len(resp.GetObservers()), "there must be no observers in the cluster")
	t.Require().Equal(0, len(resp.GetWitnesses()), "there must be no witnesses in the cluster")
	t.Require().Equal(0, len(resp.GetRemoved()), "there must be no removed replicas in the cluster")
}

func (t *RaftShardGrpcAdapterTestSuite) TestNewShard() {
	client := NewShardManagerClient(t.conn)
	_, err := client.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   t.testShardId + 1,
		ReplicaId: t.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  t.testShardManager.nh.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
}

func (t *RaftShardGrpcAdapterTestSuite) TestRemoveData() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.StopReplica(context.Background(), &raft.StopReplicaRequest{
		ShardId: t.testShardId,
	})
	t.Require().NoError(err, "there must not be an error when stopping the replica")

	_, err = client.RemoveData(context.Background(), &raft.RemoveDataRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
	})
	t.Require().NoError(err, "there must not be an error when removing replica data")
}

func (t *RaftShardGrpcAdapterTestSuite) TestRemoveReplica() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.RemoveReplica(context.Background(), &raft.DeleteReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when deleting a replica")

	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err = t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(1, len(members.Nodes), "there must be one replica in the cluster after deletion")
}

func (t *RaftShardGrpcAdapterTestSuite) TestStartReplica() {

	currentTestShard := rand.Uint64()

	// generate a new "local" shard manager
	localNodeHost := buildTestNodeHost(t.T())
	shardConfig := buildTestShardConfig(t.T())
	shardConfig.ClusterID = currentTestShard
	localShardManager := newShardManager(localNodeHost, t.logger)

	// and then wire it to a server.
	listener := bufconn.Listen(1024*1024)
	ctx := context.Background()
	srv := grpc.NewServer()
	adapter := &raftShardGrpcAdapter{
		logger:       t.logger,
		shardManager: localShardManager,
	}
	RegisterShardManagerServer(srv, adapter)
	go func() {
		if err := srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	// the "remote" host is creating a new shard
	remoteClient := NewShardManagerClient(t.conn)
	_, err := remoteClient.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   currentTestShard,
		ReplicaId: t.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  t.testShardManager.nh.RaftAddress(), // remote host address (itself)
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
	utils.Wait(t.defaultTimeout)

	// the grpc side of the "local" host
	localClient := NewShardManagerClient(conn)

	// we're telling the remote host to add our "local" host as a replica
	_, err = remoteClient.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Hostname:  localNodeHost.RaftAddress(), // local host address
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a new replica")

	// now tell the "local" host to start the replica
	_, err = localClient.StartReplica(context.Background(), &raft.StartReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
	})
	t.Require().NoError(err, "there must not be an error when starting a replica")
	utils.Wait(t.defaultTimeout)

	// fetch the members from the "local" perspective, to make sure everything is okay.
	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := localShardManager.nh.SyncGetClusterMembership(ctx, currentTestShard)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")
}

func (t *RaftShardGrpcAdapterTestSuite) TestStartReplicaObserver() {

	currentTestShard := rand.Uint64()

	// generate a new "local" shard manager
	localNodeHost := buildTestNodeHost(t.T())
	shardConfig := buildTestShardConfig(t.T())
	shardConfig.ClusterID = currentTestShard
	localShardManager := newShardManager(localNodeHost, t.logger)

	// and then wire it to a server.
	listener := bufconn.Listen(1024*1024)
	ctx := context.Background()
	srv := grpc.NewServer()
	adapter := &raftShardGrpcAdapter{
		logger:       t.logger,
		shardManager: localShardManager,
	}
	RegisterShardManagerServer(srv, adapter)
	go func() {
		if err := srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	// the "remote" host is creating a new shard
	remoteClient := NewShardManagerClient(t.conn)
	_, err := remoteClient.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   currentTestShard,
		ReplicaId: t.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  t.testShardManager.nh.RaftAddress(), // remote host address (itself)
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
	utils.Wait(t.defaultTimeout)

	// the grpc side of the "local" host
	localClient := NewShardManagerClient(conn)

	// we're telling the remote host to add our "local" host as a replica
	_, err = remoteClient.AddReplicaObserver(context.Background(), &raft.AddReplicaObserverRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Hostname:  localNodeHost.RaftAddress(), // local host address
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a new replica")
	utils.Wait(t.defaultTimeout)

	// now tell the "local" host to start the replica
	_, err = localClient.StartReplicaObserver(context.Background(), &raft.StartReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
	})
	t.Require().NoError(err, "there must not be an error when starting a replica")
	utils.Wait(t.defaultTimeout)

	// fetch the members from the "local" perspective, to make sure everything is okay.
	ctx, _ = context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := localShardManager.nh.SyncGetClusterMembership(ctx, currentTestShard)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(1, len(members.Observers), "there must be one observer in the shard")
}

func (t *RaftShardGrpcAdapterTestSuite) TestStopReplica() {

	testNodeHost := buildTestNodeHost(t.T())

	clusterConfig := buildTestShardConfig(t.T())
	clusterConfig.ClusterID = t.testShardId

	client := NewShardManagerClient(t.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   t.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when adding a replica")
	utils.Wait(t.defaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	t.Require().NoError(err, "there must not be an error when starting a test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	members, err := t.testShardManager.nh.SyncGetClusterMembership(ctx, t.testShardId)
	t.Require().NoError(err, "there must not be an error when fetching shard members")
	t.Require().NotNil(members, "the members response must not be nil")
	t.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.StopReplica(context.Background(), &raft.StopReplicaRequest{
		ShardId: t.testShardId,
	})
	t.Require().NoError(err, "there must not be an error when stopping the replica")
}
