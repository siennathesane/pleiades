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
	"hash/crc32"
	"io"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/lni/dragonboat/v3/raftpb"
	"github.com/rs/zerolog"
)

const (
	RaftStreamProtocolVersion protocol.ID = "/pleiades/raft-stream/0.0.1"
)

var (
	// todo (sienna): implement data signing
	_ raftio.IConnection = (*RaftConnectionStream)(nil)

	// ErrBadMessage is the error returned to indicate the incoming message is
	// corrupted.
	ErrBadMessage       = errors.New("invalid message")
	errPoisonReceived   = errors.New("poison received")
	MagicNumber         = [2]byte{0xAE, 0x7D}
	PoisonNumber        = [2]byte{0x0, 0x0}
	magicNumberDuration = 1 * time.Second
	headerDuration      = 2 * time.Second
	readDuration        = 5 * time.Second
	writeDuration       = 5 * time.Second
	keepAlivePeriod     = 10 * time.Second
)

type RaftConnectionStream struct {
	logger         zerolog.Logger
	stream         network.Stream
	messageHandler raftio.MessageHandler
	chunkHandler   raftio.ChunkHandler
}

func (r *RaftConnectionStream) Serve() {
	//for {
	//	msg, err := NewMessageStream(r.stream, nil, r.logger)
	//	if err != nil {
	//		r.logger.Error().Err(err).Msg("failed to readAndHandle message")
	//		return
	//	}
	//
	//	msgType, payload, err := msg.Read()
	//	if err != nil {
	//		r.logger.Error().Err(err).Msg("failed to readAndHandle message")
	//		return
	//	}
	//
	//	switch msgType{
	//	case raftType:
	//		buf := make([]byte, len(payload))
	//		batch := raftpb.MessageBatch{}
	//		if err := batch.Unmarshal(buf); err != nil {
	//			r.logger.Error().Err(err).Msg("failed to unmarshal message")
	//		}
	//		r.messageHandler(batch)
	//	case snapshotType:
	//		buf := make([]byte, len(payload))
	//		chunk := raftpb.Chunk{}
	//		if err := chunk.Unmarshal(buf); err != nil {
	//			r.logger.Error().Err(err).Msg("failed to unmarshal chunk")
	//		}
	//		r.chunkHandler(chunk)
	//	}
	//}
}

func (r *RaftConnectionStream) ReadMessage() error {
	headerBuf := make([]byte, requestHeaderSize)

	headerDeadline := time.Now().Add(headerDuration)
	if err := r.stream.SetReadDeadline(headerDeadline); err != nil {
		r.logger.Error().Err(err).Msg("failed to set readAndHandle deadline for header")
		return err
	}

	if _, err := io.ReadFull(r.stream, headerBuf); err != nil {
		r.logger.Error().Err(err).Msg("failed to readAndHandle header")
		return err
	}

	header := &requestHeader{}
	if err := header.decode(headerBuf); err != nil {
		r.logger.Error().Err(err).Msg("failed to decode header")
		return err
	}

	if header.size == 0 {
		r.logger.Error().Msg("invalid message size")
		return ErrBadMessage
	}

	buf := make([]byte, header.size)
	messageDeadline := time.Now().Add(readDuration)
	if err := r.stream.SetReadDeadline(messageDeadline); err != nil {
		r.logger.Error().Err(err).Msg("failed to set readAndHandle deadline for message")
		return err
	}

	if _, err := io.ReadFull(r.stream, buf); err != nil {
		r.logger.Error().Err(err).Msg("failed to readAndHandle message")
		return err
	}

	if crc32.ChecksumIEEE(buf) != header.crc {
		err := errors.New("invalid message checksum")
		r.logger.Error().Err(err).Msg("invalid message checksum")
		return err
	}

	if header.method == raftType {
		batch := raftpb.MessageBatch{}
		if err := batch.Unmarshal(buf); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal message")
			return err
		}
		r.messageHandler(batch)
	}
	if header.method == snapshotType {
		chunk := raftpb.Chunk{}
		if err := chunk.Unmarshal(buf); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal chunk")
			return err
		}
		r.chunkHandler(chunk)
	}

	return nil
}

func (r *RaftConnectionStream) Poison() error {
	magicNumDuration := time.Now().Add(magicNumberDuration).Add(magicNumberDuration)
	if err := r.stream.SetWriteDeadline(magicNumDuration); err != nil {
		return err
	}
	if _, err := r.stream.Write(PoisonNumber[:]); err != nil {
		return err
	}
	return nil
}

func (r *RaftConnectionStream) Close() {
	if err := r.stream.Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close stream")
	}
}

func (r *RaftConnectionStream) SendMessageBatch(batch raftpb.MessageBatch) error {
	messageBuf := make([]byte, batch.SizeUpperLimit())
	_, err := batch.MarshalTo(messageBuf)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal message batch")
		return err
	}

	header := &requestHeader{
		method: raftType,
		size:   uint64(batch.SizeUpperLimit()),
		crc: crc32.ChecksumIEEE(messageBuf),
	}

	headerBuf := make([]byte, requestHeaderSize)
	headerBuf = header.encode(headerBuf)

	magicNumDeadline := time.Now().Add(magicNumberDuration)
	if err := r.stream.SetWriteDeadline(magicNumDeadline); err != nil {
		r.logger.Error().Err(err).Msg("failed to set write deadline for magic number")
		return err
	}
	if _, err := r.stream.Write(MagicNumber[:]); err != nil {
		r.logger.Error().Err(err).Msg("failed to write magic number")
		return err
	}

	headerDeadline := time.Now().Add(headerDuration)
	if err := r.stream.SetWriteDeadline(headerDeadline); err != nil {
		r.logger.Error().Err(err).Msg("failed to set write deadline for header")
		return err
	}
	if _, err := r.stream.Write(headerBuf); err != nil {
		r.logger.Error().Err(err).Msg("failed to write header")
		return err
	}

	messageDeadline := time.Now().Add(writeDuration)
	if err := r.stream.SetWriteDeadline(messageDeadline); err != nil {
		r.logger.Error().Err(err).Msg("failed to set write deadline for message")
		return err
	}
	if _, err := r.stream.Write(messageBuf); err != nil {
		r.logger.Error().Err(err).Msg("failed to write message batch")
		return err
	}

	return nil
}
