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

	AddNode        MethodByte = 0x00
	AddObserver    MethodByte = 0x01
	AddWitness     MethodByte = 0x02
	GetId          MethodByte = 0x03
	GetLeaderId    MethodByte = 0x04
	Compact        MethodByte = 0x07
	DeleteNode     MethodByte = 0x08
	LeaderTransfer MethodByte = 0x09
	Snapshot       MethodByte = 0x10
	Stop           MethodByte = 0x11
	StopNode       MethodByte = 0x12
	Error          MethodByte = 0xfe
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

func (r *RaftControlRPCServer) handleStream(stream network.Stream) {
	if err := stream.Scope().SetService(RaftControlServiceName); err != nil {
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

		switch frame.GetMethod() {
		case AddNode:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshalling message")
			}
			r.AddNode(msg, stream)
		case AddObserver:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshalling message")
			}
			r.AddObserver(msg, stream)
		case AddWitness:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshalling message")
			}
			r.AddWitness(msg, stream)
		case GetId:
			r.GetID(context.TODO(), nil, stream)
		case GetLeaderId:
			msg, err := unmarshal[*database.GetLeaderIDRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.GetLeaderID(context.TODO(), msg, stream)
		case Compact:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.Compact(context.TODO(), msg, stream)
		case DeleteNode:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.DeleteNode(msg, stream)
		case LeaderTransfer:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.LeaderTransfer(context.TODO(), msg, stream)
		case Snapshot:
			msg, err := unmarshal[*database.RequestSnapshotRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.Snapshot(msg, stream)
		case Stop:
			msg, err := unmarshal[*database.StopRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.Stop(context.TODO(), msg, stream)
		case StopNode:
			msg, err := unmarshal[*database.ModifyNodeRequest](msgBuf)
			if err != nil {
				// todo (sienna): add error handling
				r.logger.Error().Err(err).Msg("error unmarshaling message")
			}
			r.StopNode(context.TODO(), msg, stream)
		}
	}
}

func (r *RaftControlRPCServer) handleRequestState(rs *dragonboat.RequestState, method MethodByte, stream network.Stream) {
	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(method)
	first := true
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
		if first {
			responseFrame = responseFrame.WithPayload(buf).WithState(StreamStartByte).WithMethod(method)
		} else {
			responseFrame = responseFrame.WithPayload(buf).WithState(StreamContinueByte).WithMethod(method)
		}
		sendFrame(responseFrame, r.logger, stream)

		first = false
	}

	responseFrame = responseFrame.WithState(StreamEndByte).WithMethod(AddNode)
	sendFrame(responseFrame, r.logger, stream)
}

func (r *RaftControlRPCServer) AddNode(request *database.ModifyNodeRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(AddNode)

	rs, err := r.node.RequestAddNode(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add node")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add node").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	r.handleRequestState(rs, AddNode, stream)
}

func (r *RaftControlRPCServer) AddObserver(request *database.ModifyNodeRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(AddObserver)

	rs, err := r.node.RequestAddObserver(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add observer")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add observer").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	r.handleRequestState(rs, AddObserver, stream)
}

func (r *RaftControlRPCServer) AddWitness(request *database.ModifyNodeRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	target := request.GetTarget()
	configChange := request.GetConfigChangeIndex()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(AddWitness)

	rs, err := r.node.RequestAddWitness(clusterId, nodeId, target, configChange, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't add witness")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't add witness").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	r.handleRequestState(rs, AddWitness, stream)
}

func (r *RaftControlRPCServer) GetID(ctx context.Context, _ *database.IdRequest, stream network.Stream) {
	ctx = r.logger.WithContext(ctx)
	id := r.node.ID()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(GetLeaderId)

	msg := &database.IdResponse{Id: id}
	buf, _ := msg.MarshalVT()
	responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(GetId)
	sendFrame(responseFrame, r.logger, stream)
	return
}

func (r *RaftControlRPCServer) GetLeaderID(ctx context.Context, request *database.GetLeaderIDRequest, stream network.Stream) {
	ctx = r.logger.WithContext(ctx)
	clusterId := request.GetClusterId()
	leaderId, ok, err := r.node.GetLeaderID(clusterId)

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(GetLeaderId)

	if !ok {
		r.logger.Error().Err(err).Msg("leader information is not available")
		msg := &transportv1.DBError{Type: transportv1.DBErrorType_RAFT_CONTROL, Message: err.Error()}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	if err != nil && ok {
		r.logger.Error().Err(err).Msg("failed to get leader information")
		msg := &transportv1.DBError{Type: transportv1.DBErrorType_RAFT_CONTROL, Message: "leader information not available"}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	msg := &database.GetLeaderIDResponse{
		LeaderId: leaderId,
	}

	buf, _ := msg.MarshalVT()
	responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(GetLeaderId)
	sendFrame(responseFrame, r.logger, stream)
	return
}

func (r *RaftControlRPCServer) Compact(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {
	ctx = r.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(Compact)

	state, err := r.node.RequestCompaction(clusterId, nodeId)
	if err != nil {
		r.logger.Error().Err(err).Msg("compaction can't be completed")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "compaction can't be completed").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	first := true
	for range state.ResultC() {
		resp := &database.SysOpState{}

		buf, _ := resp.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithMethod(Compact)
		if first {
			responseFrame = responseFrame.WithState(StreamStartByte)
		} else {
			responseFrame = responseFrame.WithState(StreamContinueByte)
		}
		sendFrame(responseFrame, r.logger, stream)

		first = false
	}

	responseFrame = responseFrame.WithState(StreamEndByte).WithMethod(Compact)
	sendFrame(responseFrame, r.logger, stream)
}

func (r *RaftControlRPCServer) DeleteNode(request *database.ModifyNodeRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	timeout := time.Duration(request.GetTimeout())
	configChange := request.GetConfigChangeIndex()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(DeleteNode)

	rs, err := r.node.RequestDeleteNode(clusterId, nodeId, configChange, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("node can't be deleted")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "node can't be deleted").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	r.handleRequestState(rs, DeleteNode, stream)
}

func (r *RaftControlRPCServer) LeaderTransfer(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	targetNodeId := request.GetNodeId()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(LeaderTransfer)

	err := r.node.RequestLeaderTransfer(clusterId, targetNodeId)
	if err != nil {
		r.logger.Error().Err(err).Msg("leader can't be transferred")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "leader can't be transferred").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	msg := &database.RequestLeaderTransferResponse{}
	buf, _ := msg.MarshalVT()
	responseFrame = responseFrame.WithPayload(buf)
	sendFrame(responseFrame, r.logger, stream)
}

func (r *RaftControlRPCServer) Snapshot(request *database.RequestSnapshotRequest, stream network.Stream) {

	clusterId := request.GetClusterId()
	snapOpts := request.GetOptions()
	timeout := time.Duration(request.GetTimeout())

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(Snapshot)

	opts := dragonboat.SnapshotOption{
		CompactionOverhead:         snapOpts.CompactionOverhead,
		ExportPath:                 snapOpts.ExportPath,
		Exported:                   snapOpts.Exported,
		OverrideCompactionOverhead: snapOpts.OverrideCompactionOverhead,
	}

	rs, err := r.node.RequestSnapshot(clusterId, opts, timeout)
	if err != nil {
		r.logger.Error().Err(err).Msg("snapshot can't be created")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "snapshot can't be created").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame = responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	r.handleRequestState(rs, Snapshot, stream)
}

func (r *RaftControlRPCServer) Stop(ctx context.Context, _ *database.StopRequest, stream network.Stream) {
	ctx = r.logger.WithContext(ctx)

	r.node.Stop()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(Stop)
	sendFrame(responseFrame, r.logger, stream)
}

func (r *RaftControlRPCServer) StopNode(ctx context.Context, request *database.ModifyNodeRequest, stream network.Stream) {
	ctx = r.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()

	responseFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(StopNode)

	if err := r.node.StopNode(clusterId, nodeId); err != nil {
		r.logger.Error().Err(err).Msg("can't stop node")
		msg := &transportv1.DBError{
			Type:    transportv1.DBErrorType_RAFT_CONTROL,
			Message: errors.Wrap(err, "can't stop node").Error(),
		}
		buf, _ := msg.MarshalVT()
		responseFrame := responseFrame.WithPayload(buf).WithState(ValidByte).WithMethod(Error)
		sendFrame(responseFrame, r.logger, stream)
		return
	}

	responseFrame = responseFrame.WithMethod(StopNode)
	sendFrame(responseFrame, r.logger, stream)
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
