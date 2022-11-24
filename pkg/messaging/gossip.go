/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package messaging

import (
	"path/filepath"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
)

var (
	gossipSingleton *EmbeddedGossipServer
)

type EmbeddedGossipServer struct {
	s      *serf.Serf
	logger zerolog.Logger
}

func NewEmbeddedGossip(logger zerolog.Logger) (*EmbeddedGossipServer, error) {
	l := logger.With().Str("component", "gossip").Logger()

	if gossipSingleton != nil {
		l.Debug().Msg("using singleton")
		return gossipSingleton, nil
	}

	sc := serf.DefaultConfig()
	sc.SnapshotPath = filepath.Join(configuration.Get().GetString("server.datastore.basePath"), "gossip.db")
	sc.MemberlistConfig.BindPort = configuration.Get().GetInt("server.gossip.port")

	s, err := serf.Create(sc)
	if err != nil {
		l.Error().Err(err).Msg("failed to create serf config")
		return nil, err
	}

	gossipSingleton = &EmbeddedGossipServer{
		s:      s,
		logger: l,
	}

	config := configuration.Get()
	localeTags := map[string]string{
		"continent": config.GetString("server.gossip.continent"),
		"region":    config.GetString("server.gossip.region"),
		"zone":      config.GetString("server.gossip.zone"),
	}
	l.Debug().Interface("locale-tags", localeTags).Msg("adding locale tags")

	err = gossipSingleton.AddTags(localeTags)
	if err != nil {
		l.Error().Err(err).Msg("failed to add locale tags")
		return nil, err
	}

	return gossipSingleton, nil
}

func (e *EmbeddedGossipServer) Stop() error {
	if err := e.s.Leave(); err != nil {
		e.logger.Error().Err(err).Msg("can't safely leave memberlist")
		return err
	}
	return nil
}

func (e *EmbeddedGossipServer) AddTags(tags map[string]string) error {
	return e.s.SetTags(tags)
}
