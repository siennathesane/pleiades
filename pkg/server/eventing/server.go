/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package eventing

import (
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/rs/zerolog"
)

var (
	serverSingleton *Server
)

func NewServer(logger zerolog.Logger) (*Server, error) {
	if serverSingleton != nil {
		return serverSingleton, nil
	}

	srv, err := messaging.NewEmbeddedMessagingWithDefaults(logger)
	if err != nil {
		return nil, err
	}

	srv.Start()

	serverSingleton = &Server{srv, logger.With().Str("component", "eventing").Logger()}

	return serverSingleton, nil
}

type Server struct {
	*messaging.EmbeddedMessaging
	logger zerolog.Logger
}

func (s *Server) GetRaftEventHandler() (*messaging.RaftEventHandler, error) {
	pubSubClient, err := s.EmbeddedMessaging.GetPubSubClient()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't create pubsub client")
		return nil, err
	}

	queueClient, err := s.EmbeddedMessaging.GetStreamClient()
	if err != nil {
		s.logger.Error().Err(err).Msg("can't create stream client")
		return nil, err
	}

	return messaging.NewRaftEventHandler(pubSubClient, queueClient, s.logger), nil
}