package routing

import (
	"github.com/cespare/xxhash/v2"
)

const (
	// internal systemic limit for sharding
	shardLimit int = 1024
)

//GetShardAssignment determines which shard a given key is assigned to
func GetShardAssignment(s string) int32 {
	return Hash(xxhash.Sum64String(s), shardLimit)
}
