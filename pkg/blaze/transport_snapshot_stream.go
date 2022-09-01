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
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/lni/dragonboat/v3/raftpb"
	"github.com/rs/zerolog"
)

const (
	RaftSnapshotProtocolVersion  protocol.ID = "/pleiades/raft-snapshot/0.0.1"
)

var (
	_ raftio.ISnapshotConnection = (*RaftSnapshotConnectionStream)(nil)
)

func NewRaftSnapshotConnectionStream(stream network.Stream, logger zerolog.Logger) *RaftSnapshotConnectionStream {
	return &RaftSnapshotConnectionStream{
		logger: logger,
		stream: stream,
	}
}

type RaftSnapshotConnectionStream struct{
	logger         zerolog.Logger
	stream         network.Stream
}

func (r *RaftSnapshotConnectionStream) Close() {
	if err := r.stream.Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close stream")
	}
}

func (r *RaftSnapshotConnectionStream) SendChunk(chunk raftpb.Chunk) error {
	return nil
}
