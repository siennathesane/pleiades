/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

// todo (sienna): there is a bug here. if the number of shards this host sees

package kvstore

import (
	"hash"
	"hash/fnv"
	"sync/atomic"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
)

const (
	// internal systemic limit for sharding
	shardLimit uint64 = 256
)

func NewShardRouter(host runtime.IHost, logger zerolog.Logger) *ShardRouter {
	return &ShardRouter{
		logger:     logger,
		nh:         host,
		hasher:     fnv.New64a(),
		shardCount: atomic.Int64{},
	}
}

type ShardRouter struct {
	nh         runtime.IHost
	logger     zerolog.Logger
	hasher     hash.Hash64
	shardCount atomic.Int64
}

func (s *ShardRouter) AccountToShard(accountId uint64) uint64 {
	return 0
}

// CalcShard is an implementation of Jump Consistent Hash. It computes a string key and returns the shard the key can be
// found in.
// ref: https://arxiv.org/ftp/arxiv/papers/1406/1406.2294.pdf
func (s *ShardRouter) CalcShard(key []byte) (uint64, error) {

	s.hasher.Reset()
	_, err := s.hasher.Write(key)
	if err != nil {
		return 0, err
	}
	h := s.hasher.Sum64()

	// todo (sienna): this is a pending bug
	if s.shardCount.Load() == 0 {
		s.shardCount.Store(1)
	}

	var b, j int64
	for j < s.shardCount.Load() {
		b = j
		h = h*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((h>>33)+1)))
	}

	return h, nil
}

func (s *ShardRouter) HandleShardUpdate(event *raftv1.RaftEvent) {
	s.logger.Debug().Interface("payload", event).Msg("leader update received")

	// safety check
	if event.Typ != raftv1.EventType_EVENT_TYPE_RAFT {
		s.logger.Error().Msg("event type mismatched")
		return
	}

	hi := s.nh.GetHostInfo(runtime.HostInfoOption{SkipLogInfo: true})
	if hi == nil {
		s.logger.Error().Msg("runtime host information is blank, cannot update shard router count")
		return
	}

	var count int64
	for idx := range hi.ClusterInfoList {
		if !hi.ClusterInfoList[idx].IsWitness || !hi.ClusterInfoList[idx].IsObserver {
			count++
		}
	}

	if ok := s.shardCount.CompareAndSwap(s.shardCount.Load(), count); !ok {
		s.logger.Error().Bool("swapped", ok).Msg("shard count not swapped")
	}

	return
}
