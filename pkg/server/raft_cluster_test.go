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
	"testing"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestClusterManager(t *testing.T) {
	suite.Run(t, new(ClusterManagerTests))
}

type ClusterManagerTests struct {
	suite.Suite
	logger zerolog.Logger
	node   *Node
}

func (cmt *ClusterManagerTests) SetupSuite() {
	if testing.Short() {
		cmt.T().Skipf("skipping")
	}

	cmt.logger = utils.NewTestLogger(cmt.T())

	nodeHostConfig := buildTestNodeHostConfig(cmt.T())
	node, err := NewRaftControlNode(nodeHostConfig, cmt.logger)
	cmt.Require().NoError(err, "there must not be an error when starting the node")
	cmt.Require().NotNil(node, "node must not be nil")
	cmt.node = node
}

func (cmt *ClusterManagerTests) TestStartCluster() {
	//testClusterId := uint64(0)
	//firstNodeClusterConfig := buildTestClusterConfig(cmt.T())
	//testClusterId = firstNodeClusterConfig.ClusterID
	//
	//nodeClusters := make(map[uint64]string)
	//nodeClusters[firstNodeClusterConfig.NodeID] = cmt.node.RaftAddress()
	//
	//cm := newClusterManager(cmt.node.nh, cmt.logger)
	//
	//err := cm.StartCluster(nodeClusters, false, newTestStateMachine, firstNodeClusterConfig)
	//cmt.Require().NoError(err, "there must not be an error when starting the test state machine")
	//time.Sleep(5000 * time.Millisecond)
	//
	//sm := newSessionManager(cmt.node.nh, cmt.logger)
	//cs := sm.GetNoOpSession(testClusterId)
	//
	//proposeContext, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	//_, err = cmt.node.nh.SyncPropose(proposeContext, cs, []byte("test-message"))
	//cmt.Require().NoError(err, "there must not be an error when proposing a new message")
}
