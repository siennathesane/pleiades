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
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/lni/dragonboat/v3/raftpb"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

var (
	_ raftio.ITransport = (*raftTransport)(nil)
)

const (
	messageService MethodByte = 0x0
	chunkService   MethodByte = 0x1
)

type raftTransport struct {
	host           host.Host
	logger         zerolog.Logger
	messageHandler raftio.MessageHandler
	chunkHandler   raftio.ChunkHandler
}

func (r *raftTransport) Name() string {
	return string(raftTransportProtocolVersion)
}

func (r *raftTransport) Start() error {
	r.host.SetStreamHandler(raftTransportProtocolVersion, r.connectionStreamHandler)
	return nil
}

func (r *raftTransport) connectionStreamHandler(stream network.Stream) {
	if err := stream.Scope().SetService(raftTransportService); err != nil {
		_ = stream.Reset()
	}

	for {
		frame := NewFrame()
		_, err := frame.ReadFrom(stream)
		if err != nil {
			// todo (sienna): add error handling
			r.logger.Error().Err(err).Msg("cannot read frame")
		}

		msgBuf := frame.GetPayload()
		if err != nil {
			// todo (sienna): add error handling
			r.logger.Error().Err(err).Msg("cannot get payload")
		}

		if frame.GetService() != RaftControlServiceByte {
			// todo (sienna): add error handling
			r.logger.Error().Err(err).Msg("invalid service")
			_ = stream.Reset()
		}

		switch frame.GetMethod() {
		case messageService:
			batch := raftpb.MessageBatch{}
			if err := batch.Unmarshal(msgBuf); err != nil {
				r.logger.Error().Err(err).Msg("cannot unmarshal message batch")
				_ = stream.Reset()
			}
			r.logger.Debug().Msg("received message batch")
			r.messageHandler(batch)
			break
		case chunkService:
			chunk := raftpb.Chunk{}
			if err := chunk.Unmarshal(msgBuf); err != nil {
				r.logger.Error().Err(err).Msg("cannot unmarshal chunk")
				_ = stream.Reset()
			}
			r.logger.Debug().Msg("received chunk")
			r.chunkHandler(chunk)
			break
		}
	}
}

func (r *raftTransport) Stop() {
	r.host.RemoveStreamHandler(raftTransportProtocolVersion)
	if err := r.host.Network().Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close network")
	}
}

func (r *raftTransport) GetConnection(ctx context.Context, target string) (raftio.IConnection, error) {
	r.logger.Info().Str("target", target).Msg("getting connection")
	targetHost, err := multiaddr.NewMultiaddr(target)
	if err != nil {
		r.logger.Error().Str("target", target).Err(err).Msg("cannot create multiaddr from target")
		return nil, err
	}
	return newRaftConnection(ctx, r.host, targetHost, r.logger)
}

func (r *raftTransport) GetSnapshotConnection(ctx context.Context, target string) (raftio.ISnapshotConnection, error) {
	r.logger.Info().Str("target", target).Msg("getting connection")
	targetHost, err := multiaddr.NewMultiaddr(target)
	if err != nil {
		r.logger.Error().Str("target", target).Err(err).Msg("cannot create multiaddr from target")
		return nil, err
	}
	return newRaftSnapshotConnectionStream(ctx, r.host, targetHost, r.logger)
}
