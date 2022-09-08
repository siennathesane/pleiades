/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

const (
	raftTransportService         string      = "raft-transport.pleiades"
	raftTransportProtocolVersion protocol.ID = "/pleiades/raft-transport/0.0.1"
)

var (
	_ dconfig.TransportFactory = (*RaftTransportFactory)(nil)
)

func NewRaftTransportFactory(host host.Host, logger zerolog.Logger) *RaftTransportFactory {
	return &RaftTransportFactory{
		host:   host,
		logger: logger,
	}
}

type RaftTransportFactory struct {
	host   host.Host
	logger zerolog.Logger
}

func (r *RaftTransportFactory) Create(_ dconfig.NodeHostConfig, handler raftio.MessageHandler, handler2 raftio.ChunkHandler) raftio.ITransport {
	return &raftTransport{
		logger:         r.logger.With().Str("component", "transport").Logger(),
		messageHandler: handler,
		chunkHandler:   handler2,
		host:           r.host,
	}
}

func (r *RaftTransportFactory) Validate(s string) bool {
	if !strings.Contains(s, "/") {
		return false
	}
	_, err := multiaddr.NewMultiaddr(s)
	return err == nil
}
