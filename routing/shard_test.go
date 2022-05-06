package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShardAssignment(t *testing.T) {
	shard := GetShardAssignment("abcd1234:bucket:key")
	assert.LessOrEqual(t, shard, int32(shardLimit))
}

func FuzzGetShardAssignment(f *testing.F) {
	f.Add("abcd1234")
	f.Fuzz(func(t *testing.T, s string) {
		shard := GetShardAssignment(s)
		assert.LessOrEqual(t, shard, int32(shardLimit))
	})
}

func BenchmarkGetShardAssignment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetShardAssignment("abcd1234")
	}
}
