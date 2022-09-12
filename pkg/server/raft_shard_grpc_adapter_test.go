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
	testShardManager       *shardManager
	defaultTimeout         time.Duration
	extendedDefaultTimeout time.Duration
}

// SetupTest represents a remote Pleiades host
func (r *RaftShardGrpcAdapterTestSuite) SetupTest() {
	r.logger = utils.NewTestLogger(r.T())

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	r.srv = grpc.NewServer()

	r.testShardId = rand.Uint64()
	r.testClusterConfig = buildTestShardConfig(r.T())
	r.defaultTimeout = 300 * time.Millisecond
	r.extendedDefaultTimeout = 500 * time.Millisecond

	r.testShardManager = newShardManager(buildTestNodeHost(r.T()), r.logger)

	r.adapter = &raftShardGrpcAdapter{
		logger:         r.logger,
		clusterManager: r.testShardManager,
	}

	err := r.adapter.clusterManager.NewShard(r.testShardId, r.testClusterConfig.NodeID, testStateMachineType, r.defaultTimeout)
	r.Require().NoError(err, "there must not be an error when starting the test shard")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ = context.WithTimeout(context.Background(), r.defaultTimeout)
	cs, err := r.testShardManager.nh.SyncGetSession(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when starting the setup statemachine")

	for i := 0; i < 5; i++ {
		proposeCtx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
		_, err := r.testShardManager.nh.SyncPropose(proposeCtx, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		r.Require().NoError(err, "there must not be an error when proposing a test message during setup")
		cs.ProposalCompleted()
	}

	RegisterShardManagerServer(r.srv, r.adapter)
	go func() {
		if err := r.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (r *RaftShardGrpcAdapterTestSuite) TearDownTest() {
	// safely close things.
	r.conn.Close()
	r.srv.Stop()

	// clear out the values
	r.srv = nil
	r.adapter = nil
	r.conn = nil
}

func (r *RaftShardGrpcAdapterTestSuite) TestAddReplica() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")
}

func (r *RaftShardGrpcAdapterTestSuite) TestAddReplicaObserver() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId
	clusterConfig.IsObserver = true

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplicaObserver(context.Background(), &raft.AddReplicaObserverRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(1, len(members.Observers), "there must be one observer in the shard")
}

func (r *RaftShardGrpcAdapterTestSuite) TestAddReplicaWitness() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId
	clusterConfig.IsWitness = true

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplicaWitness(context.Background(), &raft.AddReplicaWitnessRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(1, len(members.Witnesses), "there must be one observer in the shard")
}

func (r *RaftShardGrpcAdapterTestSuite) TestGetLeaderId() {
	client := NewShardManagerClient(r.conn)
	resp, err := client.GetLeaderId(context.Background(), &raft.GetLeaderIdRequest{
		ReplicaId: r.testClusterConfig.NodeID,
		ShardId:   r.testShardId,
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when getting the leader id")
	r.Require().NotNil(resp, "the response must not be nil")
	r.Require().NotEmpty(resp.GetLeader(), "the leader is must not be empty")
	r.Require().True(resp.GetAvailable(), "the leader information must be available")
	r.Require().Equal(r.testClusterConfig.NodeID, resp.GetLeader())
}

// todo (sienna): this should check for both observers and witnesses at a later point
func (r *RaftShardGrpcAdapterTestSuite) TestGetShardMembers() {
	client := NewShardManagerClient(r.conn)
	resp, err := client.GetShardMembers(context.Background(), &raft.GetShardMembersRequest{
		ShardId: r.testShardId,
	})
	r.Require().NoError(err, "there must not be an error when getting the leader id")
	r.Require().NotNil(resp, "the response must not be nil")
	r.Require().NotEmpty(resp.GetConfigChangeId(), "the config change id must not be empty")
	r.Require().Equal(1, len(resp.GetReplicas()), "there must be at least one replica in the cluster")
	r.Require().Equal(0, len(resp.GetObservers()), "there must be no observers in the cluster")
	r.Require().Equal(0, len(resp.GetWitnesses()), "there must be no witnesses in the cluster")
	r.Require().Equal(0, len(resp.GetRemoved()), "there must be no removed replicas in the cluster")
}

func (r *RaftShardGrpcAdapterTestSuite) TestNewShard() {
	client := NewShardManagerClient(r.conn)
	_, err := client.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   r.testShardId + 1,
		ReplicaId: r.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  r.testShardManager.nh.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
}

func (r *RaftShardGrpcAdapterTestSuite) TestRemoveData() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.StopReplica(context.Background(), &raft.StopReplicaRequest{
		ShardId: r.testShardId,
	})
	r.Require().NoError(err, "there must not be an error when stopping the replica")

	_, err = client.RemoveData(context.Background(), &raft.RemoveDataRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
	})
	r.Require().NoError(err, "there must not be an error when removing replica data")
}

func (r *RaftShardGrpcAdapterTestSuite) TestRemoveReplica() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.RemoveReplica(context.Background(), &raft.DeleteReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when deleting a replica")

	ctx, _ = context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err = r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(1, len(members.Nodes), "there must be one replica in the cluster after deletion")
}

func (r *RaftShardGrpcAdapterTestSuite) TestStartReplica() {

	currentTestShard := rand.Uint64()

	// generate a new "local" shard manager
	localNodeHost := buildTestNodeHost(r.T())
	shardConfig := buildTestShardConfig(r.T())
	shardConfig.ClusterID = currentTestShard
	localShardManager := newShardManager(localNodeHost, r.logger)

	// and then wire it to a server.
	listener := bufconn.Listen(1024*1024)
	ctx := context.Background()
	srv := grpc.NewServer()
	adapter := &raftShardGrpcAdapter{
		logger:         r.logger,
		clusterManager: localShardManager,
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
	remoteClient := NewShardManagerClient(r.conn)
	_, err := remoteClient.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   currentTestShard,
		ReplicaId: r.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  r.testShardManager.nh.RaftAddress(), // remote host address (itself)
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
	time.Sleep(r.extendedDefaultTimeout)

	// the grpc side of the "local" host
	localClient := NewShardManagerClient(conn)

	// we're telling the remote host to add our "local" host as a replica
	_, err = remoteClient.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Hostname:  localNodeHost.RaftAddress(), // local host address
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a new replica")

	// now tell the "local" host to start the replica
	_, err = localClient.StartReplica(context.Background(), &raft.StartReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
	})
	r.Require().NoError(err, "there must not be an error when starting a replica")
	time.Sleep(r.extendedDefaultTimeout)

	// fetch the members from the "local" perspective, to make sure everything is okay.
	ctx, _ = context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := localShardManager.nh.SyncGetClusterMembership(ctx, currentTestShard)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")
}

func (r *RaftShardGrpcAdapterTestSuite) TestStartReplicaObserver() {

	currentTestShard := rand.Uint64()

	// generate a new "local" shard manager
	localNodeHost := buildTestNodeHost(r.T())
	shardConfig := buildTestShardConfig(r.T())
	shardConfig.ClusterID = currentTestShard
	localShardManager := newShardManager(localNodeHost, r.logger)

	// and then wire it to a server.
	listener := bufconn.Listen(1024*1024)
	ctx := context.Background()
	srv := grpc.NewServer()
	adapter := &raftShardGrpcAdapter{
		logger:         r.logger,
		clusterManager: localShardManager,
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
	remoteClient := NewShardManagerClient(r.conn)
	_, err := remoteClient.NewShard(context.Background(), &raft.NewShardRequest{
		ShardId:   currentTestShard,
		ReplicaId: r.testClusterConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
		Hostname:  r.testShardManager.nh.RaftAddress(), // remote host address (itself)
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when creating a new test shard on an existing node")
	time.Sleep(r.extendedDefaultTimeout)

	// the grpc side of the "local" host
	localClient := NewShardManagerClient(conn)

	// we're telling the remote host to add our "local" host as a replica
	_, err = remoteClient.AddReplicaObserver(context.Background(), &raft.AddReplicaObserverRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Hostname:  localNodeHost.RaftAddress(), // local host address
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a new replica")
	time.Sleep(r.extendedDefaultTimeout)

	// now tell the "local" host to start the replica
	_, err = localClient.StartReplicaObserver(context.Background(), &raft.StartReplicaRequest{
		ShardId:   currentTestShard,
		ReplicaId: shardConfig.NodeID,
		Type:      raft.StateMachineType_TEST,
	})
	r.Require().NoError(err, "there must not be an error when starting a replica")
	time.Sleep(r.extendedDefaultTimeout)

	// fetch the members from the "local" perspective, to make sure everything is okay.
	ctx, _ = context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := localShardManager.nh.SyncGetClusterMembership(ctx, currentTestShard)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(1, len(members.Observers), "there must be one observer in the shard")
}

func (r *RaftShardGrpcAdapterTestSuite) TestStopReplica() {

	testNodeHost := buildTestNodeHost(r.T())

	clusterConfig := buildTestShardConfig(r.T())
	clusterConfig.ClusterID = r.testShardId

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
	members, err := r.testShardManager.nh.SyncGetClusterMembership(ctx, r.testShardId)
	r.Require().NoError(err, "there must not be an error when fetching shard members")
	r.Require().NotNil(members, "the members response must not be nil")
	r.Require().Equal(2, len(members.Nodes), "there must be two replicas in the cluster")

	_, err = client.StopReplica(context.Background(), &raft.StopReplicaRequest{
		ShardId: r.testShardId,
	})
	r.Require().NoError(err, "there must not be an error when stopping the replica")
}
