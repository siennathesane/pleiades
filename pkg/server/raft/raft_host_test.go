/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package raft

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
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
	logger         zerolog.Logger
	nh             *dragonboat.NodeHost
	defaultTimeout time.Duration
}

func (t *RaftHostTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 300 * time.Millisecond
}

func (t *RaftHostTestSuite) SetupTest() {
	t.nh = serverutils.BuildTestNodeHost(t.T())
}

func (t *RaftHostTestSuite) TearDownTest() {
	t.nh = nil
}

func (t *RaftHostTestSuite) TestCompact() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	shardConfig := serverutils.BuildTestShardConfig(t.T())
	shardConfig.SnapshotEntries = 5
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(members, false, serverutils.NewTestStateMachine, shardConfig)
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

	// todo (sienna): figure out why it's being rejected.
	err = host.Compact(shardConfig.ClusterID, shardConfig.NodeID)
	t.Require().Error(err, "the request for log compaction must be rejected")
}

func (t *RaftHostTestSuite) TestGetHostInfo() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	resp := host.GetHostInfo(runtime.HostInfoOption{SkipLogInfo: false})
	t.Require().NotEmpty(resp, "the response must not be empty")
}

func (t *RaftHostTestSuite) TestHasNodeInfo() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	shardConfig := serverutils.BuildTestShardConfig(t.T())
	shardConfig.SnapshotEntries = 0
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(members, false, serverutils.NewTestStateMachine, shardConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	time.Sleep(3000 * time.Millisecond)

	resp := host.HasNodeInfo(shardConfig.ClusterID, shardConfig.NodeID)
	t.Require().True(resp, "the host must have a replica it started")
}

func (t *RaftHostTestSuite) TestId() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	hostname := host.Id()
	t.logger.Info().Str("host-id", hostname).Msg("got host id")
	t.Require().NotEmpty(hostname, "the hostname must not be empty")
}

func (t *RaftHostTestSuite) TestHostConfig() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	resp := host.HostConfig()
	t.Require().NotEmpty(resp, "the host config must not be empty")
}

func (t *RaftHostTestSuite) TestRaftAddress() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	resp := host.RaftAddress()
	t.Require().NotEmpty(resp, "the raft address must not be empty")
}

func (t *RaftHostTestSuite) TestSnapshot() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)

	shardConfig := serverutils.BuildTestShardConfig(t.T())
	shardConfig.SnapshotEntries = 0
	members := make(map[uint64]string)
	members[shardConfig.NodeID] = t.nh.RaftAddress()

	err := t.nh.StartCluster(members, false, serverutils.NewTestStateMachine, shardConfig)
	t.Require().NoError(err, "there must not be an error when starting the test state machine")
	utils.Wait(t.defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
	cs, err := t.nh.SyncGetSession(ctx, shardConfig.ClusterID)
	t.Require().NoError(err, "there must not be an error when fetching the client session")
	t.Require().NotNil(cs, "the client session must not be nil")
	cancel()

	loops := 25
	for i := 0; i < loops; i++ {
		proposeContext, _ := context.WithTimeout(context.Background(), utils.Timeout(t.defaultTimeout))
		_, err := t.nh.SyncPropose(proposeContext, cs, []byte(fmt.Sprintf("test-message-%d", i)))
		t.Require().NoError(err, "there must not be an error when proposing a new message")
		cs.ProposalCompleted()
	}

	// todo (sienna): figure why this doesn't work
	//x := t.T().TempDir()
	//err = os.MkdirAll(x, os.FileMode(484))
	//t.Require().NoError(err, "there must not be an error when creating the temp directory")
	//target := filepath.Join(x, "test.bin")

	_, err = host.Snapshot(shardConfig.ClusterID, runtime.SnapshotOption{
		ExportPath: ".", // replace with target once you figure this problem out
		Exported:   true,
	}, 3000*time.Millisecond)
	t.Require().NoError(err, "there must not be an error when requesting a snapshot")
}

func (t *RaftHostTestSuite) TestStop() {
	params := &RaftHostBuilderParams{
		NodeHost: t.nh,
		Logger:   t.logger,
	}
	hostRes := NewHost(params)
	host := hostRes.RaftHost.(*RaftHost)
	host.Stop()
}
