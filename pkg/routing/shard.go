
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
