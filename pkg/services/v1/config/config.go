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
	"context"
	"errors"

	"github.com/mxplusb/pleiades/pkg/fsm"
	configv1 "github.com/mxplusb/pleiades/pkg/protocols/v1/config"
	"github.com/mxplusb/pleiades/pkg/services"
	"capnproto.org/go/capnp/v3"
	"github.com/rs/zerolog"
)

var (
	_ configv1.ConfigService_Server = (*ConfigServer)(nil)
)

type ConfigServer struct {
	logger      zerolog.Logger
	raftManager *fsm.ConfigServiceStoreManager
}

// NewConfigServer creates a instance of the configuration service. This is a singleton.
// The configuration service is responsible for managing all the service available on a deployed host.
func NewConfigServer(store *services.StoreManager, logger zerolog.Logger) (*ConfigServer, error) {
	l := logger.With().Str("service", "config").Logger()
	manager, err := fsm.NewConfigServiceStoreManager(logger, store)
	if err != nil {
		return nil, err
	}

	return &ConfigServer{
		logger:      l,
		raftManager: manager,
	}, nil
}

func (c *ConfigServer) GetConfig(ctx context.Context, call configv1.ConfigService_getConfig) error {
	req, err := call.Args().Request()
	if err != nil {
		return err
	}

	what := req.What()
	amount := req.Amount()
	switch what {
	case configv1.GetConfigurationRequest_Type_all:
	case configv1.GetConfigurationRequest_Type_raft:
		switch amount {
		case configv1.GetConfigurationRequest_Specificity_one:
			key, err := req.Id()
			if err != nil {
				c.logger.Err(err).Str("key", key).Msg("error reading key")
				return err
			}
			return c.getRaftConfig(ctx, key, call)
		case configv1.GetConfigurationRequest_Specificity_everything:
			return c.getAllRaftConfigs(ctx, call)
		}
	}

	return nil
}

func (c *ConfigServer) getRaftConfig(ctx context.Context, key string, call configv1.ConfigService_getConfig) error {
	if key == "" {
		return errors.New("cannot request a named record without a key")
	}

	// allocate the results
	res, err := call.AllocResults()
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error allocating results")
		return err
	}

	// acknowledge the request now that we can allocate space for the results
	call.Ack()

	// find the target
	val, err := c.raftManager.Get(key)
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error getting config")
		return err
	}

	// prep the results
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error creating message")
		return err
	}

	// create the root config response
	config, err := configv1.NewRootGetConfigurationResponse(seg)
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error creating getRaftConfigurationResponse")
		return err
	}

	// create the config list
	slice, err := configv1.NewRaftConfiguration_List(seg, 1)
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error creating raftConfiguration_List")
		return err
	}

	// set the config into the slice
	err = slice.Set(0, *val)
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error setting raftConfiguration_List")
		return err
	}

	// set the config list into the config response
	err = config.SetRaft(slice)
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error setting config")
		return err
	}

	// map it!
	return res.SetResponse(config)
}

func (c *ConfigServer) getAllRaftConfigs(ctx context.Context, call configv1.ConfigService_getConfig) error {

	// allocate the results
	res, err := call.AllocResults()
	if err != nil {
		c.logger.Err(err).Msg("error allocating getAllRaftConfigs results")
		return err
	}

	// acknowledge the request now that we can allocate space for the results
	call.Ack()

	// find the target
	val, err := c.raftManager.GetAll()
	if err != nil {
		c.logger.Err(err).Msg("error getting all configs")
		return err
	}

	// prep the results
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		c.logger.Err(err).Msg("error creating message")
		return err
	}

	// create the root config response
	config, err := configv1.NewRootGetConfigurationResponse(seg)
	if err != nil {
		c.logger.Err(err).Msg("error creating getRaftConfigurationResponse")
		return err
	}

	// create the config lists
	slice, err := configv1.NewRaftConfiguration_List(seg, int32(len(val)))
	if err != nil {
		c.logger.Err(err).Msg("error creating raftConfiguration_List")
		return err
	}

	// set the config into the slice
	for idx := range val {
		err = slice.Set(0, *val[idx])
		if err != nil {
			c.logger.Err(err).Msg("error setting raftConfiguration_List")
			return err
		}
	}

	// set the config list into the config response
	err = config.SetRaft(slice)
	if err != nil {
		c.logger.Err(err).Msg("error setting config")
		return err
	}

	// map it!
	return res.SetResponse(config)
}

func (c *ConfigServer) PutConfig(ctx context.Context, call configv1.ConfigService_putConfig) error {
	// acknowledge the request now that we can allocate space for the results
	call.Ack()

	// get the request
	req, err := call.Args().Request()
	if err != nil {
		c.logger.Err(err).Msg("error getting request")
		return err
	}

	switch req.Which() {
	case configv1.PutConfigurationRequest_Which_raft:
		return c.putRaftConfig(ctx, call)
	case configv1.PutConfigurationRequest_Which_nodeHost:
	}

	return nil
}

func (c *ConfigServer) putRaftConfig(ctx context.Context, call configv1.ConfigService_putConfig) error {
	res, err := call.AllocResults()
	if err != nil {
		c.logger.Err(err).Msg("error allocating putRaftConfig results")
		return err
	}

	call.Ack()

	req, err := call.Args().Request()
	if err != nil {
		c.logger.Err(err).Msg("error getting request")
		return err
	}

	raftConfig, err := req.Raft()
	if err != nil {
		c.logger.Err(err).Msg("error getting raft config")
		return err
	}

	configId, err := raftConfig.Id()
	if err != nil {
		c.logger.Err(err).Msg("error getting raft config id")
		return err
	}

	err = c.raftManager.Put(configId, &raftConfig)
	if err != nil {
		c.logger.Err(err).Msg("error putting raft config")
		return err
	}

	configResponse, err := res.NewResponse()
	if err != nil {
		c.logger.Err(err).Msg("error building response")
		return err
	}

	err = configResponse.SetRaft(raftConfig)
	if err != nil {
		c.logger.Err(err).Msg("error setting raft config")
		return err
	}
	configResponse.SetSuccess(true)
	err = configResponse.SetStatus("")
	if err != nil {
		c.logger.Err(err).Msg("error setting status")
		return err
	}

	return res.SetResponse(configResponse)
}
