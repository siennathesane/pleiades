
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
	"fmt"
	"net"
	"testing"

	"github.com/mxplusb/pleiades/pkg/fsm"
	configv1 "github.com/mxplusb/pleiades/pkg/protocols/v1/config"
	"github.com/mxplusb/pleiades/pkg/services"
	"github.com/mxplusb/pleiades/pkg/utils"
	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestConfigService(t *testing.T) {
	suite.Run(t, new(ConfigServiceTests))
}

type ConfigServiceTests struct {
	suite.Suite
	logger      zerolog.Logger
	store       *services.StoreManager
	raftManager *fsm.ConfigServiceStoreManager
}

func (cst *ConfigServiceTests) SetupSuite() {
	cst.logger = utils.NewTestLogger(cst.T())

	var err error
	cst.store = services.NewStoreManager(cst.T().TempDir(), cst.logger)
	cst.raftManager, err = fsm.NewConfigServiceStoreManager(cst.logger, cst.store)
	cst.Require().NoError(err, "there must not be an error creating the raft manager")
}

func (cst *ConfigServiceTests) TestConfigService_GetRaft() {
	for i := 0; i < 10; i++ {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		cst.Require().NoError(err, "there must not be an error creating a message")

		conf, err := configv1.NewRootRaftConfiguration(seg)
		cst.Require().NoError(err, "there must not be an error creating a configuration")

		testId := fmt.Sprintf("test-%d", i)
		err = conf.SetId(testId)
		cst.Require().NoError(err, "there must not be an error setting the id")

		err = cst.raftManager.Put(testId, &conf)
		cst.Require().NoError(err, "there must not be an error putting the configuration")
	}

	configServiceImpl, err := NewConfigServer(cst.store, cst.logger)
	cst.Require().NoError(err, "there must not be an error creating a config service")
	cst.Require().NotNil(configServiceImpl, "the config service must not be nil")

	input, output := net.Pipe()

	clientFactory := configv1.ConfigService_ServerToClient(configServiceImpl, &server.Policy{MaxConcurrentCalls: 10})
	srvConn := rpc.NewConn(rpc.NewStreamTransport(input), &rpc.Options{
		BootstrapClient: clientFactory.Client,
	})

	clientConn := rpc.NewConn(rpc.NewStreamTransport(output), nil)

	ctx := context.Background()
	client := configv1.ConfigService{Client: clientConn.Bootstrap(ctx)}
	cst.Require().NotNil(client, "the client must not be nil")

	resp, rel := client.GetConfig(ctx, func(getConfigParams configv1.ConfigService_getConfig_Params) error {
		req, err := getConfigParams.NewRequest()
		if err != nil {
			return err
		}

		req.SetWhat(configv1.GetConfigurationRequest_Type_raft)
		req.SetAmount(configv1.GetConfigurationRequest_Specificity_one)
		if err := req.SetId("test-1"); err != nil {
			return err
		}
		return getConfigParams.SetRequest(req)
	})

	configResponse := resp.Response()
	cst.Require().NotNil(configResponse, "the response must not be nil")

	results, err := configResponse.Struct()
	cst.Require().NoError(err, "there must not be an error getting the results")
	cst.Require().NotNil(results, "the results must not be nil")

	ok := results.HasRaft()
	cst.Require().True(ok, "the results must have a raft configuration")

	raftConfig, err := results.Raft()
	cst.Require().NoError(err, "there must not be an error getting the raft configuration")
	cst.Require().NotNil(raftConfig, "the raft configuration must not be nil")
	cst.Require().Equal(raftConfig.Len(), 1, "the raft configuration must have a length of 1")

	config := raftConfig.At(0)
	cst.Require().NotNil(config, "the raft configuration must not be nil")

	id, err := config.Id()
	cst.Require().NoError(err, "there must not be an error getting the id")
	cst.Require().Equal(id, "test-1", "the raft configuration must have an id of 'test'")

	cst.Require().NotPanics(func() { rel() }, "there must not be an error releasing the response")

	resp, rel = client.GetConfig(ctx, func(getConfigParams configv1.ConfigService_getConfig_Params) error {
		req, err := getConfigParams.NewRequest()
		if err != nil {
			return err
		}

		req.SetWhat(configv1.GetConfigurationRequest_Type_raft)
		req.SetAmount(configv1.GetConfigurationRequest_Specificity_everything)
		return getConfigParams.SetRequest(req)
	})

	configResponse = resp.Response()
	cst.Require().NotNil(configResponse, "the response must not be nil")

	results, err = configResponse.Struct()
	cst.Require().NoError(err, "there must not be an error getting the results")
	cst.Require().NotNil(results, "the results must not be nil")

	ok = results.HasRaft()
	cst.Require().True(ok, "the results must have a raft configuration")

	raftConfig, err = results.Raft()
	cst.Require().NoError(err, "there must not be an error getting the raft configuration")
	cst.Require().NotNil(raftConfig, "the raft configuration must not be nil")
	cst.Require().Equal(raftConfig.Len(), 10, "the raft configuration must have a length of 10")

	err = clientConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the client pipe")

	err = srvConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the server pipe")
}

func (cst *ConfigServiceTests) TestConfigService_PutRaft() {
	configServiceImpl, err := NewConfigServer(cst.store, cst.logger)
	cst.Require().NoError(err, "there must not be an error creating a config service")
	cst.Require().NotNil(configServiceImpl, "the config service must not be nil")

	input, output := net.Pipe()

	clientFactory := configv1.ConfigService_ServerToClient(configServiceImpl, &server.Policy{MaxConcurrentCalls: 10})
	srvConn := rpc.NewConn(rpc.NewStreamTransport(input), &rpc.Options{
		BootstrapClient: clientFactory.Client,
	})

	clientConn := rpc.NewConn(rpc.NewStreamTransport(output), nil)

	ctx := context.Background()
	client := configv1.ConfigService{Client: clientConn.Bootstrap(ctx)}
	cst.Require().NotNil(client, "the client must not be nil")

	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	cst.Require().NoError(err, "there must not be an error creating a message")

	raftConf, err := configv1.NewRootRaftConfiguration(seg)
	cst.Require().NoError(err, "there must not be an error creating a configuration")

	err = raftConf.SetId("test-put")
	cst.Require().NoError(err, "there must not be an error setting the id")

	resp, free := client.PutConfig(ctx, func(params configv1.ConfigService_putConfig_Params) error {
		req, err := params.NewRequest()
		if err != nil {
			return err
		}
		err = req.SetRaft(raftConf)
		if err != nil {
			return err
		}

		return params.SetRequest(req)
	})

	results := resp.Response()
	cst.Require().NotNil(results, "the response must not be nil")

	configResponse, err := results.Struct()
	cst.Require().NoError(err, "there must not be an error getting the rpc results")
	cst.Require().NotNil(configResponse, "the config response must not be nil")

	raftConfig, err := configResponse.Raft()
	cst.Require().NoError(err, "there must not be an error getting the raft configuration")
	cst.Require().NotNil(raftConfig, "the raft configuration must not be nil")

	id, err := raftConfig.Id()
	cst.Require().NoError(err, "there must not be an error getting the id")
	cst.Require().Equal(id, "test-put", "the fetched results must be the same")

	free()

	err = clientConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the client pipe")

	err = srvConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the server pipe")
}
