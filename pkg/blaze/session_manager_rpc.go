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

	transportv1 "github.com/mxplusb/pleiades/api/v1"
	"github.com/mxplusb/pleiades/api/v1/database"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

const (
	SessionProtocolVersion protocol.ID = "/pleiades/session/0.0.1"
)

var (
	SessionRPCReadTimeout  time.Duration = 1 * time.Second
	SessionRPCWriteTimeout time.Duration = 1 * time.Second
)

type SessionManagerRPCService struct {
	logger         zerolog.Logger
	sessionManager ISession
	stream         network.Stream
}


func (s *SessionManagerRPCService) handleStream(stream network.Stream) {
	s.stream = stream
	s.readAndHandle()
}

func (s *SessionManagerRPCService) readAndHandle() {
	for {
		// verify the stream state
		if err := VerifyStreamState(s.stream); err != nil {
			s.logger.Error().Err(err).Msg("cannot readAndHandle stream state")
			_ = SendStreamState(s.stream, Invalid, false)
			continue
		}

		// get the header
		if err := s.stream.SetReadDeadline(time.Now().Add(RaftControlRPCReadTimeout)); err != nil {
			s.logger.Error().Err(err).Msg("cannot set read deadline")
			_ = SendStreamState(s.stream, Invalid, false)
		}

		headerBuf := make([]byte, headerSize)
		if _, err := io.ReadFull(s.stream, headerBuf); err != nil {
			s.logger.Error().Err(err).Msg("cannot readAndHandle raft control header")
			_ = SendStreamState(s.stream, Invalid, false)
			continue
		}

		// marshall the header
		header := &transportv1.Header{}
		if err := header.UnmarshalVT(headerBuf); err != nil {
			s.logger.Error().Err(err).Msg("cannot unmarshal header")
			_ = SendStreamState(s.stream, Invalid, false)
			continue
		}

		// prep the message buffer
		msgBuf := make([]byte, header.Size)
		if _, err := io.ReadFull(s.stream, msgBuf); err != nil {
			s.logger.Error().Err(err).Msg("cannot readAndHandle message payload")
			_ = SendStreamState(s.stream, Invalid, false)
		}

		// verify the message is intact
		checked := crc32.ChecksumIEEE(msgBuf)
		if checked != header.Checksum {
			s.logger.Error().Msg("checksums do not match")
			_ = SendStreamState(s.stream, InvalidMessageChecksum, false)
		}

		// unmarshal the payload
		msg := &database.SessionPayload{}
		if err := msg.UnmarshalVT(msgBuf); err != nil {
			s.logger.Error().Err(err).Msg("cannot unmarshal payload")
			_ = SendStreamState(s.stream, Invalid, false)
		}

		switch msg.Method {
		case database.SessionPayload_NEW_SESSION:
			s.newSessionHandler(msg.GetNewSessionRequest())
		}
	}
}

func (s *SessionManagerRPCService) writePayloads(payloadStream <-chan []byte, isStream bool) {
	count := 0
	for {
		if payload, ok := <-payloadStream; ok {
			// send the proper state
			//goland:noinspection GoBoolExpressions
			if count < 1 && isStream {
				if err := SendStreamState(s.stream, StreamStart, true); err != nil {
					s.logger.Error().Err(err).Msg("cannot send stream start state, unrecoverable")
					return
				}
			} else if count > 1 && isStream {
				if err := SendStreamState(s.stream, StreamContinue, true); err != nil {
					s.logger.Error().Err(err).Msg("cannot send stream continue state, unrecoverable")
					return
				}
			} else {
				if err := SendStreamState(s.stream, Valid, true); err != nil {
					s.logger.Error().Err(err).Msg("cannot send stream valid state, unrecoverable")
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
				s.logger.Error().Err(err).Msg("cannot marshal header")
			}

			// set the write deadline
			deadline := time.Now().Add(RaftControlRPCWriteTimeout)
			if err := s.stream.SetWriteDeadline(deadline); err != nil {
				s.logger.Error().Err(err).Msg("cannot set write timeout, unrecoverable")
				return
			}

			// write the header
			if _, err := s.stream.Write(headerBuf); err != nil {
				s.logger.Error().Err(err).Msg("cannot write header to stream, unrecoverable")
				return
			}

			// set the write deadline
			deadline = time.Now().Add(RaftControlRPCWriteTimeout)
			if err := s.stream.SetWriteDeadline(deadline); err != nil {
				s.logger.Error().Err(err).Msg("cannot set write timeout, unrecoverable")
			}

			// write the header
			if _, err := s.stream.Write(payload); err != nil {
				s.logger.Error().Err(err).Msg("cannot write header to stream, unrecoverable")
				return
			}

			count++
		} else if !ok {
			if isStream {
				_ = SendStreamState(s.stream, StreamEnd, false)
			} else {
				_ = SendStreamState(s.stream, Valid, false)
			}
			return
		}
	}
}

func (s *SessionManagerRPCService) newSessionHandler(request *database.NewSessionRequest) {
	clientSession := &client.Session{
		ClusterID: request.GetClusterId(),
	}

	payloadWriter := make(chan []byte)
	defer close(payloadWriter)

	payload := &database.SessionPayload{}

	rs, err := s.sessionManager.ProposeSession(clientSession, 3 * time.Second)
	if err != nil {
		s.logger.Error().Err(err).Msg("error proposing client session")

		payload.Type = &database.SessionPayload_Error{
			Error: &transportv1.DBError{
				Type: transportv1.DBErrorType_SESSION,
				Message: err.Error(),
			},
		}

		go s.writePayloads(payloadWriter, false)
		buf, err := payload.MarshalVT()

		if err != nil {
			s.logger.Error().Err(err).Msg("cannot marshal payload, unrecoverable")
			return
		}

		payloadWriter <- buf
		return
	}

	go s.writePayloads(payloadWriter, true)
	indexState := &database.IndexState{}

	count := 0
	select {
	case response := <-rs.ResultC():
		results := response.GetResult()

		indexState.Results = &database.Result{
			Value: results.Value,
			Data:  results.Data,
		}
		indexState.SnapshotIndex = response.SnapshotIndex()
		indexState.Status = requestStateCodeToResultCode(response)

		buf, err := indexState.MarshalVT()
		if err != nil {
			s.logger.Error().Err(err).Msg("cannot marshal payload, unrecoverable")
			return
		}

		payloadWriter <- buf
		count += 1

		if count == 2 {
			s.logger.Debug().Msg("returned both results")
			return
		}
	}
}