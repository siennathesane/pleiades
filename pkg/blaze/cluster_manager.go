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
	nodeHost *dragonboat.NodeHost
}

func (c *ClusterManager) StartCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateStateMachineFunc, cfg config.Config) error {
	return c.nodeHost.StartCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StartConcurrentCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateConcurrentStateMachineFunc, cfg config.Config) error {
	return c.nodeHost.StartConcurrentCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StartOnDiskCluster(initialMembers map[uint64]dragonboat.Target, join bool, create statemachine.CreateOnDiskStateMachineFunc, cfg config.Config) error {
	return c.nodeHost.StartOnDiskCluster(initialMembers, join, create, cfg)
}

func (c *ClusterManager) StopCluster(clusterID uint64) error {
	return c.nodeHost.StopCluster(clusterID)
}

func (c *ClusterManager) SyncGetClusterMembership(ctx context.Context, clusterID uint64) (*dragonboat.Membership, error) {
	return c.nodeHost.SyncGetClusterMembership(ctx, clusterID)
}

