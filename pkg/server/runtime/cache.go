/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

import (
	"github.com/allegro/bigcache"
	"github.com/rs/zerolog"
)

var (
	cacheSingleton *Cache
)

func NewCache(logger zerolog.Logger) (*Cache, error) {
	if cacheSingleton != nil {
		return cacheSingleton, nil
	}

	var err error
	cacheSingleton.BigCache, err = bigcache.NewBigCache(bigcache.DefaultConfig(0))
	if err != nil {
		return nil, err
	}

	return cacheSingleton, nil
}

// Cache is a global singleton used to store and retrieve volatile information at runtime.
type Cache struct {
	*bigcache.BigCache
}
