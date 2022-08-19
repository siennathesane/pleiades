/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"testing"

	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftControl(t *testing.T) {
	suite.Run(t, new(RaftControlTests))
}

type RaftControlTests struct {
	suite.Suite
	logger zerolog.Logger
}

func (rct *RaftControlTests) SetupSuite() {
	rct.logger = conf.NewRootLogger()
}

func (rct *RaftControlTests) TestRaftAddress() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	addr := node.RaftAddress()
	rct.Require().NotNil(addr, "the raft address must not be nil")
	rct.T().Logf("raft addr: %s", addr)
}

func (rct *RaftControlTests) TestId() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	id := node.ID()
	rct.Require().NotEmpty(id, "the node id must not be empty")
	rct.T().Logf("node id: %s", id)
}

func (rct *RaftControlTests) TestGetNodeUser() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	nodeUser, err := node.GetNodeUser(clusterConfig.ClusterID)
	rct.Require().NoError(err, "there must not be an error when fetching a node user")
	rct.Require().NotNil(nodeUser, "the node user must not be nil")

	rct.Require().Equal(clusterConfig.ClusterID, nodeUser.ClusterID(), "the nodeUser.ClusterID must be equal to the configured cluster ID")
}

func (rct *RaftControlTests) TestGetId() {
	node, err := NewRaftControlNode(buildTestNodeHostConfig(rct.T()), rct.logger)
	rct.Require().NoError(err, "there must not be an error when starting a new control node")
	rct.Require().NotNil(node, "the node must not be nil")

	clusterConfig := buildTestClusterConfig(rct.T())

	err = node.nh.StartCluster(nil, true, newTestStateMachine, clusterConfig)
	rct.Require().NoError(err, "there must not be an error when starting the test state machine")

	leader, ok, err := node.GetLeaderID(clusterConfig.ClusterID)
	rct.Require().NoError(err, "there must not be an error when fetching the leader")

	// until there's a cluster in place, it's okay to skip the rest of it since we're really just validating it works
	// todo (sienna): implement clustering logic for this test
	if !ok && leader == 0 {
		rct.T().Skipf("skipping due to unimplemented cluster functionality for test")
	}

	rct.Assert().True(ok, "the leader information should be available")
	rct.Assert().NotEmpty(leader, "the leader is should not be 0")
}
