/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package fsm

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	db "github.com/mxplusb/pleiades/api/v1/database"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3/statemachine"
	"go.etcd.io/bbolt"
)

type op int

const (
	get op = 1
	put op = 2
)

var _ statemachine.IOnDiskStateMachine = &BBoltStateMachine{}

type BBoltStateMachine struct {
	ClusterId uint64
	NodeId    uint64
	BasePath  string
	Options   *bbolt.Options

	db *bbolt.DB
	mu sync.RWMutex
}

func NewBBoltStateMachine(clusterId uint64, nodeId uint64, basePath string, options *bbolt.Options) *BBoltStateMachine {
	return &BBoltStateMachine{ClusterId: clusterId, NodeId: nodeId, BasePath: basePath, Options: options}
}

// dbPath returns the database path with or without appending the database file name.
func (b *BBoltStateMachine) dbPath(withDb bool) string {
	core := filepath.Join(b.BasePath,
		fmt.Sprintf("cluster-%s", strconv.FormatUint(b.ClusterId, 10)),
		fmt.Sprintf("node-%s", strconv.FormatUint(b.NodeId, 10)))
	if !withDb {
		return core
	}
	return filepath.Join(core, "store.db")
}

// Open the bbolt backend and read the last index.
// todo (sienna): leverage stopc at some point on bbolt.Open
func (b *BBoltStateMachine) Open(stopc <-chan struct{}) (uint64, error) {

	dbPath := b.dbPath(true)
	noDb := b.dbPath(false)
	err := os.MkdirAll(noDb, os.FileMode(dbDirModeVal))
	if err != nil {
		return uint64(0), err
	}

	b.db, err = bbolt.Open(dbPath, os.FileMode(dbFileModeVal), b.Options)
	if err != nil {
		return 0, err
	}

	var index uint64

	b.mu.Lock()
	err = b.db.Update(func(tx *bbolt.Tx) error {
		// todo (sienna): implement db stats on open
		//tx.Stats()

		internalBucket, err := tx.CreateBucketIfNotExists([]byte(monotonicLogBucket))
		if err != nil {
			return err
		}

		// todo (sienna): add createIfNotExists to the key.
		key, val := internalBucket.Cursor().Last()
		if key == nil || val == nil {
			index = 0
			return nil
		}

		index = binary.LittleEndian.Uint64(val)
		return nil
	})
	b.mu.Unlock()

	if err != nil {
		return 0, err
	}

	return index, nil
}

func (b *BBoltStateMachine) Update(entries []statemachine.Entry) ([]statemachine.Entry, error) {
	var lastApplied uint64
	applied := make([]statemachine.Entry, 0)

	b.mu.Lock()
	err := b.db.Batch(func(tx *bbolt.Tx) error {
		monotonicBucket, err := tx.CreateBucketIfNotExists([]byte(monotonicLogBucket))
		if err != nil {
			return err
		}

		// prep the last known good applied commit
		lastAppliedVal := monotonicBucket.Get([]byte(monotonicLogKey))
		if lastAppliedVal == nil {
			lastApplied = uint64(0)
		} else {
			lastApplied = binary.LittleEndian.Uint64(lastAppliedVal)
		}

		for idx := range entries {
			kvp := db.KeyValue{}
			if err := kvp.UnmarshalVT(entries[idx].Cmd); err != nil {
				return err
			}

			key := kvp.Key
			if len(key) == 0 {
			}

			parentBucketName, childBucketNames, err := prepBucket(kvp)
			if err != nil {
				return err
			}

			parentBucket, err := tx.CreateBucketIfNotExists([]byte(parentBucketName))
			if err != nil {
				return err
			}

			retVal := make(chan []byte)
			if err := keyOp(parentBucket, childBucketNames, entries[idx].Cmd, put, &retVal); err != nil {
				return err
			}

			// store the current index as the last applied commit
			lastAppliedPayload := make([]byte, 8)
			binary.LittleEndian.PutUint64(lastAppliedPayload, entries[idx].Index)
			if err := monotonicBucket.Put([]byte(monotonicLogKey), lastAppliedPayload); err != nil {
				return err
			}

			entries[idx].Result = statemachine.Result{Value: uint64(len(entries[idx].Cmd))}
			applied = append(applied, entries[idx])
		}

		return nil
	})
	b.mu.Unlock()

	if err != nil {
		return make([]statemachine.Entry, 0), err
	}

	err = b.db.View(func(tx *bbolt.Tx) error {
		monotonicBucket := tx.Bucket([]byte(monotonicLogBucket))
		if monotonicBucket == nil {
			return fmt.Errorf("the %s bucket does not exist, no writes applied", monotonicLogBucket)
		}

		// prep the last known good applied commit
		lastAppliedVal := monotonicBucket.Get([]byte(monotonicLogKey))
		currentIndex := binary.LittleEndian.Uint64(lastAppliedVal)
		if currentIndex == lastApplied {
			return errors.New("none of the commits were applied")
		}

		if currentIndex != entries[len(entries)-1].Index {
			return errors.New("not all entries were applied")
		}
		return nil
	})
	if err != nil {
		return applied, err
	}

	return entries, err
}

// prepBucket verifies the key signature. the string is the root bucket, the string slice is the rest of the bucket hierarchy, and the error is any parsing errs
func prepBucket(kvp db.KeyValue) (string, []string, error) {
	// verify we're not trying to create an empty bucket and skip the first item
	key := string(kvp.Key)

	if len(key) == 0 {
		return "", []string{}, fmt.Errorf("key is empty")
	}

	bucketHierarchy := strings.Split(string(key[:]), "/")[1:]
	bucketHierarchyLen := len(bucketHierarchy)
	if bucketHierarchy[bucketHierarchyLen-1] == "" {
		return "", []string{}, errors.New("cannot create empty bucket")
	}

	if bucketHierarchyLen < fsmRootKeyCount {
		return "", []string{}, errors.New("the fsm root key count is not correct")
	}

	if bucketHierarchyLen+3 > maxKeyDepth {
		return "", []string{}, fmt.Errorf("the nested key cannot be more than %d levels deep", maxKeyDepth)
	}

	return bucketHierarchy[0], bucketHierarchy[1:], nil
}

// keyOp recursively creates buckets until it can put or get the key
func keyOp(parentBucket *bbolt.Bucket, bucketHierarchy []string, val []byte, operation op, retVal *chan []byte) error {
	// the last value in the bucketHierarchy is 1, it's the key, which makes the parent bucket the desired bucket
	if len(bucketHierarchy) == 1 {
		switch operation {
		case get:
			targetVal := parentBucket.Get([]byte(bucketHierarchy[0]))
			*retVal <- targetVal
			//sent = true
			return nil
		case put:
			return parentBucket.Put([]byte(bucketHierarchy[0]), val)
		}
		return nil
	}

	// there's more hierarchy, keep expanding.
	if len(bucketHierarchy) >= 2 {
		childBucket, err := parentBucket.CreateBucketIfNotExists([]byte(bucketHierarchy[0]))
		if err != nil {
			return err
		}
		return keyOp(childBucket, bucketHierarchy[1:], val, operation, retVal)
	}

	return nil
}

func (b *BBoltStateMachine) Lookup(i interface{}) (interface{}, error) {
	payload := db.KeyValue{}

	b.mu.Lock()
	err := b.db.Update(func(tx *bbolt.Tx) error {
		kvp := db.KeyValue{}
		if err := kvp.UnmarshalVT(i.([]byte)); err != nil {
			return err
		}

		key := kvp.Key

		if string(key) == "" {
			return errors.New("cannot find an empty key")
		}

		parentBucketName, childBucketNames, err := prepBucket(kvp)
		if err != nil {
			return err
		}

		bucket := tx.Bucket([]byte(parentBucketName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("the root bucket %s does not exist", parentBucketName))
		}

		retVal := make(chan []byte, 1)
		if err := keyOp(bucket, childBucketNames, make([]byte, 0), get, &retVal); err != nil {
			return err
		}

		// if it's empty, that's okay, we just don't want a serialized value
		target := <-retVal
		if target == nil {
			return nil
		}

		if err := payload.UnmarshalVT(target); err != nil {
			return err
		}

		return nil
	})
	b.mu.Unlock()

	// todo (sienna): figure this out, it will vail `go vet`
	//goland:noinspection ALL
	return payload, err
}

func (b *BBoltStateMachine) Sync() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.db.Sync()
}

func (b *BBoltStateMachine) PrepareSnapshot() (interface{}, error) {
	return nil, nil
}

func (b *BBoltStateMachine) SaveSnapshot(ctx interface{}, writer io.Writer, done <-chan struct{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.WriteTo(writer)
		return err
	})
}

func (b *BBoltStateMachine) RecoverFromSnapshot(reader io.Reader, i <-chan struct{}) error {
	fn := func(r io.Reader) error {
		target, err := os.Create(b.dbPath(true))
		if err != nil {
			return err
		}
		_, err = io.Copy(target, reader)
		if err != nil {
			return err
		}
		return nil
	}

	// verify the existing database is closed
	b.mu.Lock()
	err := b.db.Close()
	if err != nil {
		b.mu.Unlock()
		return err
	}
	b.mu.Unlock()

	b.mu.Lock()
	_, err = os.Stat(b.dbPath(true))
	if err != nil {
		if os.IsNotExist(err) {
			b.mu.Unlock()
			return fn(reader)
		}
		b.mu.Unlock()
		return err
	}
	b.mu.Unlock()

	b.mu.Lock()
	err = os.Remove(b.dbPath(true))
	if err != nil {
		b.mu.Unlock()
		return err
	}
	b.mu.Unlock()

	return fn(reader)
}

func (b *BBoltStateMachine) Close() error {
	b.mu.Lock()

	err := b.db.Close()
	if err != nil {
		b.mu.Unlock()
		return err
	}

	b.db = nil
	b.mu.Unlock()

	return nil
}
