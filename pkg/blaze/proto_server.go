/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
)

type ProtoServer struct {
	logger zerolog.Logger
	registry *Registry
	streamReceiver *StreamReceiver
}

func NewProtoServer(logger zerolog.Logger) *ProtoServer {
	l := logger.With().Str("component", "server").Logger()

	reg, err := NewRegistry(l)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create registry")
	}

	streamReceiver, err := NewStreamReceiver(l, reg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create stream receiver")
	}

	return &ProtoServer{
		logger: logger,
		registry: reg,
		streamReceiver: streamReceiver,
	}
}

func (ps *ProtoServer) Register(name string, srv *server.Server) error {
	return ps.registry.Put(name, srv)
}


