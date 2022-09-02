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

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/rs/zerolog"
)

const (
	RaftTransportProtocolVersion protocol.ID = "/pleiades/raft-transport/0.0.1"
)

var (
	_ raftio.ITransport = (*RaftTransport)(nil)
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

func (r *RaftTransportFactory) Create(config dconfig.NodeHostConfig, handler raftio.MessageHandler, handler2 raftio.ChunkHandler) raftio.ITransport {
	//TODO implement me
	panic("implement me")
}

func (r *RaftTransportFactory) Validate(s string) bool {
	//TODO implement me
	panic("implement me")
}

type RaftTransport struct {
	host           host.Host
	logger         zerolog.Logger
	messageHandler raftio.MessageHandler
	chunkHandler   raftio.ChunkHandler
}

func (r *RaftTransport) Name() string {
	return string(RaftTransportProtocolVersion)
}

func (r *RaftTransport) Start() error {
	r.host.SetStreamHandler(RaftStreamProtocolVersion, r.connectionStreamHandler)
	return nil
}

func (r *RaftTransport) connectionStreamHandler(stream network.Stream) {
	streamer := &RaftConnectionStream{
		logger:         r.logger.With().Str("peer", stream.Conn().RemotePeer().String()).Logger(),
		messageHandler: r.messageHandler,
		chunkHandler:   r.chunkHandler,
		stream:         stream,
	}

	streamer.Serve()
}

func (r *RaftTransport) Stop() {

	r.host.RemoveStreamHandler(RaftStreamProtocolVersion)
	if err := r.host.Network().Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close network")
	}
}

func (r *RaftTransport) GetConnection(ctx context.Context, target string) (raftio.IConnection, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RaftTransport) GetSnapshotConnection(ctx context.Context, target string) (raftio.ISnapshotConnection, error) {
	//TODO implement me
	panic("implement me")
}
