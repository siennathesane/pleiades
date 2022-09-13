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
	"testing"

	"github.com/lithammer/go-jump-consistent-hash"
	"github.com/stretchr/testify/assert"
)

const (
	fuzzLimit = 256
)

func FuzzGetShardAssignment(f *testing.F) {
	for i := uint64(0); i < fuzzLimit; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, a uint64) {
		shardRouter := &Shard{}
		shard := shardRouter.AccountToShard(a)
		assert.LessOrEqual(t, shard, shardLimit, "the shard must be within the shard range")
	})
}

func BenchmarkGetShardAssignment(b *testing.B) {
	shardRouter := &Shard{
		j: jump.New(int(shardLimit), &farmHash{}),
	}
	for i := 0; i < b.N; i++ {
		shardRouter.AccountToShard(uint64(i))
	}
}
