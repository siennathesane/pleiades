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
	"time"

	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/rs/zerolog"
)

var (
	_ INodeHost = (*Node)(nil)
)

func init() {
	dlog.SetLoggerFactory(conf.DragonboatLoggerFactory)
}

type Node struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost

	notifyOnCommit bool
	clusterManager IShardManager
	sessionManager ITransactionManager
	storeManager   IStore
}

// NewRaftControlNode creates a new Node instance.
func NewRaftControlNode(nodeHostConfig dconfig.NodeHostConfig, logger zerolog.Logger) (*Node, error) {
	l := logger.With().Str("component", "node").Logger()

	nh, err := dragonboat.NewNodeHost(nodeHostConfig)
	if err != nil {
		l.Error().Err(err).Msg("failed to create node host")
		return nil, err
	}

	node := &Node{logger: l, nh: nh}

	return node, nil
}

// NewOrGetClusterManager creates a new IShardManager instance or gets the existing one.
func (n *Node) NewOrGetClusterManager() (IShardManager, error) {
	if n.clusterManager == nil {
		n.clusterManager = newClusterManager(n.nh,n.logger)
	}
	return n.clusterManager, nil
}

// NewOrGetSessionManager creates a new ITransactionManager instance or gets the existing one.
func (n *Node) NewOrGetSessionManager() (ITransactionManager, error) {
	if n.sessionManager == nil {
		n.sessionManager = newSessionManager(n.nh, n.logger)
	}
	return n.sessionManager, nil
}

// NewOrGetStoreManager creates a new IStore instance or gets the existing one.
func (n *Node) NewOrGetStoreManager() (IStore, error) {
	if n.storeManager == nil {
		n.storeManager = newStoreManager(n.logger, n.nh)
	}
	return n.storeManager, nil
}

func (n *Node) GetLeaderID(clusterID uint64) (uint64, bool, error) {
	return n.nh.GetLeaderID(clusterID)
}

func (n *Node) GetNodeUser(clusterID uint64) (dragonboat.INodeUser, error) {
	return n.nh.GetNodeUser(clusterID)
}

func (n *Node) ID() string {
	return n.nh.ID()
}

func (n *Node) RaftAddress() string {
	return n.nh.RaftAddress()
}

func (n *Node) RemoveData(clusterID uint64, nodeID uint64) error {
	return n.nh.RemoveData(clusterID, nodeID)
}

func (n *Node) AddNode(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	return n.nh.RequestAddNode(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) AddObserver(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	return n.nh.RequestAddObserver(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) AddWitness(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	return n.nh.RequestAddWitness(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) Compact(clusterID uint64, nodeID uint64) (*dragonboat.SysOpState, error) {
	return n.nh.RequestCompaction(clusterID, nodeID)
}

func (n *Node) DeleteNode(clusterID uint64, nodeID uint64, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	return n.nh.RequestDeleteNode(clusterID, nodeID, configChangeIndex, timeout)
}

func (n *Node) LeaderTransfer(clusterID uint64, targetNodeID uint64) error {
	return n.nh.RequestLeaderTransfer(clusterID, targetNodeID)
}

func (n *Node) Snapshot(clusterID uint64, opt dragonboat.SnapshotOption, timeout time.Duration) (*dragonboat.RequestState, error) {
	return n.nh.RequestSnapshot(clusterID, opt, timeout)
}

func (n *Node) Stop() {
	n.nh.Stop()
}

func (n *Node) StopNode(clusterID uint64, nodeID uint64) error {
	return n.nh.StopNode(clusterID, nodeID)
}
