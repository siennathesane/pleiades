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
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/lni/dragonboat/v3/raftpb"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

const (
	raftSnapshotProtocolVersion protocol.ID = "/pleiades/raft-snapshot/0.0.1"
)

var (
	_ raftio.ISnapshotConnection = (*raftSnapshotConnectionStream)(nil)
)

func newRaftSnapshotConnectionStream(ctx context.Context, host host.Host, addr multiaddr.Multiaddr, logger zerolog.Logger) (*raftSnapshotConnectionStream, error) {
	l := logger.With().Str("component", "raft-snapshot-connection").Str("target", addr.String()).Logger()
	rc := &raftSnapshotConnectionStream{
		logger: l,
		host:   host,
		target: addr,
	}

	targetInfo, err := peer.AddrInfosFromP2pAddrs(addr)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse target address")
		return nil, err
	}

	// this is safe to use the 0th value due to only having a single address
	rc.stream, err = rc.host.NewStream(ctx, targetInfo[0].ID, raftTransportProtocolVersion)
	if err != nil {
		l.Error().Err(err).Msg("failed to open stream")
		return nil, err
	}

	return rc, nil
}

type raftSnapshotConnectionStream struct {
	host   host.Host
	logger zerolog.Logger
	target multiaddr.Multiaddr
	stream network.Stream
}

func (r *raftSnapshotConnectionStream) Close() {
	if err := r.stream.Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close stream")
	}
}

func (r *raftSnapshotConnectionStream) SendChunk(chunk raftpb.Chunk) error {
	frame := NewFrame().WithService(RaftTransportService).WithMethod(chunkService)

	buf, err := chunk.Marshal()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal message batch")
		return err
	}

	frame.WithPayload(buf)
	msgBuf, err := frame.Marshal()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal frame")
		return err
	}

	_, err = r.stream.Write(msgBuf)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to write frame")
		return err
	}

	return nil
}
