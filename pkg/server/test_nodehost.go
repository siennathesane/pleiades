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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/lni/goutils/vfs"
)

func buildTestNodeHostConfig(t *testing.T) dconfig.NodeHostConfig {
	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)

	// this is needed or you'll hit the ulimit for the temp dirs.
	expertConf := dconfig.ExpertConfig{
		FS: vfs.NewMem(),
	}

	return dconfig.NodeHostConfig{
		WALDir:         t.TempDir(),
		NodeHostDir:    t.TempDir(),
		RTTMillisecond: 1,
		RaftAddress:    fmt.Sprintf("localhost:%d", port),
		NotifyCommit:   false,
		Expert: expertConf,
	}
}

func buildTestShardConfig(t *testing.T) dconfig.Config {
	rand.Seed(time.Now().UTC().UnixNano())
	nodeId := rand.Intn(10_000)
	clusterId := rand.Intn(10_000)

	return dconfig.Config{
		NodeID:       uint64(nodeId),
		ClusterID:    uint64(clusterId),
		HeartbeatRTT: 10,
		ElectionRTT:  100,
	}
}

func buildTestNodeHost(t *testing.T) *dragonboat.NodeHost {
	host, err := dragonboat.NewNodeHost(buildTestNodeHostConfig(t))
	if err != nil {
		t.Fatalf(err.Error())
	}

	return host
}

func build3NodeTestCluster(t *testing.T) ([]dconfig.Config, []dconfig.NodeHostConfig, []*dragonboat.NodeHost) {
	clusterConfigs := make([]dconfig.Config, 3)
	nodeConfigs := make([]dconfig.NodeHostConfig, 3)
	nodeHosts := make([]*dragonboat.NodeHost, 3)

	initialMembers := make(map[uint64]dragonboat.Target)
	for i := 0; i < 3; i++ {
		clusterConfig := buildTestShardConfig(t)
		nodeConfig := buildTestNodeHostConfig(t)
		clusterConfigs[i] = clusterConfig
		nodeConfigs[i] = nodeConfig

		host, err := dragonboat.NewNodeHost(nodeConfig)
		if err != nil {
			t.Fatalf(err.Error())
		}
		nodeHosts[i] = host
		initialMembers[clusterConfig.NodeID] = host.RaftAddress()
	}

	for i := 0; i < len(nodeHosts); i++ {
		if err := nodeHosts[i].StartCluster(initialMembers, false, newTestStateMachine, buildTestShardConfig(t)); err != nil {
			t.Fatalf(err.Error())
		}
	}

	return clusterConfigs, nodeConfigs, nodeHosts
}
