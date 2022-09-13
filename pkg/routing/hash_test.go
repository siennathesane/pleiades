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
	"testing"

	"github.com/lithammer/go-jump-consistent-hash"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	hasher := jump.New(int(shardLimit), &farmHash{})

	target := strconv.FormatUint(10, 10)
	val := hasher.Hash(target)
	assert.Equal(t, 63, val, "the value must be 63")
}