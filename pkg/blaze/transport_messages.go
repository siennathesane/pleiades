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
	"bytes"
	"hash/crc32"
	"io"
	"time"

	transportv1 "github.com/mxplusb/pleiades/api/v1"
	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/rs/zerolog"
)

func NewMessageStream(stream network.Stream, msg []byte, logger zerolog.Logger) (*MessageStream, error) {
	if len(msg) == 0 || msg == nil {
		err := errors.New("no message to send")
		logger.Error().Err(err).Msg("empty message")
		return nil, err
	}

	return &MessageStream{
		stream: stream,
		body:   msg,
	}, nil
}

type MessageStream struct {
	stream network.Stream
	logger zerolog.Logger
	header requestHeader
	body   []byte
}

func (m *MessageStream) VerifyMagicNumber() error {
	readDeadline := time.Now().Add(magicNumberDuration)
	if err := m.stream.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	magicLen := make([]byte, len(MagicNumber))
	if _, err := io.ReadFull(m.stream, magicLen); err != nil {
		return err
	}

	if !bytes.Equal(magicLen, PoisonNumber[:]) {
		return errPoisonReceived
	}

	if !bytes.Equal(magicLen, MagicNumber[:]) {
		return ErrBadMessage
	}

	return nil
}

func (m *MessageStream) Send(t uint16) error {
	if m.body == nil {
		err := errors.New("no body to send")
		m.logger.Error().Err(err).Msg("body is nil")
		return err
	}

	header := &requestHeader{
		method: t,
		size:   uint64(len(m.body)),
		crc:    crc32.ChecksumIEEE(m.body),
	}

	headerBuf := make([]byte, requestHeaderSize)
	headerBuf = header.encode(headerBuf)

	magicNumDeadline := time.Now().Add(magicNumberDuration)
	if err := m.stream.SetWriteDeadline(magicNumDeadline); err != nil {
		m.logger.Error().Err(err).Msg("failed to set write deadline for magic number")
		return err
	}
	if _, err := m.stream.Write(MagicNumber[:]); err != nil {
		m.logger.Error().Err(err).Msg("failed to write magic number")
		return err
	}

	headerDeadline := time.Now().Add(headerDuration)
	if err := m.stream.SetWriteDeadline(headerDeadline); err != nil {
		m.logger.Error().Err(err).Msg("failed to set write deadline for header")
		return err
	}
	if _, err := m.stream.Write(headerBuf); err != nil {
		m.logger.Error().Err(err).Msg("failed to write header")
		return err
	}

	messageDeadline := time.Now().Add(writeDuration)
	if err := m.stream.SetWriteDeadline(messageDeadline); err != nil {
		m.logger.Error().Err(err).Msg("failed to set write deadline for message")
		return err
	}
	if _, err := m.stream.Write(m.body); err != nil {
		m.logger.Error().Err(err).Msg("failed to write message batch")
		return err
	}

	return nil
}

func (m *MessageStream) Read() (uint16, []byte, error) {
	if err := m.VerifyMagicNumber(); err != nil {
		if err == errPoisonReceived {
			m.logger.Error().Err(err).Msg("poison received")
		}
		if err == ErrBadMessage {
			m.logger.Error().Err(err).Msg("bad message")
			return 0, nil, err
		}
	}

	headerBuf := make([]byte, requestHeaderSize)

	headerDeadline := time.Now().Add(headerDuration)
	if err := m.stream.SetReadDeadline(headerDeadline); err != nil {
		m.logger.Error().Err(err).Msg("failed to set readAndHandle deadline for header")
		return 0, nil, err
	}

	if _, err := io.ReadFull(m.stream, headerBuf); err != nil {
		m.logger.Error().Err(err).Msg("failed to readAndHandle header")
		return 0, nil, err
	}

	header := &requestHeader{}
	if err := header.decode(headerBuf); err != nil {
		m.logger.Error().Err(err).Msg("failed to decode header")
		return 0, nil, err
	}

	if header.size == 0 {
		m.logger.Error().Msg("invalid message size")
		return 0, nil, ErrBadMessage
	}

	buf := make([]byte, header.size)
	messageDeadline := time.Now().Add(readDuration)
	if err := m.stream.SetReadDeadline(messageDeadline); err != nil {
		m.logger.Error().Err(err).Msg("failed to set readAndHandle deadline for message")
		return 0, nil, err
	}

	if _, err := io.ReadFull(m.stream, buf); err != nil {
		m.logger.Error().Err(err).Msg("failed to readAndHandle message")
		return 0, nil, err
	}

	if crc32.ChecksumIEEE(buf) != header.crc {
		err := errors.New("invalid message checksum")
		m.logger.Error().Err(err).Msg("invalid message checksum")
		return 0, nil, err
	}

	return header.method, buf, nil
}

func payloadWriter(payloadStream <-chan []byte, isStream bool, stream network.Stream) {
	count := 0
	for {
		payload, open := <-payloadStream
		if !open {
			if isStream {
				if count > 1 {
					_ = SendStreamState(stream, StreamNoLongerValid, false)
				}
				_ = SendStreamState(stream, StreamEnd, false)
			} else {
				_ = SendStreamState(stream, Valid, false)
			}
			return
		}

		// send the proper state
		if count < 1 && isStream {
			if err := SendStreamState(stream, StreamStart, true); err != nil {
				_ = stream.Reset()
				return
			}
		} else if count > 1 && isStream {
			if err := SendStreamState(stream, StreamContinue, true); err != nil {
				_ = stream.Reset()
				return
			}
		} else {
			if err := SendStreamState(stream, Valid, true); err != nil {
				_ = stream.Reset()
				return
			}
		}

		// set the header
		header := transportv1.Header{
			Size:     uint32(len(payload)),
			Checksum: crc32.ChecksumIEEE(payload),
		}
		headerBuf, err := header.MarshalVT()
		if err != nil {
			if err := SendStreamState(stream, Valid, true); err != nil {
				_ = stream.Reset()
				return
			}
		}

		// set the write deadline
		deadline := time.Now().Add(RaftControlRPCWriteTimeout)
		if err := stream.SetWriteDeadline(deadline); err != nil {
			_ = stream.Reset()
		}

		// write the header
		if _, err := stream.Write(headerBuf); err != nil {
			_ = stream.Reset()
		}

		// set the write deadline
		deadline = time.Now().Add(RaftControlRPCWriteTimeout)
		if err := stream.SetWriteDeadline(deadline); err != nil {
			_ = stream.Reset()
		}

		// write the payload
		if _, err := stream.Write(payload); err != nil {
			_ = stream.Reset()
			return
		}

		count++
	}
}
