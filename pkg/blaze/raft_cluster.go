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

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/rs/zerolog"
)

var (
	_ ICluster = (*ClusterManager)(nil)
)

func newClusterManager(logger zerolog.Logger, nodeHost *dragonboat.NodeHost) *ClusterManager {
	l := logger.With().Str("component", "cluster-manager").Logger()
	return &ClusterManager{l, nodeHost}
}

type ClusterManager struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (c *ClusterManager) StartCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateStateMachineFunc, cfg config.Config) error {
	return c.nh.StartCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StartConcurrentCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateConcurrentStateMachineFunc, cfg config.Config) error {
	return c.nh.StartConcurrentCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StartOnDiskCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateOnDiskStateMachineFunc, cfg config.Config) error {
	return c.nh.StartOnDiskCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StopCluster(clusterID uint64) error {
	return c.nh.StopCluster(clusterID)
}

func (c *ClusterManager) SyncGetClusterMembership(ctx context.Context, clusterID uint64) (*dragonboat.Membership, error) {
	return c.nh.SyncGetClusterMembership(ctx, clusterID)
}

func (c *ClusterManager) ReadIndex(clusterID uint64, timeout time.Duration) (*dragonboat.RequestState, error) {
	return c.nh.ReadIndex(clusterID, timeout)
}

func (c *ClusterManager) ReadLocalNode(rs *dragonboat.RequestState, query interface{}) (interface{}, error) {
	return c.nh.ReadLocalNode(rs, query)
}

func (c *ClusterManager) NAReadLocalNode(rs *dragonboat.RequestState, query []byte) ([]byte, error) {
	return c.nh.NAReadLocalNode(rs, query)
}

func (c *ClusterManager) StaleRead(clusterID uint64, query interface{}) (interface{}, error) {
	return c.nh.StaleRead(clusterID, query)
}
