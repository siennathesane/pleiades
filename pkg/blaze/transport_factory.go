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
	RaftTransportService string = "raft-transport.pleiades"
)

var (
	_ raftio.ITransport = (*raftTransport)(nil)
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

type raftTransport struct {
	host           host.Host
	logger         zerolog.Logger
	messageHandler raftio.MessageHandler
	chunkHandler   raftio.ChunkHandler
}

func (r *raftTransport) Name() string {
	return string(RaftTransportProtocolVersion)
}

func (r *raftTransport) Start() error {
	r.host.SetStreamHandler(RaftStreamProtocolVersion, r.connectionStreamHandler)
	return nil
}

func (r *raftTransport) connectionStreamHandler(stream network.Stream) {
	if err := stream.Scope().SetService(RaftTransportService); err != nil {
		_ = stream.Reset()
	}

	for {
		frame := NewFrame()
		_, err := frame.ReadFrom(stream)
		if err != nil {
			// todo (sienna): add error handling
			r.logger.Error().Err(err).Msg("cannot read frame")
		}

		msgBuf, err := frame.GetPayload()
		if err != nil {
			// todo (sienna): add error handling
			r.logger.Error().Err(err).Msg("cannot get payload")
		}
	}
}

func (r *raftTransport) Stop() {
	r.host.RemoveStreamHandler(RaftStreamProtocolVersion)
	if err := r.host.Network().Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close network")
	}
}

func (r *raftTransport) GetConnection(ctx context.Context, target string) (raftio.IConnection, error) {
	//TODO implement me
	panic("implement me")
}

func (r *raftTransport) GetSnapshotConnection(ctx context.Context, target string) (raftio.ISnapshotConnection, error) {
	//TODO implement me
	panic("implement me")
}
