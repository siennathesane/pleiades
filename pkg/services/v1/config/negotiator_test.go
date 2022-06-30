/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package config

import (
	"context"
	"net"
	"testing"

	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/services"
	"r3t.io/pleiades/pkg/utils"
)

func TestNegotiator(t *testing.T) {
	suite.Run(t, new(NegotiatorTests))
}

type NegotiatorTests struct {
	suite.Suite
	logger       zerolog.Logger
	registry     *Registry
	store        *services.StoreManager
	configServer *ConfigServer
}

func (nt *NegotiatorTests) SetupSuite() {
	nt.logger = utils.NewTestLogger(nt.T())
	var err error
	nt.registry, err = NewRegistry(nt.logger)
	nt.Require().NoError(err, "there must not be an error creating the Registry")

	nt.store = services.NewStoreManager(nt.T().TempDir(), nt.logger)
	nt.Require().NotNil(nt.store, "the store must not be nil")

	nt.configServer, err = NewConfigServer(nt.store, nt.logger)
	err = nt.registry.PutServer(config.ServiceType_Type_configService, nt.configServer)
	nt.Require().NoError(err, "there must not be an error adding the config server to the registry")
}

func (nt *NegotiatorTests) TestConfigServerReturn() {
	neg := NewNegotiator(nt.logger, nt.registry)
	nt.Require().NotNil(neg, "the negotiator must not be nil")

	serverPipe, clientPipe := net.Pipe()

	clientFactory := config.Negotiator_ServerToClient(neg, &server.Policy{
		MaxConcurrentCalls: 250,
	})

	serverConn := rpc.NewConn(rpc.NewStreamTransport(serverPipe), &rpc.Options{
		BootstrapClient: clientFactory.Client,
	})
	defer serverConn.Close()

	clientConn := rpc.NewConn(rpc.NewStreamTransport(clientPipe), nil)

	ctx := context.Background()
	client := config.Negotiator{Client: clientConn.Bootstrap(ctx)}
	nt.Require().NotNil(client, "the client must not be nil")

	configServiceResponse, free := client.ConfigService(ctx, nil)
	nt.Require().NotNil(configServiceResponse, "the configServiceResponse must not be nil")
	defer free()

	configService := configServiceResponse.Svc()
	nt.Require().NotNil(configService, "the configService must not be nil")

	resp, free := configService.PutConfig(ctx, func(params config.ConfigService_putConfig_Params) error {
		req, err := params.NewRequest()
		if err != nil {
			nt.logger.Error().Err(err).Msg("failed to allocate request")
			return err
		}

		raft, err := req.NewRaft()
		if err != nil {
			nt.logger.Error().Err(err).Msg("failed to allocate raft")
			return err
		}
		err = raft.SetId("test")
		if err != nil {
			nt.logger.Error().Err(err).Msg("failed to set raft id")
			return err
		}

		err = req.SetRaft(raft)
		if err != nil {
			nt.logger.Error().Err(err).Msg("failed to set raft")
			return err
		}

		return params.SetRequest(req)
	})

	res := resp.Response()
	nt.Require().NotNil(res, "the response must not be nil")

	response, err := res.Struct()
	nt.Require().NoError(err, "there must not be an error getting the response")
	nt.Require().True(response.HasRaft(), "the response must have a raft configuration")
	nt.Require().True(response.Success(), "putting a configuration must be successful")

	raftConfigs, err := response.Raft()
	nt.Require().NoError(err, "there must not be an error getting the raft configuration")

	id, err := raftConfigs.Id()
	nt.Require().NoError(err, "there must not be an error getting the raft id")
	nt.Require().Equal("test", id, "the raft id must be equal")
}
