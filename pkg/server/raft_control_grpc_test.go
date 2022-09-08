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
	"net"

	"github.com/mxplusb/pleiades/pkg/api/v1/raft"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type RaftControlGrpcAdapterTestSuite struct {
	suite.Suite
	logger zerolog.Logger
	adapter *raftControlGrpcAdapter
	conn   *grpc.ClientConn
	srv   *grpc.Server
}

func (suite *RaftControlGrpcAdapterTestSuite) SetupSuite() {
	suite.logger = utils.NewTestLogger(suite.T())

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	suite.srv = grpc.NewServer()

	suite.adapter = &raftControlGrpcAdapter{
		logger:         suite.logger,
		clusterManager: newShardManager(buildTestNodeHost(suite.T()), suite.logger),
	}

	RegisterShardManagerServer(suite.srv, suite.adapter)
	go func() {
		if err := suite.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	suite.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (suite *RaftControlGrpcAdapterTestSuite) TearDownSuite() {
	suite.conn.Close()
	suite.srv.Stop()
}

// AddReplica(cfg IClusterConfig, newHost string, timeout time.Duration) error
//	AddShardObserver(cfg IClusterConfig, newHost string, timeout time.Duration) error
//	AddShardWitness(cfg IClusterConfig, newHost string, timeout time.Duration) error
//	DeleteReplica(cfg IClusterConfig, timeout time.Duration) error
//	GetLeaderId(shardId uint64) (leader uint64, ok bool, err error)
//	GetShardMembers(shardId uint64) (*MembershipEntry, error)
//	NewShard(cfg IClusterConfig) error
//	RemoveData(shardId, replicaId uint64) error
//	StopReplica(shardId uint64) (*OperationResult, error)

func (suite *RaftControlGrpcAdapterTestSuite) TestAddReplica() {
	testNodeHost := buildTestNodeHost(suite.T())

	clusterConfig := buildTestClusterConfig(suite.T())
	others := make(map[uint64]string)
	others[clusterConfig.ClusterID] = testNodeHost.RaftAddress()

	err := testNodeHost.StartCluster(others, false, newTestStateMachine, clusterConfig)
	suite.Require().NoError(err, "there must not be an error when starting a test state machine")

	//suite.adapter.shardManager.

	client := NewShardManagerClient(suite.conn)
	_, err = client.AddReplica(context.Background(), &raft.AddReplicaRequest{})
}