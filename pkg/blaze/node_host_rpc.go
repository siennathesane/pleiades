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

	"gitlab.com/anthropos-labs/pleiades/api/v1/database"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
)

var (
	_ database.SRPCRaftControlServiceServer = (*NodeRPCServer)(nil)
)

func NewNodeHostRPCServer(logger zerolog.Logger, node INodeHost) *NodeRPCServer {
	return &NodeRPCServer{
		logger: logger,
		node:   node,
	}
}

type NodeRPCServer struct {
	database.SRPCRaftControlServiceUnimplementedServer
	logger zerolog.Logger
	node   INodeHost
}

func (n *NodeRPCServer) GetLeaderID(ctx context.Context, request *database.GetLeaderIDRequest) (*database.GetLeaderIDResponse, error) {
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

func (n *NodeRPCServer) GetID(ctx context.Context, _ *database.IdRequest) (*database.IdResponse, error) {
	ctx = n.logger.WithContext(ctx)
	id := n.node.ID()
	return &database.IdResponse{Id: id}, nil
}

func (n *NodeRPCServer) ReadIndex(request *database.ReadIndexRequest, stream database.SRPCRaftControlService_ReadIndexStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *NodeRPCServer) ReadLocalNode(ctx context.Context, request *database.ReadLocalNodeRequest) (*database.KeyValue, error) {
	ctx = n.logger.WithContext(ctx)

	query, err := request.GetQuery().MarshalVT()
	if err != nil {
		return nil, err
	}

	var rs dragonboat.RequestState
	data, err := n.node.ReadLocalNode(&rs, query)
	if err != nil {
		n.logger.Error().Err(err).Msg("can't read from local node")
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

func (n *NodeRPCServer) AddNode(request *database.ModifyNodeRequest, stream database.SRPCRaftControlService_AddNodeStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *NodeRPCServer) AddObserver(request *database.ModifyNodeRequest, stream database.SRPCRaftControlService_AddObserverStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *NodeRPCServer) AddWitness(request *database.ModifyNodeRequest, stream database.SRPCRaftControlService_AddWitnessStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

// note (sienna): this blocks until the request has been resolved
func (n *NodeRPCServer) RequestCompaction(ctx context.Context, request *database.ModifyNodeRequest) (*database.SysOpState, error) {
	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()
	state, err := n.node.RequestCompaction(clusterId, nodeId)
	if err != nil {
		return nil, err
	}

	select {
	case <- state.ResultC():
		return &database.SysOpState{}, nil
	}
}

func (n *NodeRPCServer) RequestDeleteNode(request *database.ModifyNodeRequest, stream database.SRPCRaftControlService_RequestDeleteNodeStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *NodeRPCServer) RequestLeaderTransfer(ctx context.Context, request *database.ModifyNodeRequest) (*database.RequestLeaderTransferResponse, error) {
	clusterId := request.GetClusterId()
	targetNodeId := request.GetNodeId()
	err := n.node.RequestLeaderTransfer(clusterId, targetNodeId)
	if err != nil {
		return nil, err
	}

	return &database.RequestLeaderTransferResponse{}, nil
}

func (n *NodeRPCServer) RequestSnapshot(request *database.RequestSnapshotRequest, stream database.SRPCRaftControlService_RequestSnapshotStream) error {
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

		if err := errors.Wrap(stream.Send(indexState), "error sending index state"); err != nil {
			return err
		}

		if n.node.NotifyOnCommit() && count == 2 {
			n.logger.Debug().Msg("returned both results")
			return nil
		}
	}
	return nil
}

func (n *NodeRPCServer) Stop(ctx context.Context, request *database.StopRequest) (*database.StopResponse, error) {
	ctx = n.logger.WithContext(ctx)

	n.node.Stop()

	return nil, nil
}

func (n *NodeRPCServer) StopNode(ctx context.Context, request *database.ModifyNodeRequest) (*database.StopNodeResponse, error) {
	ctx = n.logger.WithContext(ctx)

	clusterId := request.GetClusterId()
	nodeId := request.GetNodeId()

	return nil, errors.Wrap(n.node.StopNode(clusterId ,nodeId), "could not stop node")
}

func (n *NodeRPCServer) requestStateCodeToResultCode(result dragonboat.RequestResult) database.IndexState_ResultCode {
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
