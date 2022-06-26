/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"fmt"

	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
)

type Registry struct {
	logger zerolog.Logger
	cache map[string]*server.Server
}

func NewRegistry(logger zerolog.Logger) (*Registry, error) {
	l := logger.With().Str("component", "registry").Logger()
	return &Registry{logger: l, cache: make(map[string]*server.Server)}, nil
}

func (r *Registry) Get(key string) (*server.Server, error) {
	val, ok := r.cache[key]
	if !ok {
		return nil, fmt.Errorf("no server found for key: %s", key)
	}
	return val, nil
}

func (r *Registry) Put(key string, srv *server.Server) error {
	r.cache[key] = srv
	return nil
}
