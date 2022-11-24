/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package systemstore

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sort"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/fsm"
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"go.etcd.io/bbolt"
)

const (
	ShardConfigBucket   string = "shards"
	dbDirModeVal int = 484
)

var (
	ErrNoShards = errors.New("no shards configured")
	ErrNoKeys = errors.New("no keys found")

	systemStoreSingleton *SystemStore
)

func NewSystemStore(logger zerolog.Logger) (*SystemStore, error) {
	if systemStoreSingleton != nil {
		return systemStoreSingleton, nil
	}

	basePath := configuration.Get().GetString("server.datastore.basePath")
	dbPath := filepath.Join(basePath, "system.db")

	db, err := bbolt.Open(dbPath, os.FileMode(dbDirModeVal), nil)
	if err != nil {
		return nil, err
	}

	return &SystemStore{
		logger: logger.With().Str("component", "system-config").Logger(),
		db:     db,
	}, nil
}

type SystemStore struct {
	logger zerolog.Logger
	db     *bbolt.DB
}

func (s *SystemStore) Close() error {
	s.logger.Debug().Msg("shutting down shard storage")
	return s.db.Close()
}

func (s *SystemStore) GetAllShards() ([]*raftv1.ShardState, error) {
	reqs := make([]*raftv1.ShardState, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(fsm.ShardConfigBucket))
		if bucket == nil {
			return ErrNoShards
		}

		return bucket.ForEach(func(k, v []byte) error {
			req := &raftv1.ShardState{}
			err := req.UnmarshalVT(v)
			if err != nil {
				s.logger.Trace().Err(err).Msg("can't unmarshal shard configuration")
			}

			s.logger.Trace().Interface("shard-state", req).Msg("found shard configuration")
			reqs = append(reqs, req)

			return nil
		})
	})

	sort.SliceStable(reqs, func(i, j int) bool {
		return reqs[i].GetShardId() < reqs[j].GetShardId()
	})

	return reqs, err
}

// Deprecated: use Get
func (s *SystemStore) GetShard(shardId uint64) (*raftv1.ShardState, error) {
	req := &raftv1.ShardState{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(fsm.ShardConfigBucket))
		if bucket == nil {
			return ErrNoShards
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf, shardId)

		payload := bucket.Get(shardIdBuf)
		if payload == nil {
			return kv.ErrKeyNotFound
		}

		err := req.UnmarshalVT(payload)
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't unmarshal shard configuration")
		}
		return err
	})
	return req, err
}

func (s *SystemStore) GetAll(bucket string) (map[string][]byte, error) {
	reqs := make(map[string][]byte, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return errors.New("bucket doesn't exist")
		}

		return bucket.ForEach(func(k, v []byte) error {
			reqs[string(k)] = v
			return nil
		})
	})

	return reqs, err
}

func (s *SystemStore) Get(bucket, key string) ([]byte, error) {
	var resp []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return ErrNoKeys
		}

		payload := bucket.Get([]byte(key))
		if payload == nil {
			return kv.ErrKeyNotFound
		}

		resp = make([]byte, len(payload))
		copy(resp, payload)
		return nil
	})
	return resp, err
}

func (s *SystemStore) Put(bucket, key string, payload []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		if bucket == "" {
			return errors.New("empty bucket")
		}
		if key == "" {
			return errors.New("empty key")
		}
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't open bucket")
			return err
		}

		return bucket.Put([]byte(key), payload)
	})
}

func (s *SystemStore) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return errors.New("bucket doesn't exist")
		}

		return bucket.Delete([]byte(key))
	})
}

func (s *SystemStore) DeleteBucket(bucket string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// Deprecated: use Put
func (s *SystemStore) PutShard(req *raftv1.ShardState) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(fsm.ShardConfigBucket))
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't open shard config bucket")
			return err
		}

		s.logger.Trace().Interface("request", req).Msg("storing shard state")

		payload, err := req.MarshalVT()
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't unmarshal request")
			return err
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf, req.GetShardId())

		return bucket.Put(shardIdBuf, payload)
	})
}

// Deprecated: use Delete
func (s *SystemStore) DeleteShard(shardId uint64) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(fsm.ShardConfigBucket))
		if bucket == nil {
			return ErrNoShards
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf, shardId)

		return bucket.Delete(shardIdBuf)
	})
}
