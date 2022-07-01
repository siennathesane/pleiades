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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
)

var (
	validateHostConfig = true
)

type Environment string

const (
	Development            Environment = "development"
	Production             Environment = "production"
	HostConfigPathTemplate string      = "hosts/%s/config/env"
)

func ProvideEnvironmentConfig() fx.Option {
	return fx.Provide(NewEnvironmentConfig)
}

type EnvironmentConfig struct {
	Environment     Environment `json:"environment"`
	GCPProjectId    string      `json:"gcp-project-id"`
	BaseDir         string      `json:"base-dir"`
	LocalClusterId  uint64      `json:"local-cluster-id"`
	LocalExchangeId uint64      `json:"local-exchange-id"`
	BasePort        int         `json:"base-port"`
	MaxPort         int         `json:"max-port"`
	Hostname        string      `json:"hostname,omitempty"`
}

func NewEnvironmentConfig(client *api.Client) (*EnvironmentConfig, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	hostKey := fmt.Sprintf(HostConfigPathTemplate, hostname)
	pair, _, err := client.KV().Get(hostKey, nil)
	if err != nil {
		return nil, err
	}

	var e EnvironmentConfig
	if err := json.Unmarshal(pair.Value, &e); err != nil {
		return nil, err
	}

	e.Hostname = hostname
	if validateHostConfig {
		if err := e.validate(); err != nil {
			return nil, err
		}
	}

	return &e, nil
}

func (e EnvironmentConfig) validate() error {
	if e.BasePort <= 1024 {
		return errors.New("base_port needs to be more than 1024")
	}
	if (e.MaxPort - e.BasePort) < 4000 {
		return errors.New("max_port needs to be at least 4000 higher than base_port")
	}
	return nil
}
