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

	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	configv1 "r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/services/v1/config"
)

const (
	ConfigService = configv1.ServiceType_Type_configService
	TestService   = configv1.ServiceType_Type_test
)

// StreamReceiver takes a new stream and returns a *server.Server reference from the registry
type StreamReceiver struct {
	logger   zerolog.Logger
	registry *config.Registry
}

func NewStreamReceiver(logger zerolog.Logger, registry *config.Registry) (*StreamReceiver, error) {
	l := logger.With().Str("component", "stream_receiver").Logger()

	return &StreamReceiver{logger: l, registry: registry}, nil
}

func (sr *StreamReceiver) Receive(ctx context.Context, logger zerolog.Logger, stream quic.Stream) {
	localLogger := logger.With().Int64("stream-id", int64(stream.StreamID())).Logger()

	neg := config.NewNegotiator(localLogger, sr.registry)

	clientFactory := configv1.Negotiator_ServerToClient(neg, &server.Policy{
		MaxConcurrentCalls: 250,
	})

	serverConn := rpc.NewConn(rpc.NewStreamTransport(stream), &rpc.Options{
		BootstrapClient: clientFactory.Client,
	})

	select {
	case <-ctx.Done():
		return
	case <-ctx.Done():
		err := serverConn.Close()
		if err != nil {
			localLogger.Err(err).Msg("error closing connection")
		}
		return
	}
}
