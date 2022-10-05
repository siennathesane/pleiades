/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/lni/goutils/vfs"
)

func BuildTestNodeHostConfig(t *testing.T) dconfig.NodeHostConfig {
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

func BuildTestShardConfig(t *testing.T) dconfig.Config {
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

func BuildTestNodeHost(t *testing.T) *dragonboat.NodeHost {
	host, err := dragonboat.NewNodeHost(BuildTestNodeHostConfig(t))
	if err != nil {
		t.Fatalf(err.Error())
	}

	return host
}

func BuildTestShard(t *testing.T) (*dragonboat.NodeHost, dconfig.Config) {
	firstTestHost := BuildTestNodeHost(t)
	firstNodeClusterConfig := BuildTestShardConfig(t)
	nodeClusters := make(map[uint64]string)
	nodeClusters[firstNodeClusterConfig.NodeID] = firstTestHost.RaftAddress()

	err := firstTestHost.StartCluster(nodeClusters, false, NewTestStateMachine, firstNodeClusterConfig)
	if err != nil {
		t.Fatalf("cannot start test state machine: %s", err)
	}
	Wait(100*time.Millisecond)

	return firstTestHost, firstNodeClusterConfig
}
