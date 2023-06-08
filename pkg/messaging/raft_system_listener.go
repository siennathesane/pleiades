/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package messaging

import (
	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	RaftHostSubject       string = "system.raftv1.host"
	RaftLogSubject        string = "system.raftv1.log"
	RaftNodeSubject       string = "system.raftv1.node"
	RaftSnapshotSubject   string = "system.raftv1.snapshot"
	RaftConnectionSubject string = "system.raftv1.connection"
	RaftSubject           string = "system.raftv1.raft"
	SystemStreamName      string = "system"
)

var _ raftio.ISystemEventListener = (*RaftSystemListener)(nil)
var _ raftio.IRaftEventListener = (*RaftSystemListener)(nil)

func NewRaftSystemListener(client *EmbeddedMessagingPubSubClient, queueClient *EmbeddedMessagingStreamClient, logger zerolog.Logger) (*RaftSystemListener, error) {
	rl := &RaftSystemListener{
		logger:            logger.With().Str("component", "raft-listener").Logger(),
		eventStreamClient: client,
		queueClient: &EmbeddedMessagingStreamClient{
			JetStreamContext: queueClient,
		},
	}

	return rl, nil
}

type RaftSystemListener struct {
	logger            zerolog.Logger
	eventStreamClient *EmbeddedMessagingPubSubClient
	queueClient       *EmbeddedMessagingStreamClient
}

func (r *RaftSystemListener) LeaderUpdated(info raftio.LeaderInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_RAFT,
		Action:    raftv1.Event_EVENT_LEADER_UPDATED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_LeaderUpdate{
			LeaderUpdate: &raftv1.RaftLeaderInfo{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				Term:      info.Term,
				LeaderId:  info.LeaderID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft host shut down event")
	}

	r.logger.Info().Interface("payload", payload).Msg("leader updated")
}

func (r *RaftSystemListener) NodeHostShuttingDown() {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_HOST,
		Action:    raftv1.Event_EVENT_NODE_HOST_SHUTTING_DOWN,
		Timestamp: timestamppb.Now(),
		Event:     &raftv1.RaftEvent_HostShutdown{},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftHostSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft host shut down event")
	}

	r.logger.Info().Interface("payload", payload).Msg("node shutting down")
}

func (r *RaftSystemListener) NodeUnloaded(info raftio.NodeInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_NODE,
		Action:    raftv1.Event_EVENT_NODE_UNLOADED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Node{
			Node: &raftv1.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node unload event")
	}

	r.logger.Info().Interface("payload", payload).Msg("node unloading")
}

func (r *RaftSystemListener) NodeReady(info raftio.NodeInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_NODE,
		Action:    raftv1.Event_EVENT_NODE_READY,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Node{
			Node: &raftv1.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node ready event")
	}

	r.logger.Info().Interface("payload", payload).Msg("node ready")
}

func (r *RaftSystemListener) MembershipChanged(info raftio.NodeInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_NODE,
		Action:    raftv1.Event_EVENT_MEMBERSHIP_CHANGED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Node{
			Node: &raftv1.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node membership change event")
	}

	r.logger.Info().Interface("payload", payload).Msg("membership changed")
}

func (r *RaftSystemListener) ConnectionEstablished(info raftio.ConnectionInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_CONNECTION,
		Action:    raftv1.Event_EVENT_CONNECTION_ESTABLISHED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Connection{
			Connection: &raftv1.RaftConnectionEvent{
				Address:    info.Address,
				IsSnapshot: info.SnapshotConnection,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftConnectionSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft connection established event")
	}
}

func (r *RaftSystemListener) ConnectionFailed(info raftio.ConnectionInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_CONNECTION,
		Action:    raftv1.Event_EVENT_CONNECTION_FAILED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Connection{
			Connection: &raftv1.RaftConnectionEvent{
				Address:    info.Address,
				IsSnapshot: info.SnapshotConnection,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftConnectionSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft connection failed event")
	}

	r.logger.Info().Interface("payload", payload).Msg("connection failed")
}

func (r *RaftSystemListener) SendSnapshotStarted(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SEND_SNAPSHOT_STARTED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot started event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot started")
}

func (r *RaftSystemListener) SendSnapshotCompleted(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SEND_SNAPSHOT_COMPLETED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot completed event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot completed")
}

func (r *RaftSystemListener) SendSnapshotAborted(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SEND_SNAPSHOT_ABORTED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot aborted event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot aborted")
}

func (r *RaftSystemListener) SnapshotReceived(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SNAPSHOT_RECEIVED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot received event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot received")
}

func (r *RaftSystemListener) SnapshotRecovered(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SNAPSHOT_RECOVERED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot recovered event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot recovered")
}

func (r *RaftSystemListener) SnapshotCreated(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SNAPSHOT_CREATED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot created event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot created")
}

func (r *RaftSystemListener) SnapshotCompacted(info raftio.SnapshotInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_SNAPSHOT,
		Action:    raftv1.Event_EVENT_SNAPSHOT_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Snapshot{
			Snapshot: &raftv1.RaftSnapshotEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				FromIndex: info.From,
				ToIndex:   info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot compacted event")
	}

	r.logger.Info().Interface("payload", payload).Msg("snapshot compacted")
}

func (r *RaftSystemListener) LogCompacted(info raftio.EntryInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_LOG_ENTRY,
		Action:    raftv1.Event_EVENT_LOG_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_LogEntry{
			LogEntry: &raftv1.RaftLogEntryEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				Index:     info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftLogSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft log compacted event")
	}

	r.logger.Info().Interface("payload", payload).Msg("log compacted")
}

func (r *RaftSystemListener) LogDBCompacted(info raftio.EntryInfo) {
	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_LOG_ENTRY,
		Action:    raftv1.Event_EVENT_LOGDB_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_LogEntry{
			LogEntry: &raftv1.RaftLogEntryEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
				Index:     info.Index,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(RaftLogSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft logdb compacted event")
	}

	r.logger.Info().Interface("payload", payload).Msg("logdb compacted")
}
