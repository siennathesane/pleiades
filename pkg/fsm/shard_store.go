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
	"os"
	"path/filepath"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/fsm/kv"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"go.etcd.io/bbolt"
)

func NewShardStore(logger zerolog.Logger) (*ShardStore, error) {
	basePath := configuration.Get().GetString("server.datastore.basePath")
	dbPath := filepath.Join(basePath, "shard-config.db")

	db, err := bbolt.Open(dbPath, os.FileMode(dbDirModeVal), nil)
	if err != nil {
		return nil, err
	}

	return &ShardStore{
		logger: logger.With().Str("component", "shard-config").Logger(),
		db:     db,
	}, nil
}

type ShardStore struct {
	logger zerolog.Logger
	db *bbolt.DB
}

func (s *ShardStore) Get(shardId uint64) (*raftv1.NewShardRequest, error) {
	req := &raftv1.NewShardRequest{}
	err :=  s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(shardConfigBucket))
		if bucket == nil {
			return errors.New("no shards configured")
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf,shardId)

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

func (s *ShardStore) Put(req *raftv1.NewShardRequest) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(shardConfigBucket))
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't open shard config bucket")
			return err
		}

		payload, err := req.MarshalVT()
		if err != nil {
			s.logger.Trace().Err(err).Msg("can't unmarshal request")
			return err
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf,req.GetShardId())

		return bucket.Put(shardIdBuf, payload)
	})
}

func (s *ShardStore) Delete(shardId uint64) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(shardConfigBucket))
		if bucket == nil {
			return errors.New("no shards configured")
		}

		shardIdBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardIdBuf,shardId)

		return bucket.Delete(shardIdBuf)
	})
}
