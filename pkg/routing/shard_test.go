/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package routing

import (
	crand "crypto/rand"
	"math/rand"
	"testing"
)

const (
	fuzzLimit = 256
)

//func FuzzGetShardAssignment(f *testing.F) {
//	for i := 0; i < fuzzLimit; i++ {
//		buf := make([]byte, rand.Uint32())
//		_, err := crand.Read(buf)
//		if err != nil {
//			f.Fatal(err)
//		}
//		f.Add(i)
//	}
//	f.Fuzz(func(t *testing.T, k []byte) {
//		shardRouter := NewShardRouter()
//		shardRouter.shardCount = 128
//		shard, err := shardRouter.CalcShard(k)
//		assert.LessOrEqual(t, shard, shardLimit, "the shard must be within the shard range")
//		assert.NoError(t, err, "there must not be an error trying to calculate the shard assignment")
//	})
//}

func BenchmarkGetShardAssignment(b *testing.B) {
	shardRouter := NewShardRouter()
	shardRouter.shardCount = 128
	for i := 0; i < b.N; i++ {

		// this is <10ns, it will hang the benchmark if the timer is excluded
		buf := make([]byte, rand.Intn(128))
		_, err := crand.Read(buf)
		if err != nil {
			b.Fatal(err)
		}

		_, err = shardRouter.CalcShard(buf)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}
