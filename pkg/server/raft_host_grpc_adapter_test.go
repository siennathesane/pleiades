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

	"github.com/mxplusb/pleiades/api/v1/raft"
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
func (r *raftHostGrpcAdapterTestSuite) SetupTest() {
	r.logger = utils.NewTestLogger(r.T())
	r.defaultTimeout = 3000 * time.Millisecond

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	r.srv = grpc.NewServer()

	r.nh = buildTestNodeHost(r.T())
	r.rh = &raftHost{
		nh:     r.nh,
		logger: r.logger,
	}

	RegisterRaftHostServer(r.srv, &raftHostGrpcAdapter{
		logger: r.logger,
		host:   r.rh,
	})

	go func() {
		if err := r.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (r *raftHostGrpcAdapterTestSuite) TearDownTest() {
	// safely close things.
	r.conn.Close()
	r.srv.Stop()

	// clear out the values
	r.srv = nil
	r.conn = nil
	r.rh = nil
	r.nh = nil
}

func (r *raftHostGrpcAdapterTestSuite) TestCompact() {
	if testing.Short() {
		r.T().Skipf("skipping")
	}

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(3000*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	cs, err := r.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	r.Require().NoError(err, "there must not be an error when fetching the client session")
	r.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	for i := 0; i < 25; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		_, err := r.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		r.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	client := NewRaftHostClient(r.conn)

	// todo (sienna): figure out why it's being rejected.
	_, err = client.Compact(context.Background(), &raft.CompactRequest{
		ReplicaId: shardConfig.NodeID,
		ShardId:   shardConfig.ClusterID,
	})
	r.Require().Error(err, "the request for log compaction must be rejected")
}

func (r *raftHostGrpcAdapterTestSuite) TestGetHostConfig() {
	client := NewRaftHostClient(r.conn)
	resp, err := client.GetHostConfig(context.Background(), &raft.GetHostConfigRequest{})
	r.Require().NoError(err, "there must not be an error getting the host config")
	r.Require().NotEmpty(resp, "the response must not be empty")
}

func (r *raftHostGrpcAdapterTestSuite) TestSnapshot() {
	if testing.Short() {
		r.T().Skipf("skipping")
	}

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(3000*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	cs, err := r.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	r.Require().NoError(err, "there must not be an error when fetching the client session")
	r.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	for i := 0; i < 25; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		_, err := r.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		r.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	client := NewRaftHostClient(r.conn)

	// todo (sienna): figure out why it's failing to write to disk
	_, err = client.Snapshot(context.Background(), &raft.SnapshotRequest{
		ShardId:   shardConfig.ClusterID,
		Timeout: int64(r.defaultTimeout),
	})
	r.Require().NoError(err, "there must not be an error when trying to create a snapshot")
}

func (r *raftHostGrpcAdapterTestSuite) TestStop() {
	client := NewRaftHostClient(r.conn)
	_, err := client.Stop(context.Background(), &raft.StopRequest{})
	r.Require().NoError(err, "there must not be an error when trying to stop the host")
}

