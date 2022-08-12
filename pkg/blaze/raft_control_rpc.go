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
	"hash/crc32"
	"io"
	"time"

	transportv1 "github.com/mxplusb/pleiades/pkg/api/v1"
	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
)

const (
	RaftControlProtocolVersion protocol.ID = "pleiades/raft-control/0.0.1"
)

var (
	RaftControlRPCReadTimeout  time.Duration = 1 * time.Second
	RaftControlRPCWriteTimeout time.Duration = 1 * time.Second
)

func NewRaftControlRPCServer(node INodeHost, logger zerolog.Logger) *RaftControlRPCServer {
	return &RaftControlRPCServer{
		logger: logger,
		node:   node,
	}
}

type RaftControlRPCServer struct {
	logger zerolog.Logger
	node   INodeHost
	stream network.Stream
}

func (n *RaftControlRPCServer) handleStream(stream network.Stream) {
	n.stream = stream
	n.readAndHandle()
}

func (n *RaftControlRPCServer) readAndHandle() {
	for {
		// verify the stream state
		if err := VerifyStreamState(n.stream); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle stream state")
			_ = SendStreamState(n.stream, Invalid, false)
			continue
		}

		// get the header
		if err := n.stream.SetReadDeadline(time.Now().Add(RaftControlRPCReadTimeout)); err != nil {
			n.logger.Error().Err(err).Msg("cannot set read deadline")
			_ = SendStreamState(n.stream, Invalid, false)
		}

		headerBuf := make([]byte, headerSize)
		if _, err := io.ReadFull(n.stream, headerBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle raft control header")
			_ = SendStreamState(n.stream, Invalid, false)
			continue
		}

		// marshall the header
		header := &transportv1.Header{}
		if err := header.UnmarshalVT(headerBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot unmarshal header")
			_ = SendStreamState(n.stream, Invalid, false)
			continue
		}

		// prep the message buffer
		msgBuf := make([]byte, header.Size)
		if _, err := io.ReadFull(n.stream, msgBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle message payload")
			_ = SendStreamState(n.stream, Invalid, false)
		}

		// verify the message is intact
		checked := crc32.ChecksumIEEE(msgBuf)
		if checked != header.Checksum {
			n.logger.Error().Msg("checksums do not match")
			_ = SendStreamState(n.stream, InvalidMessageChecksum, false)
		}

		// unmarshal the payload
		msg := &database.RaftControlPayload{}
		if err := msg.UnmarshalVT(msgBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot unmarshal payload")
			_ = SendStreamState(n.stream, Invalid, false)
		}

		switch msg.Method {
		case database.RaftControlPayload_ADD_NODE:
			n.addNodeHandler(msg.GetModifyNodeRequest())
		}
	}
}

func (n *RaftControlRPCServer) writePayloads(payloadStream <-chan []byte, isStream bool) {
	count := 0
	for {
		if payload, ok := <-payloadStream; ok {
			// send the proper state
			//goland:noinspection GoBoolExpressions
			if count < 1 && isStream {
				if err := SendStreamState(n.stream, StreamStart, true); err != nil {
					n.logger.Error().Err(err).Msg("cannot send stream start state, unrecoverable")
					return
				}
			} else if count > 1 && isStream {
				if err := SendStreamState(n.stream, StreamContinue, true); err != nil {
					n.logger.Error().Err(err).Msg("cannot send stream continue state, unrecoverable")
					return
				}
			} else {
				if err := SendStreamState(n.stream, Valid, true); err != nil {
					n.logger.Error().Err(err).Msg("cannot send stream valid state, unrecoverable")
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
				n.logger.Error().Err(err).Msg("cannot marshal header")
			}

			// set the write deadline
			deadline := time.Now().Add(RaftControlRPCWriteTimeout)
			if err := n.stream.SetWriteDeadline(deadline); err != nil {
				n.logger.Error().Err(err).Msg("cannot set write timeout, unrecoverable")
			}

			// write the header
			if _, err := n.stream.Write(headerBuf); err != nil {
				n.logger.Error().Err(err).Msg("cannot write header to stream, unrecoverable")
				return
			}

			// set the write deadline
			deadline = time.Now().Add(RaftControlRPCWriteTimeout)
			if err := n.stream.SetWriteDeadline(deadline); err != nil {
				n.logger.Error().Err(err).Msg("cannot set write timeout, unrecoverable")
			}

			// write the header
			if _, err := n.stream.Write(payload); err != nil {
				n.logger.Error().Err(err).Msg("cannot write header to stream, unrecoverable")
				return
			}

			count++
		} else if !ok {
			if isStream {
				_ = SendStreamState(n.stream, StreamEnd, false)
			} else {
				_ = SendStreamState(n.stream, Valid, false)
			}
			return
		}
	}
}

func (n *RaftControlRPCServer) GetLeaderID(ctx context.Context, request *database.GetLeaderIDRequest) (*database.GetLeaderIDResponse, error) {
	ctx = n.logger.WithContext(ctx)
	clusterId := request.GetClusterId()
	leaderId, ok, err := n.node.GetLeaderID(clusterId)
	if !ok {
		n.logger.Error().Err(err).Msg("leader information is not available")
		return nil, err
	}
	if err != nil && ok {
		n.logger.Error().Err(err).Msg("failed to get leader information")
		return nil, err
	}
	return &database.GetLeaderIDResponse{
		LeaderId: leaderId,
	}, err
}

func (n *RaftControlRPCServer) GetID(ctx context.Context, _ *database.IdRequest) (*database.IdResponse, error) {
	ctx = n.logger.WithContext(ctx)
	id := n.node.ID()
	return &database.IdResponse{Id: id}, nil
}

func (n *RaftControlRPCServer) ReadIndex(request *database.ReadIndexRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	timeout := time.Duration(request.GetTimeout())
	rs, err := n.node.ReadIndex(clusterId, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *RaftControlRPCServer) ReadLocalNode(ctx context.Context, request *database.ReadLocalNodeRequest) (*database.KeyValue, error) {
	ctx = n.logger.WithContext(ctx)

	query, err := request.GetQuery().MarshalVT()
	if err != nil {
		return nil, err
	}

	var rs dragonboat.RequestState
	data, err := n.node.ReadLocalNode(&rs, query)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't readAndHandle from local node")
		return nil, err
	}

	if data == nil {
		err := errors.New("key not found")
		n.logger.Error().Err(err).Msg("incorrect query parameters")
		return nil, err
	}

	kv := &database.KeyValue{}
	if err := kv.UnmarshalVT(data.([]byte)); err != nil {
		n.logger.Error().Err(err).Msg("can't unmarshal key from local fsm")
		return nil, err
	}

	return kv, nil
}

func (n *RaftControlRPCServer) addNodeHandler(request *database.ModifyNodeRequest) {
	resultsChan := make(chan *database.IndexState, 2)
	writerChan := make(chan []byte)

	go func(idxState chan *database.IndexState, writerPayloads chan<- []byte) {
		count := 0
		for count < 3 {
			if idx, ok := <- idxState; ok {
				msg := &database.RaftControlPayload{
					Types: &database.RaftControlPayload_IndexState{IndexState: idx},
				}
				buf, _ := msg.MarshalVT()
				writerChan <- buf
				count++
			} else if !ok {
				return
			}
			return
		}
	}(resultsChan, writerChan)
	go n.writePayloads(writerChan, true)

	if err := n.AddNode(request, resultsChan); err != nil {
		n.logger.Error().Err(err).Msg("")
	}
}

func (n *RaftControlRPCServer) AddNode(request *database.ModifyNodeRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddNode(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *RaftControlRPCServer) AddObserver(request *database.ModifyNodeRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddObserver(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *RaftControlRPCServer) AddWitness(request *database.ModifyNodeRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddWitness(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

// note (sienna): this blocks until the request has been resolved
func (n *RaftControlRPCServer) RequestCompaction(ctx context.Context, request *database.ModifyNodeRequest) (*database.SysOpState, error) {
	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	state, err := n.node.RequestCompaction(clusterId, nodeId)
	if err != nil {
		return nil, err
	}

	select {
	case <-state.ResultC():
		return &database.SysOpState{}, nil
	}
}

func (n *RaftControlRPCServer) RequestDeleteNode(request *database.ModifyNodeRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestDeleteNode(clusterId, nodeId, configChange, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *RaftControlRPCServer) RequestLeaderTransfer(ctx context.Context, request *database.ModifyNodeRequest) (*database.RequestLeaderTransferResponse, error) {
	clusterId := request.GetClusterId()
	targetNodeId := request.GetNodeId()
	err := n.node.RequestLeaderTransfer(clusterId, targetNodeId)
	if err != nil {
		return nil, err
	}

	return &database.RequestLeaderTransferResponse{}, nil
}

func (n *RaftControlRPCServer) RequestSnapshot(request *database.RequestSnapshotRequest, stream chan<- *database.IndexState) error {
	clusterId := request.GetClusterId()
	snapOpts := request.GetOptions()
	timeout := time.Duration(request.GetTimeout())

	opts := dragonboat.SnapshotOption{
		CompactionOverhead:         snapOpts.CompactionOverhead,
		ExportPath:                 snapOpts.ExportPath,
		Exported:                   snapOpts.Exported,
		OverrideCompactionOverhead: snapOpts.OverrideCompactionOverhead,
	}

	rs, err := n.node.RequestSnapshot(clusterId, opts, timeout)
	if err != nil {
		return err
	}

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
		indexState.Status = n.requestStateCodeToResultCode(response)

		count += 1

		stream <- indexState

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *RaftControlRPCServer) Stop(ctx context.Context, request *database.StopRequest) (*database.StopResponse, error) {
	ctx = n.logger.WithContext(ctx)

	n.node.Stop()

	return nil, nil
}

func (n *RaftControlRPCServer) StopNode(ctx context.Context, request *database.ModifyNodeRequest) (*database.StopNodeResponse, error) {
	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()

	return nil, errors.Wrap(n.node.StopNode(clusterId, nodeId), "could not stop node")
}

func (n *RaftControlRPCServer) requestStateCodeToResultCode(result dragonboat.RequestResult) database.IndexState_ResultCode {
	switch {
	case result.Aborted():
		return database.IndexState_Aborted
	case result.Committed():
		return database.IndexState_Committed
	case result.Dropped():
		return database.IndexState_Dropped
	case result.Rejected():
		return database.IndexState_Rejected
	case result.Terminated():
		return database.IndexState_Terminated
	case result.Timeout():
		return database.IndexState_Timeout
	default:
		return database.IndexState_Completed
	}
}
