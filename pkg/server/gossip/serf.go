/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package gossip

import (
	"context"
	"sync"

	"github.com/hashicorp/serf/serf"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type GossipServerBuilderParams struct {
	fx.In

	PubSubClient *messaging.EmbeddedMessagingPubSubClient
	Logger       zerolog.Logger
}

func NewServer(lc fx.Lifecycle, params GossipServerBuilderParams) (*Server, error) {
	s := &Server{
		serfConfig: serf.DefaultConfig(),
	}

	// set the loggers
	s.serfConfig.MemberlistConfig.LogOutput = s.logger
	s.serfConfig.MemberlistConfig.EnableCompression = true
	s.serfConfig.LogOutput = s.logger

	var err error
	s.serf, err = serf.Create(s.serfConfig)
	if err != nil {
		return &Server{}, err
	}

	return s, nil
}

type Server struct {
	logger            zerolog.Logger
	serf              *serf.Serf
	serfConfig        *serf.Config
	eventChan         chan serf.Event
	eventHandlers     map[string]struct{}
	eventHandlersLock sync.Mutex
	shutdown          bool
	shutdownCh        chan struct{}
	shutdownLock      sync.Mutex
}

func (s *Server) StartHook(ctx context.Context) error {
	return nil
}
