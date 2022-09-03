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
	"time"

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
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
	sessionManager ITransactionManager
	stream         network.Stream
}


func (s *SessionManagerRPCService) handleStream(stream network.Stream) {
	s.stream = stream
	s.readAndHandle()
}

func (s *SessionManagerRPCService) readAndHandle() {

}

func (s *SessionManagerRPCService) writePayloads(payloadStream <-chan []byte, isStream bool) {

}

func (s *SessionManagerRPCService) newSessionHandler(request *database.NewSessionRequest) {
	//clientSession := &client.Session{
	//	ClusterID: request.GetClusterId(),
	//}
	//
	//payloadWriter := make(chan []byte)
	//defer close(payloadWriter)
	//
	//payload := &database.SessionPayload{}
	//
	//rs, err := s.sessionManager.ProposeSession(clientSession, 3 * time.Second)
	//if err != nil {
	//	s.logger.Error().Err(err).Msg("error proposing client session")
	//
	//	payload.Type = &database.SessionPayload_Error{
	//		Error: &transportv1.DBError{
	//			Type: transportv1.DBErrorType_SESSION,
	//			Message: err.Error(),
	//		},
	//	}
	//
	//	go s.writePayloads(payloadWriter, false)
	//	buf, err := payload.MarshalVT()
	//
	//	if err != nil {
	//		s.logger.Error().Err(err).Msg("cannot marshal payload, unrecoverable")
	//		return
	//	}
	//
	//	payloadWriter <- buf
	//	return
	//}
	//
	//go s.writePayloads(payloadWriter, true)
	//indexState := &database.IndexState{}
	//
	//count := 0
	//select {
	//case response := <-rs.ResultC():
	//	results := response.GetResult()
	//
	//	indexState.Results = &database.Result{
	//		Value: results.Value,
	//		Data:  results.Data,
	//	}
	//	indexState.SnapshotIndex = response.SnapshotIndex()
	//	indexState.Status = requestStateCodeToResultCode(response)
	//
	//	buf, err := indexState.MarshalVT()
	//	if err != nil {
	//		s.logger.Error().Err(err).Msg("cannot marshal payload, unrecoverable")
	//		return
	//	}
	//
	//	payloadWriter <- buf
	//	count += 1
	//
	//	if count == 2 {
	//		s.logger.Debug().Msg("returned both results")
	//		return
	//	}
	//}
}