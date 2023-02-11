/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package routing

import (
	"hash"
	"hash/fnv"
	"sync"
)

const (
	// internal systemic limit for sharding
	shardLimit uint64 = 256
)

func NewShardRouter() *ShardRouter {
	return &ShardRouter{
		hasher:     fnv.New64a(),
		shardCount: 0,
	}
}

type ShardRouter struct {
	hasher     hash.Hash64
	shardCount int64
	mu         sync.Mutex
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
	if s.shardCount == 0 {
		s.shardCount = 1
	}

	var b, j int64
	for j < s.shardCount {
		b = j
		h = h*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((h>>33)+1)))
	}

	return h, nil
}
