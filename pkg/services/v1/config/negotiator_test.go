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
	"r3t.io/pleiades/pkg/servers/services"
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

	configService := configServiceResponse.Svc()
	nt.Require().NotNil(configService, "the configService must not be nil")

	free()
}
