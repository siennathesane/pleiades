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
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/lni/goutils/vfs"
)

func buildTestNodeHostConfig(t *testing.T) dconfig.NodeHostConfig {
	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)

	expertCfg := dconfig.GetDefaultExpertConfig()
	expertCfg.LogDB.Shards = 4
	expertCfg.FS = vfs.NewMem()

	return dconfig.NodeHostConfig{
		WALDir:         t.TempDir(),
		NodeHostDir:    t.TempDir(),
		RTTMillisecond: 10,
		RaftAddress:    fmt.Sprintf("localhost:%d", port),
		Expert:         expertCfg,
		NotifyCommit:   false,
	}
}

func buildTestClusterConfig(t *testing.T) dconfig.Config {
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

	if err = host.StartCluster(nil, true, newTestStateMachine, buildTestClusterConfig(t)); err != nil {
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
		clusterConfig := buildTestClusterConfig(t)
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
		if err := nodeHosts[i].StartCluster(initialMembers, false, newTestStateMachine, buildTestClusterConfig(t)); err != nil {
			t.Fatalf(err.Error())
		}
	}

	return clusterConfigs, nodeConfigs, nodeHosts
}

type testStateMachine struct {
	ShardID   uint64
	ReplicaID uint64
	Count     uint64
}

func newTestStateMachine(shardID uint64,
	replicaID uint64) sm.IStateMachine {
	return &testStateMachine{
		ShardID:   shardID,
		ReplicaID: replicaID,
		Count:     0,
	}
}

func (s *testStateMachine) Lookup(query interface{}) (interface{}, error) {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result, nil
}

func (s *testStateMachine) Update(data []byte) (sm.Result, error) {
	s.Count++
	fmt.Printf("from testStateMachine.Update(), msg: %s, count:%d\n",
		string(data), s.Count)
	return sm.Result{Value: uint64(len(data))}, nil
}

func (s *testStateMachine) SaveSnapshot(w io.Writer, fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	return err
}

func (s *testStateMachine) RecoverFromSnapshot(r io.Reader,
	files []sm.SnapshotFile,
	done <-chan struct{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v
	return nil
}

func (s *testStateMachine) Close() error { return nil }
