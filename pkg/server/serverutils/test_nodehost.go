/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package serverutils

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/lni/goutils/vfs"
)

var (
	called = false
)

func WaitForReadyCluster(t *testing.T, shardId uint64, host *dragonboat.NodeHost, timeout time.Duration) {
	now := time.Now()
	expiry := time.Now().Add(timeout)
	for {
		if now.UnixMilli() > expiry.UnixMilli() {
			t.Error("timeout: can't get shard session")
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		_, err := host.SyncGetSession(ctx, shardId)
		if err != nil {
			if err == dragonboat.ErrClusterNotReady {
				cancel()
				continue
			}
			cancel()
			continue
		}
		cancel()
		return
	}
}

func BuildTestNodeHostConfig(t *testing.T) dconfig.NodeHostConfig {
	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)

	// this is needed, or you'll hit the ulimit for the temp dirs.
	expertConf := dconfig.ExpertConfig{
		FS: vfs.NewMem(),
	}

	msg, err := messaging.NewEmbeddedMessagingWithDefaults(utils.NewTestLogger(t))
	if err != nil {
		t.Fatalf("error starting nats: %s", err)
	}
	msg.Start()

	pubSubClient, err := msg.GetPubSubClient()
	if err != nil {
		t.Fatalf("error getting pubsub client: %s", err)
	}

	queueClient, err := msg.GetStreamClient()
	if err != nil {
		t.Fatalf("error getting queue client: %s", err)
	}

	sysListener, err := messaging.NewRaftSystemListener(pubSubClient, queueClient, utils.NewTestLogger(t))

	logger := utils.NewTestLogger(t)
	SetRootLogger(logger)
	if !called {
		dlog.SetLoggerFactory(DragonboatLoggerFactory)
		called = true
	}

	return dconfig.NodeHostConfig{
		WALDir:              t.TempDir(),
		NodeHostDir:         t.TempDir(),
		RTTMillisecond:      1,
		RaftAddress:         fmt.Sprintf("localhost:%d", port),
		NotifyCommit:        false,
		Expert:              expertConf,
		SystemEventListener: sysListener,
		RaftEventListener:   sysListener,
	}
}

func BuildTestShardConfig(t *testing.T) dconfig.Config {
	rand.Seed(time.Now().UTC().UnixNano())
	nodeId := rand.Intn(10_000)
	clusterId := rand.Intn(10_000)

	return dconfig.Config{
		NodeID:       uint64(nodeId),
		ClusterID:    uint64(clusterId),
		HeartbeatRTT: 1,
		ElectionRTT:  10,
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
	utils.Wait(100 * time.Millisecond)

	return firstTestHost, firstNodeClusterConfig
}
