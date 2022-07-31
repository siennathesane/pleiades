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
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
)

var (
	_ INodeHost = (*Node)(nil)
)

type Node struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost

	started        bool
	notifyOnCommit bool
	clusterManager ICluster
	sessionManager ISession
	storeManager   IStore
}

// NewNode creates a new Node instance.
func NewNode(conf dconfig.NodeHostConfig, logger zerolog.Logger) (*Node, error) {
	l := logger.With().Str("component", "node").Logger()

	nh, err := dragonboat.NewNodeHost(conf)
	if err != nil {
		l.Error().Err(err).Msg("failed to create node host")
		return nil, err
	}

	node := &Node{logger: l, nh: nh}

	if conf.NotifyCommit {
		node.notifyOnCommit = true
	}

	return node, nil
}

// NewOrGetClusterManager creates a new ICluster instance or gets the existing one.
func (n *Node) NewOrGetClusterManager() (ICluster, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	if n.clusterManager != nil {
		n.clusterManager = newClusterManager(n.logger, n.nh)
	}
	return n.clusterManager, nil
}

// NewOrGetSessionManager creates a new ISession instance or gets the existing one.
func (n *Node) NewOrGetSessionManager() (ISession, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	if n.sessionManager != nil {
		n.sessionManager = newSessionManager(n.logger, n.nh)
	}
	return n.sessionManager, nil
}

// NewOrGetStoreManager creates a new IStore instance or gets the existing one.
func (n *Node) NewOrGetStoreManager() (IStore, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	if n.storeManager != nil {
		n.storeManager = newStoreManager(n.logger, n.nh)
	}
	return n.storeManager, nil
}

func (n *Node) NotifyOnCommit() bool {
	return n.notifyOnCommit
}

func (n *Node) GetLeaderID(clusterID uint64) (uint64, bool, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return 0, false, err
	}

	return n.nh.GetLeaderID(clusterID)
}

func (n *Node) GetNodeUser(clusterID uint64) (dragonboat.INodeUser, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.GetNodeUser(clusterID)
}

func (n *Node) ID() string {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return err.Error()
	}

	return n.nh.ID()
}

func (n *Node) NAReadLocalNode(rs *dragonboat.RequestState, query []byte) ([]byte, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.NAReadLocalNode(rs, query)
}

func (n *Node) RaftAddress() string {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return err.Error()
	}

	return n.nh.RaftAddress()
}

func (n *Node) ReadIndex(clusterID uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.ReadIndex(clusterID, timeout)
}

func (n *Node) ReadLocalNode(rs *dragonboat.RequestState, query interface{}) (interface{}, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.ReadLocalNode(rs, query)
}

func (n *Node) RemoveData(clusterID uint64, nodeID uint64) error {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return err
	}

	return n.nh.RemoveData(clusterID, nodeID)
}

func (n *Node) RequestAddNode(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestAddNode(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) RequestAddObserver(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestAddObserver(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) RequestAddWitness(clusterID uint64, nodeID uint64, target dragonboat.Target, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestAddWitness(clusterID, nodeID, target, configChangeIndex, timeout)
}

func (n *Node) RequestCompaction(clusterID uint64, nodeID uint64) (*dragonboat.SysOpState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestCompaction(clusterID, nodeID)
}

func (n *Node) RequestDeleteNode(clusterID uint64, nodeID uint64, configChangeIndex uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestDeleteNode(clusterID, nodeID, configChangeIndex, timeout)
}

func (n *Node) RequestLeaderTransfer(clusterID uint64, targetNodeID uint64) error {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return err
	}

	return n.nh.RequestLeaderTransfer(clusterID, targetNodeID)
}

func (n *Node) RequestSnapshot(clusterID uint64, opt dragonboat.SnapshotOption, timeout time.Duration) (*dragonboat.RequestState, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.RequestSnapshot(clusterID, opt, timeout)
}

func (n *Node) StaleRead(clusterID uint64, query interface{}) (interface{}, error) {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return nil, err
	}

	return n.nh.StaleRead(clusterID, query)
}

func (n *Node) Stop() {
	if !n.started {
		return
	}
	n.nh.Stop()
}

func (n *Node) StopNode(clusterID uint64, nodeID uint64) error {
	err, ok := n.verifyStarted()
	if !ok || err != nil {
		return err
	}

	return n.nh.StopNode(clusterID, nodeID)
}

func (n *Node) SyncRemoveData(ctx context.Context, clusterID uint64, nodeID uint64) error {
	return n.nh.SyncRemoveData(ctx, clusterID, nodeID)
}

func (n *Node) SyncRequestAddNode(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error {
	return n.nh.SyncRequestAddNode(ctx, clusterID, nodeID, target, configChangeIndex)
}

func (n *Node) SyncRequestAddObserver(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error {
	return n.nh.SyncRequestAddObserver(ctx, clusterID, nodeID, target, configChangeIndex)
}

func (n *Node) SyncRequestAddWitness(ctx context.Context, clusterID uint64, nodeID uint64, target string, configChangeIndex uint64) error {
	return n.nh.SyncRequestAddWitness(ctx, clusterID, nodeID, target, configChangeIndex)
}

func (n *Node) SyncRequestDeleteNode(ctx context.Context, clusterID uint64, nodeID uint64, configChangeIndex uint64) error {
	return n.nh.SyncRequestDeleteNode(ctx, clusterID, nodeID, configChangeIndex)
}

func (n *Node) SyncRequestSnapshot(ctx context.Context, clusterID uint64, opt dragonboat.SnapshotOption) (uint64, error) {
	return n.nh.SyncRequestSnapshot(ctx, clusterID, opt)
}

// verifyStarted checks if the node host is started and okay
func (n *Node) verifyStarted() (error, bool) {
	if n.nh == nil {
		return errors.New("node host is nil, has it been started?"), false
	}
	if !n.started {
		return errors.New("node host isn't running"), true
	}
	return nil, true
}
