/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog"
)

type Registry struct {
	logger zerolog.Logger
	cache *bigcache.BigCache
}

func NewRegistry(logger zerolog.Logger) (*Registry, error) {
	l := logger.With().Str("component", "registry").Logger()
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(1024 * 1024 * 1024))
	if err != nil {
		return nil, err
	}
	return &Registry{logger: l, cache: cache}, nil
}

func (r *Registry) Get(key string) (interface{}, error) {
	return r.cache.Get(key)
}

func (r *Registry) Put(key string, payload []byte) error {
	return r.cache.Set(key, payload)
}