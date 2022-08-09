
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"context"
	"testing"

	hostv1 "github.com/mxplusb/pleiades/pkg/protocols/v1/host"
	"github.com/mxplusb/pleiades/pkg/services/v1/config"
	"github.com/mxplusb/pleiades/pkg/utils"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// StreamReceiverTest tests the StreamReceiver class.
func TestStreamReceiver(t *testing.T) {
	suite.Run(t, new(StreamReceiverTest))
}

type StreamReceiverTest struct {
	suite.Suite
	tk       *QuicTestKit
	logger   zerolog.Logger
	registry *config.Registry
}

func (s *StreamReceiverTest) SetupSuite() {
	s.logger = utils.NewTestLogger(s.T())
	s.tk = NewQuicTestKit(s.T())

	var err error
	s.registry, err = config.NewRegistry(s.logger)
	s.Require().NoError(err, "failed to create registry")

	err = s.registry.PutServer(
		hostv1.ServiceType_Type_test,
		&server.Server{},
	)
	s.Require().NoError(err, "failed to register test server")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver() {

	sr, err := NewStreamReceiver(s.logger, s.registry)
	s.Require().NoError(err, "there must not be an error when creating a stream receiver")

	conn, err := s.tk.GetListener().Accept(context.Background())
	s.Require().NoError(err, "there must not be an error when accepting a connection")
	s.Require().NotNil(conn, "the connection must not be nil")

	ctx := context.Background()
	go func(ctx context.Context, s *StreamReceiverTest) {
		for {
			receivingTestStream, err := conn.AcceptStream(ctx)
			s.Require().NoError(err, "there must not be an error when accepting a stream")
			s.Require().NotNil(receivingTestStream, "the stream must not be nil")
			sr.Receive(ctx, s.logger, receivingTestStream)
		}
	}(ctx, s)

	dialer := s.tk.GetConnection()
	stream, err := dialer.OpenStream()
	s.Require().NoError(err, "there must not be an error when opening a stream")

	rpcConn := rpc.NewConn(rpc.NewStreamTransport(stream), nil)
	s.Require().NotNil(rpcConn, "the rpc connection must not be nil")

	client := hostv1.Negotiator{Client: rpcConn.Bootstrap(ctx)}
	s.Require().NotNil(client, "the client must not be nil")

	future, free := client.ConfigService(ctx, nil)
	configService := future.Svc()
	s.Require().NotNil(configService, "the host service must not be nil")

	configSrv := configService.AddRef()
	s.Require().NotNil(configSrv, "the host service must not be nil")

	free()
	s.Require().NotNil(configSrv, "the host service reference must not be nil after calling free()")

	ctx.Done()
}
