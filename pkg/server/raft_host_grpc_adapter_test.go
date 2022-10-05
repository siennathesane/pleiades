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
	"net"
	"testing"
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestRaftHostGrpcAdapter(t *testing.T) {
	suite.Run(t, new(raftHostGrpcAdapterTestSuite))
}

type raftHostGrpcAdapterTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	conn           *grpc.ClientConn
	srv            *grpc.Server
	nh             *dragonboat.NodeHost
	rh             *raftHost
	defaultTimeout time.Duration
}

// SetupTest represents a remote Pleiades host
func (t *raftHostGrpcAdapterTestSuite) SetupTest() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 300 * time.Millisecond

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	t.srv = grpc.NewServer()

	t.nh = utils.BuildTestNodeHost(t.T())
	t.rh = &raftHost{
		nh:     t.nh,
		logger: t.logger,
	}

	raftv1.RegisterHostServiceServer(t.srv, &raftHostGrpcAdapter{
		logger: t.logger,
		host:   t.rh,
	})

	go func() {
		if err := t.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	t.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (t *raftHostGrpcAdapterTestSuite) TearDownTest() {
	// safely close things.
	t.conn.Close()
	t.srv.Stop()

	// clear out the values
	t.srv = nil
	t.conn = nil
	t.rh = nil
	t.nh = nil
}

func (t *raftHostGrpcAdapterTestSuite) TestCompact() {
	if testing.Short() {
		t.T().Skipf("skipping")
	}

	shardConfig := utils.BuildTestShardConfig(t.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(members, false, utils.NewTestStateMachine, shardConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := t.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	t.Require().NoError(err, "there must not be an error when fetching the client session")
	t.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	for i := 0; i < 25; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := t.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	client := raftv1.NewHostServiceClient(t.conn)

	// todo (sienna): figure out why it's being rejected.
	_, err = client.Compact(context.Background(), &raftv1.CompactRequest{
		ReplicaId: shardConfig.NodeID,
		ShardId:   shardConfig.ClusterID,
	})
	t.Require().Error(err, "the request for log compaction must be rejected")
}

func (t *raftHostGrpcAdapterTestSuite) TestGetHostConfig() {
	client := raftv1.NewHostServiceClient(t.conn)
	resp, err := client.GetHostConfig(context.Background(), &raftv1.GetHostConfigRequest{})
	t.Require().NoError(err, "there must not be an error getting the host config")
	t.Require().NotEmpty(resp, "the response must not be empty")
}

func (t *raftHostGrpcAdapterTestSuite) TestSnapshot() {
	if testing.Short() {
		t.T().Skipf("skipping")
	}

	shardConfig := utils.BuildTestShardConfig(t.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(members, false, utils.NewTestStateMachine, shardConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := t.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	t.Require().NoError(err, "there must not be an error when fetching the client session")
	t.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	for i := 0; i < 25; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := t.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	client := raftv1.NewHostServiceClient(t.conn)

	// todo (sienna): figure out why it's failing to write to disk
	_, err = client.Snapshot(context.Background(), &raftv1.SnapshotRequest{
		ShardId:   shardConfig.ClusterID,
		Timeout: int64(t.defaultTimeout),
	})
	t.Require().NoError(err, "there must not be an error when trying to create a snapshot")
}

func (t *raftHostGrpcAdapterTestSuite) TestStop() {
	client := raftv1.NewHostServiceClient(t.conn)
	_, err := client.Stop(context.Background(), &raftv1.StopRequest{})
	t.Require().NoError(err, "there must not be an error when trying to stop the host")
}
