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
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/rs/zerolog"
)

// TestStateMachine is the IStateMachine implementation used in the
// helloworld example.
// See https://github.com/lni/dragonboat/blob/master/statemachine/rsm.go for
// more details of the IStateMachine interface.
type TestStateMachine struct {
	ClusterID uint64
	NodeID    uint64
	Count     uint64
	logger  zerolog.Logger
}

// NewExampleStateMachine creates and return a new TestStateMachine object.
func NewTestStateMachine(clusterID uint64,
	nodeID uint64) sm.IStateMachine {
	return &TestStateMachine{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
		logger: zerolog.New(os.Stdout).With().Uint64("shard", clusterID).Uint64("replica", nodeID).Logger(),
	}
}

// Lookup performs local lookup on the testStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *TestStateMachine) Lookup(query interface{}) (interface{}, error) {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result, nil
}

// Update updates the object using the specified committed raft entry.
func (s *TestStateMachine) Update(data []byte) (sm.Result, error) {
	// in this example, we print out the following hello world message for each
	// incoming update request. we also increase the counter by one to remember
	// how many updates we have applied
	s.Count++
	s.logger.Debug().Uint64("count", s.Count).Str("payload", string(data)).Msg("testStateMachine.Update() called")
	return sm.Result{Value: uint64(len(data))}, nil
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *TestStateMachine) SaveSnapshot(w io.Writer,
	fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	return err
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *TestStateMachine) RecoverFromSnapshot(r io.Reader,
	files []sm.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable, that is the only state we maintain in this
	// example, the input files is expected to be empty
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v
	return nil
}

// Close closes the IStateMachine instance. There is nothing for us to cleanup
// or release as this is a pure in memory data store. Note that the Close
// method is not guaranteed to be called as node can crash at any time.
func (s *TestStateMachine) Close() error { return nil }
