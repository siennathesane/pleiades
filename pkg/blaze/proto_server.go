/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/services/v1/config"
)

type ProtoServer struct {
	logger zerolog.Logger
	registry *config.Registry
	streamReceiver *StreamReceiver
}

func NewProtoServer(logger zerolog.Logger) *ProtoServer {
	l := logger.With().Str("component", "server").Logger()

	reg, err := config.NewRegistry(l)
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
