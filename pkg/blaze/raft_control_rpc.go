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

	transportv1 "github.com/mxplusb/pleiades/api/v1"
	"github.com/mxplusb/pleiades/api/v1/database"
	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
)

const (
	RaftControlProtocolVersion protocol.ID = "/pleiades/raft-control/0.0.1"
	RaftControlServiceName     string      = "raft-control.pleiades"
)

var (
	RaftControlRPCReadTimeout  time.Duration = 1 * time.Second
	RaftControlRPCWriteTimeout time.Duration = 1 * time.Second
)

func NewRaftControlRPCServer(node INodeHost, host host.Host, logger zerolog.Logger) *RaftControlRPCServer {
	rcrs := &RaftControlRPCServer{
		logger: logger,
		node:   node,
		host:   host,
	}
	rcrs.host.SetStreamHandler(RaftControlProtocolVersion, rcrs.handleStream)

	return rcrs
}

type RaftControlRPCServer struct {
	logger zerolog.Logger
	node   INodeHost
	host   host.Host
}

func (n *RaftControlRPCServer) handleStream(stream network.Stream) {
	if err := stream.Scope().SetService(RaftControlServiceName); err != nil {
		_ = stream.Reset()
	}

	for {
		// verify the stream state
		if err := VerifyStreamState(stream); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle stream state")
			_ = SendStreamState(stream, Invalid, false)
			_ = stream.Reset()
			return
		}

		// get the header
		if err := stream.SetReadDeadline(time.Now().Add(RaftControlRPCReadTimeout)); err != nil {
			n.logger.Error().Err(err).Msg("cannot set read deadline")
			_ = SendStreamState(stream, Invalid, false)
		}

		headerBuf := make([]byte, headerSize)
		if _, err := io.ReadFull(stream, headerBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle raft control header")
			_ = SendStreamState(stream, Invalid, false)
			continue
		}

		// marshall the header
		header := &transportv1.Header{}
		if err := header.UnmarshalVT(headerBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot unmarshal header")
			_ = SendStreamState(stream, Invalid, false)
			continue
		}

		// prep the message buffer
		msgBuf := make([]byte, header.Size)
		if _, err := io.ReadFull(stream, msgBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot readAndHandle message payload")
			_ = SendStreamState(stream, Invalid, false)
		}

		// verify the message is intact
		checked := crc32.ChecksumIEEE(msgBuf)
		if checked != header.Checksum {
			n.logger.Error().Msg("checksums do not match")
			_ = SendStreamState(stream, InvalidMessageChecksum, false)
		}

		// unmarshal the payload
		msg := &database.RaftControlPayload{}
		if err := msg.UnmarshalVT(msgBuf); err != nil {
			n.logger.Error().Err(err).Msg("cannot unmarshal payload")
			_ = SendStreamState(stream, Invalid, false)
		}

		switch msg.Method {
		case database.RaftControlPayload_ADD_NODE:
			n.AddNode(msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_ADD_OBSERVER:
			n.AddObserver(msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_ADD_WITNESS:
			n.AddWitness(msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_GET_ID:
			n.GetID(context.TODO(), nil, stream)
		case database.RaftControlPayload_GET_LEADER_ID:
			n.GetLeaderID(context.TODO(), msg.GetGetLeaderIdRequest(), stream)
		case database.RaftControlPayload_READ_INDEX:
			n.ReadIndex(msg.GetReadIndexRequest(), stream)
		case database.RaftControlPayload_READ_LOCAL_NODE:
			n.ReadLocalNode(context.TODO(), msg.GetReadLocalNodeRequest(), stream)
		case database.RaftControlPayload_REQUEST_COMPACTION:
			n.RequestCompaction(context.TODO(), msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_REQUEST_DELETE_NODE:
			n.RequestDeleteNode(msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_REQUEST_LEADER_TRANSFER:
			n.RequestLeaderTransfer(context.TODO(), msg.GetModifyNodeRequest(), stream)
		case database.RaftControlPayload_REQUEST_SNAPSHOT:
			n.RequestSnapshot(msg.GetRequestSnapshotRequest(), stream)
		case database.RaftControlPayload_STOP:
			n.Stop(context.TODO(), msg.GetStopRequest(), stream)
		case database.RaftControlPayload_STOP_NODE:
			n.StopNode(context.TODO(), msg.GetModifyNodeRequest(), stream)
		}
	}
}

func (n *RaftControlRPCServer) GetLeaderID(ctx context.Context, request *database.GetLeaderIDRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, false, stream)

	ctx = n.logger.WithContext(ctx)
	clusterId := request.GetClusterId()
	leaderId, ok, err := n.node.GetLeaderID(clusterId)
	if !ok {
		n.logger.Error().Err(err).Msg("leader information is not available")
		msg := &transportv1.DBError{Type: transportv1.DBErrorType_RAFT_CONTROL, Message: err.Error()}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	if err != nil && ok {
		n.logger.Error().Err(err).Msg("failed to get leader information")
		msg := &transportv1.DBError{Type: transportv1.DBErrorType_RAFT_CONTROL, Message: "leader information not available"}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	msg := &database.GetLeaderIDResponse{
		LeaderId: leaderId,
	}

	buf, _ := msg.MarshalVT()
	writerChan <- buf
	return
}

func (n *RaftControlRPCServer) GetID(ctx context.Context, _ *database.IdRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, false, stream)

	ctx = n.logger.WithContext(ctx)
	id := n.node.ID()

	msg := &database.IdResponse{Id: id}
	buf, _ := msg.MarshalVT()
	writerChan <- buf
	return
}

func (n *RaftControlRPCServer) ReadIndex(request *database.ReadIndexRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	clusterId := request.GetClusterId()
	timeout := time.Duration(request.GetTimeout())
	rs, err := n.node.ReadIndex(clusterId, timeout)
	if err != nil {
		n.logger.Error().Err(err).Msg("error reading index")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "error reading index").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) ReadLocalNode(ctx context.Context, request *database.ReadLocalNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, false, stream)

	ctx = n.logger.WithContext(ctx)

	query, err := request.GetQuery().MarshalVT()
	if err != nil {
		n.logger.Error().Err(err).Msg("error marshalling query")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "error marshalling query").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	var rs dragonboat.RequestState
	data, err := n.node.ReadLocalNode(&rs, query)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't read from local node")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't read from local node").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	if data == nil {
		err := errors.New("key not found")
		n.logger.Error().Err(err).Msg("incorrect query parameters")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "incorrect query parameters").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	kv := &database.KeyValue{}
	if err := kv.UnmarshalVT(data.([]byte)); err != nil {
		n.logger.Error().Err(err).Msg("can't unmarshal key from local fsm")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't unmarshal key from local fsm").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	buf, _ := kv.MarshalVT()
	writerChan <- buf
}

func (n *RaftControlRPCServer) AddNode(request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddNode(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't add node")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add node").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) AddObserver(request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddObserver(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't add observer")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add observer").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) AddWitness(request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestAddWitness(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't add witness")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add witness").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) RequestCompaction(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	state, err := n.node.RequestCompaction(clusterId, nodeId)
	if err != nil {
		n.logger.Error().Err(err).Msg("compaction can't be completed")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "compaction can't be completed").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for range state.ResultC() {
		msg := &database.SysOpState{}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) RequestDeleteNode(request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	configChange := request.GetConfigChangeIndex()

	rs, err := n.node.RequestDeleteNode(clusterId, nodeId, configChange, timeout)
	if err != nil {
		n.logger.Error().Err(err).Msg("node can't be deleted")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "node can't be deleted").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) RequestLeaderTransfer(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, false, stream)

	clusterId := request.GetClusterId()
	targetNodeId := request.GetNodeId()
	err := n.node.RequestLeaderTransfer(clusterId, targetNodeId)
	if err != nil {
		n.logger.Error().Err(err).Msg("leader can't be transferred")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "leader can't be transferred").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	msg := &database.RequestLeaderTransferResponse{}
	buf, _ := msg.MarshalVT()
	writerChan <- buf
}

func (n *RaftControlRPCServer) RequestSnapshot(request *database.RequestSnapshotRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

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
		n.logger.Error().Err(err).Msg("snapshot can't be created")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "snapshot can't be created").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
		return
	}

	for response := range rs.ResultC() {
		results := response.GetResult()

		indexState := &database.IndexState{
			Results: &database.Result{
				Value: results.Value,
				Data:  results.Data,
			},
			SnapshotIndex: response.SnapshotIndex(),
			Status:        requestStateCodeToResultCode(response),
		}

		buf, _ := indexState.MarshalVT()
		writerChan <- buf
	}
}

func (n *RaftControlRPCServer) Stop(ctx context.Context, request *database.StopRequest, stream network.Stream) {
	ctx = n.logger.WithContext(ctx)

	n.node.Stop()

	_ = SendStreamState(stream, Valid, false)
}

func (n *RaftControlRPCServer) StopNode(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {
	// payload writer
	writerChan := make(chan []byte)
	defer close(writerChan)

	go payloadWriter(writerChan, true, stream)

	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()

	if err := n.node.StopNode(clusterId, nodeId); err != nil {
		n.logger.Error().Err(err).Msg("can't stop node")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't stop node").Error(),
		}
		buf, _ := msg.MarshalVT()
		writerChan <- buf
	}
}

func requestStateCodeToResultCode(result dragonboat.RequestResult) database.IndexState_ResultCode {
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
