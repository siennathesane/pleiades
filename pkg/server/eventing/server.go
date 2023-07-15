/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package eventing

import (
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	serverSingleton *EventServer
)

type EventServerBuilderParams struct {
	fx.In

	Logger            zerolog.Logger
	EmbeddedMessaging *messaging.EmbeddedMessaging
}

type EventServerBuilderResults struct {
	fx.Out
	Server *EventServer
}

func NewEventServer(params EventServerBuilderParams) EventServerBuilderResults {
	serverSingleton = &EventServer{
		params.EmbeddedMessaging,
		params.Logger.With().Str("component", "eventing").Logger(),
	}

	return EventServerBuilderResults{
		Server: serverSingleton,
	}
}

type EventServer struct {
	*messaging.EmbeddedMessaging
	logger zerolog.Logger
}

func (s *EventServer) GetRaftEventHandler() (*messaging.RaftEventHandler, error) {
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

func (s *EventServer) GetRaftSystemEventListener() (*messaging.RaftSystemListener, error) {
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

	return messaging.NewRaftSystemListener(pubSubClient, queueClient, s.logger)
}

type NewPubSubClientBuilderParams struct {
	fx.In

	Server *EventServer
}

func NewPubSubClient(params NewPubSubClientBuilderParams) (*messaging.EmbeddedMessagingPubSubClient, error) {
	return params.Server.GetPubSubClient()
}

type NewStreamClientBuilderParams struct {
	fx.In

	Server *EventServer
}

func NewStreamClient(params NewStreamClientBuilderParams) (*messaging.EmbeddedMessagingStreamClient, error) {
	return params.Server.GetStreamClient()
}
