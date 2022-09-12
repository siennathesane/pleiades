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
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftHost(t *testing.T) {
	suite.Run(t, new(RaftHostTestSuite))
}

type RaftHostTestSuite struct {
	suite.Suite
	logger zerolog.Logger
	nh *dragonboat.NodeHost
	defaultTimeout time.Duration
}

func (r *RaftHostTestSuite) SetupSuite() {
	r.logger = utils.NewTestLogger(r.T())
	r.defaultTimeout = 300*time.Millisecond
}

func (r *RaftHostTestSuite) SetupTest() {
	r.nh = buildTestNodeHost(r.T())
}

func (r *RaftHostTestSuite) TearDownTest() {
	r.nh = nil
}

func (r *RaftHostTestSuite) TestCompact() {
	host := newRaftHost(r.nh, r.logger)

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(r.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), r.defaultTimeout)
	cs, err := r.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	r.Require().NoError(err, "there must not be an error when fetching the client session")
	r.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	for i := 0; i < 25; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
		_, err := r.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		r.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	// todo (sienna): figure out why it's being rejected.
	err = host.Compact(shardConfig.ClusterID, shardConfig.NodeID)
	r.Require().Error(err, "the request for log compaction must be rejected")
}

func (r *RaftHostTestSuite) TestGetHostInfo() {
	host := newRaftHost(r.nh, r.logger)

	resp := host.GetHostInfo(HostInfoOption{SkipLogInfo: false})
	r.Require().NotEmpty(resp, "the response must not be empty")
}

func (r *RaftHostTestSuite) TestHasNodeInfo() {
	host := newRaftHost(r.nh, r.logger)

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 0
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(3000*time.Millisecond)

	resp := host.HasNodeInfo(shardConfig.ClusterID, shardConfig.NodeID)
	r.Require().True(resp, "the host must have a replica it started")
}

func (r *RaftHostTestSuite) TestId() {
	host := newRaftHost(r.nh, r.logger)

	hostname := host.Id()
	r.logger.Info().Str("host-id", hostname).Msg("got host id")
	r.Require().NotEmpty(hostname, "the hostname must not be empty")
}

func (r *RaftHostTestSuite) TestHostConfig() {
	host := newRaftHost(r.nh, r.logger)

	resp := host.HostConfig()
	r.Require().NotEmpty(resp, "the host config must not be empty")
}

func (r *RaftHostTestSuite) TestRaftAddress() {
	host := newRaftHost(r.nh, r.logger)

	resp := host.RaftAddress()
	r.Require().NotEmpty(resp, "the raft address must not be empty")
}

func (r *RaftHostTestSuite) TestSnapshot() {
	host := newRaftHost(r.nh, r.logger)

	shardConfig := buildTestShardConfig(r.T())
	shardConfig.SnapshotEntries = 0
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = r.nh.RaftAddress()

	err := r.nh.StartCluster(members, false, newTestStateMachine, shardConfig)
	r.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(r.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), r.defaultTimeout)
	cs, err := r.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	r.Require().NoError(err, "there must not be an error when fetching the client session")
	r.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	loops := 25
	for i := 0; i < loops; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), r.defaultTimeout)
		_, err := r.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		r.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	// todo (sienna): figure why this doesn't work
	//x := r.T().TempDir()
	//err = os.MkdirAll(x, os.FileMode(484))
	//r.Require().NoError(err, "there must not be an error when creating the temp directory")
	//target := filepath.Join(x, "test.bin")

	_, err = host.Snapshot(shardConfig.ClusterID, SnapshotOption{
		ExportPath:                 ".", // replace with target once you figure this problem out
		Exported:                   true,
	}, 3000*time.Millisecond)
	r.Require().NoError(err, "there must not be an error when requesting a snapshot")
}

func (r *RaftHostTestSuite) TestStop() {
	host := newRaftHost(r.nh, r.logger)
	host.Stop()
}
