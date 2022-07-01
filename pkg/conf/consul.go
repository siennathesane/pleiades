/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package conf

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
)

// ProvideConsulClient is the main constructor for providing the Consul agent client
func ProvideConsulClient() fx.Option {
	return fx.Provide(NewConsulClient)
}

// NewConsulClient creates a new client for Consul
func NewConsulClient(lifecycle fx.Lifecycle) (*api.Client, error) {
	return api.NewClient(api.DefaultConfig())
}
