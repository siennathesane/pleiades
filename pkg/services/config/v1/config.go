/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package v1

import (
	"context"
	"errors"

	"capnproto.org/go/capnp/v3"
	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/fsm"
	configv1 "r3t.io/pleiades/pkg/protocols/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
)

var (
	_ configv1.ConfigService_Server = (*ConfigService)(nil)
)

type ConfigService struct {
	store       *services.StoreManager
	logger      zerolog.Logger
	raftManager *fsm.ConfigServiceStoreManager
}

func NewConfigService(store *services.StoreManager, logger zerolog.Logger) (*ConfigService, error) {
	l := logger.With().Str("service", "config").Logger()
	manager, err := fsm.NewConfigServiceStoreManager(logger, store)
	if err != nil {
		return nil, err
	}

	return &ConfigService{
		store:       store,
		logger:      l,
		raftManager: manager,
	}, nil
}

func (c *ConfigService) GetConfig(ctx context.Context, call configv1.ConfigService_getConfig) error {
	req, err := call.Args().Request()
	if err != nil {
		return err
	}

	switch req.What() {
	case configv1.GetConfigurationRequest_Type_all:
	case configv1.GetConfigurationRequest_Type_raft:
		switch req.Amount() {
		case configv1.GetConfigurationRequest_Specificity_one:
			key, err := req.Id()
			if err != nil {
				c.logger.Err(err).Str("key", key).Msg("error reading key")
				return err
			}
			return c.getRaftConfig(ctx, key, call)
		}
	}

	return nil
}

func (c *ConfigService) getRaftConfig(ctx context.Context, key string, call configv1.ConfigService_getConfig) error {
	if key == "" {
		return errors.New("cannot request a named record without a key")
	}

	// allocate the results
	res, err := call.AllocResults()
	if err != nil {
		c.logger.Err(err).Str("key", key).Msg("error allocating results")
		return err
	}

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
