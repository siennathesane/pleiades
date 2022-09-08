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
	suite.Run(t, new(RaftShardGrpcAdapterTestSuite))
}

type RaftShardGrpcAdapterTestSuite struct {
	suite.Suite
	logger zerolog.Logger
	adapter *raftControlGrpcAdapter
	conn   *grpc.ClientConn
	srv   *grpc.Server
	testShardId uint64
	testClusterConfig dconfig.Config
	testShardManager *shardManager
	defaultTimeout         time.Duration
	extendedDefaultTimeout time.Duration
}

func (r *RaftShardGrpcAdapterTestSuite) SetupTest() {
	r.logger = utils.NewTestLogger(r.T())

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	r.srv = grpc.NewServer()

	r.testShardId = rand.Uint64()
	r.testClusterConfig = buildTestClusterConfig(r.T())
	r.defaultTimeout = 3000 * time.Millisecond
	r.extendedDefaultTimeout = 5000 * time.Millisecond

	r.testShardManager = newShardManager(buildTestNodeHost(r.T()), r.logger)

	r.adapter = &raftControlGrpcAdapter{
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

	clusterConfig := buildTestClusterConfig(r.T())
	others := make(map[uint64]string)
	others[clusterConfig.ClusterID] = testNodeHost.RaftAddress()

	client := NewShardManagerClient(r.conn)
	_, err := client.AddReplica(context.Background(), &raft.AddReplicaRequest{
		ReplicaId: clusterConfig.NodeID,
		ShardId:   r.testShardId,
		Type:      raft.StateMachineType_TEST,
		Hostname:  testNodeHost.RaftAddress(),
		Timeout:   int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when adding a replica")
	time.Sleep(r.extendedDefaultTimeout)

	err = testNodeHost.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	r.Require().NoError(err, "there must not be an error when starting a test state machine")
	time.Sleep(r.extendedDefaultTimeout)
}

func (r *RaftShardGrpcAdapterTestSuite) TestAddShardObserver() {}
func (r *RaftShardGrpcAdapterTestSuite) TestAddShardWitness()  {}
func (r *RaftShardGrpcAdapterTestSuite) TestDeleteReplica()    {}
func (r *RaftShardGrpcAdapterTestSuite) TestGetLeaderId()      {}
func (r *RaftShardGrpcAdapterTestSuite) TestGetShardMembers()  {}
func (r *RaftShardGrpcAdapterTestSuite) TestNewShard()         {}
func (r *RaftShardGrpcAdapterTestSuite) TestRemoveData()       {}
func (r *RaftShardGrpcAdapterTestSuite) TestStopReplica()      {}
