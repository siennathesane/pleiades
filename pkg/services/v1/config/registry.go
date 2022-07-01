/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package config

import (
	"fmt"

	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/protocols/v1/config"
)

type Registry struct {
	logger           zerolog.Logger
	serverMap        map[config.ServiceType_Type]any
	clientFactoryMap map[string]any
}

func NewRegistry(logger zerolog.Logger) (*Registry, error) {
	l := logger.With().Str("component", "registry").Logger()
	return &Registry{logger: l, serverMap: make(map[config.ServiceType_Type]any)}, nil
}

func (r *Registry) GetServer(key config.ServiceType_Type) (any, error) {
	val, ok := r.serverMap[key]
	if !ok {
		return nil, fmt.Errorf("no server found for key: %s", key)
	}
	return val, nil
}

func (r *Registry) GetClientFactory(key string) (any, error) {
	val, ok := r.clientFactoryMap[key]
	if !ok {
		return nil, fmt.Errorf("no client factory found for key: %s", key)
	}
	return val, nil
}

func (r *Registry) PutServer(key config.ServiceType_Type, srv any) error {
	r.serverMap[key] = srv
	return nil
}

func (r *Registry) PutClientFactory(key string, f any) error {
	r.clientFactoryMap[key] = f
	return nil
}
