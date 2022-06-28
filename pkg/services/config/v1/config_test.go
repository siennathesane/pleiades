/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package v1

import (
	"context"
	"net"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/fsm"
	v1 "r3t.io/pleiades/pkg/protocols/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
	"r3t.io/pleiades/pkg/utils"
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

func (cst *ConfigServiceTests) TestConfigService() {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	cst.Require().NoError(err, "there must not be an error creating a message")

	conf, err := v1.NewRootRaftConfiguration(seg)
	cst.Require().NoError(err, "there must not be an error creating a configuration")
	err = conf.SetId("test")
	cst.Require().NoError(err, "there must not be an error setting the id")

	err = cst.raftManager.Put("test", &conf)
	cst.Require().NoError(err, "there must not be an error putting the configuration")

	configServiceImpl, err := NewConfigService(cst.store, cst.logger)
	cst.Require().NoError(err, "there must not be an error creating a config service")
	cst.Require().NotNil(configServiceImpl, "the config service must not be nil")

	input, output := net.Pipe()

	clientFactory := v1.ConfigService_ServerToClient(configServiceImpl, &server.Policy{MaxConcurrentCalls: 10})
	srvConn := rpc.NewConn(rpc.NewStreamTransport(input), &rpc.Options{
		BootstrapClient: clientFactory.Client,
	})

	clientConn := rpc.NewConn(rpc.NewStreamTransport(output), nil)

	ctx := context.Background()
	client := v1.ConfigService{Client: clientConn.Bootstrap(ctx)}
	cst.Require().NotNil(client, "the client must not be nil")

	resp, rel := client.GetConfig(ctx, func(getConfigParams v1.ConfigService_getConfig_Params) error {
		_, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
		req, err := v1.NewRootGetConfigurationRequest(s)
		if err != nil {
			return err
		}

		req.SetWhat(v1.GetConfigurationRequest_Type_raft)
		req.SetAmount(v1.GetConfigurationRequest_Specificity_one)
		if err := req.SetId("test"); err != nil {
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
	cst.Require().Equal(id, "test", "the raft configuration must have an id of 'test'")

	cst.Require().NotPanics(func() { rel() }, "there must not be an error releasing the response")

	err = clientConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the client pipe")

	err = srvConn.Close()
	cst.Require().NoError(err, "there must not be an error closing the server pipe")
}
