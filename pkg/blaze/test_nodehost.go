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
	"github.com/lni/dragonboat/v3/config"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/lni/goutils/vfs"
)

func buildTestNodeHost(t *testing.T) *dragonboat.NodeHost {
	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)

	cfg := config.Config{
		NodeID:    1,
		ClusterID: 1,
		HeartbeatRTT: 10,
		ElectionRTT: 100,
	}

	expertCfg := config.GetDefaultExpertConfig()
	expertCfg.LogDB.Shards = 4
	expertCfg.FS = vfs.NewMem()

	nhConf := config.NodeHostConfig{
		WALDir:         t.TempDir(),
		NodeHostDir:    t.TempDir(),
		RTTMillisecond: 10,
		RaftAddress:    fmt.Sprintf("localhost:%d", port),
		Expert:         expertCfg,
	}

	host, err := dragonboat.NewNodeHost(nhConf)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if err = host.StartCluster(nil, true, newTestStateMachine, cfg); err != nil {
		t.Fatalf(err.Error())
	}

	return host
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
