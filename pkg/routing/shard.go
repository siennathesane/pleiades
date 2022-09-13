
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
	"strconv"

	"github.com/lithammer/go-jump-consistent-hash"
)

const (
	// internal systemic limit for sharding
	shardLimit uint64 = 256
)

type Shard struct {
	j *jump.Hasher
}

func (s *Shard) AccountToShard(accountId uint64) uint64 {
	if s.j == nil {
		s.j = jump.New(int(shardLimit), &farmHash{})
	}
	return uint64(s.j.Hash(strconv.FormatUint(accountId, 10)))
}
