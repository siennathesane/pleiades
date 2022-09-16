/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package pubsub

import (
	"github.com/mxplusb/pleiades/api/v1/raft"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	raftHostSubject       string = "system.raft.host"
	raftLogSubject        string = "system.raft.log"
	raftNodeSubject       string = "system.raft.node"
	raftSnapshotSubject   string = "system.raft.snapshot"
	raftConnectionSubject string = "system.raft.connection"
	systemStreamName      string = "SYSTEM"
)

var _ raftio.ISystemEventListener = (*RaftListener)(nil)

func NewRaftListener(client *EmbeddedMessagingPubSubClient, queueClient *EmbeddedMessagingStreamClient, logger zerolog.Logger) (*RaftListener, error) {
	rl := &RaftListener{
		logger:            logger.With().Str("component", "raft-listener").Logger(),
		eventStreamClient: client,
		queueClient: &EmbeddedMessagingStreamClient{
			JetStreamContext: queueClient,
		},
	}

	_, err := rl.queueClient.AddStream(&nats.StreamConfig{
		Name:        systemStreamName,
		Description: "Work queue for Raft system notifications",
		Subjects:    []string{raftHostSubject, raftLogSubject, raftNodeSubject, raftSnapshotSubject, raftConnectionSubject},
		Retention:   nats.WorkQueuePolicy,
		Discard:     nats.DiscardOld,
		Storage:     nats.MemoryStorage,
		MaxMsgs:     1000,
	})
	if err != nil {
		return nil, err
	}

	return rl, nil
}

type RaftListener struct {
	logger            zerolog.Logger
	eventStreamClient *EmbeddedMessagingPubSubClient
	queueClient       *EmbeddedMessagingStreamClient
}

func (r *RaftListener) NodeHostShuttingDown() {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_HOST,
		Action:    raft.Event_NODE_HOST_SHUTTING_DOWN,
		Timestamp: timestamppb.Now(),
		Event:     &raft.RaftEvent_HostShutdown{},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftHostSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft host shut down event")
	}
}

func (r *RaftListener) NodeUnloaded(info raftio.NodeInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_NODE,
		Action:    raft.Event_NODE_UNLOADED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Node{
			Node: &raft.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node unload event")
	}
}

func (r *RaftListener) NodeReady(info raftio.NodeInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_NODE,
		Action:    raft.Event_NODE_READY,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Node{
			Node: &raft.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node ready event")
	}
}

func (r *RaftListener) MembershipChanged(info raftio.NodeInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_NODE,
		Action:    raft.Event_MEMBERSHIP_CHANGED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Node{
			Node: &raft.RaftNodeEvent{
				ShardId:   info.ClusterID,
				ReplicaId: info.NodeID,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftNodeSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft node membership change event")
	}
}

func (r *RaftListener) ConnectionEstablished(info raftio.ConnectionInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_CONNECTION,
		Action:    raft.Event_CONNECTION_ESTABLISHED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Connection{
			Connection: &raft.RaftConnectionEvent{
				Address:    info.Address,
				IsSnapshot: info.SnapshotConnection,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftConnectionSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft connection established event")
	}
}

func (r *RaftListener) ConnectionFailed(info raftio.ConnectionInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_CONNECTION,
		Action:    raft.Event_CONNECTION_FAILED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Connection{
			Connection: &raft.RaftConnectionEvent{
				Address:    info.Address,
				IsSnapshot: info.SnapshotConnection,
			},
		},
	}
	buf, err := payload.MarshalVT()
	if err != nil {
		r.logger.Error().Err(err).Msg("error marshalling payload")
	}

	err = r.eventStreamClient.Publish(raftConnectionSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft connection failed event")
	}
}

func (r *RaftListener) SendSnapshotStarted(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SEND_SNAPSHOT_STARTED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot started event")
	}
}

func (r *RaftListener) SendSnapshotCompleted(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SEND_SNAPSHOT_COMPLETED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot completed event")
	}
}

func (r *RaftListener) SendSnapshotAborted(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SEND_SNAPSHOT_ABORTED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot aborted event")
	}
}

func (r *RaftListener) SnapshotReceived(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SNAPSHOT_RECEIVED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot received event")
	}
}

func (r *RaftListener) SnapshotRecovered(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SNAPSHOT_RECOVERED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot recovered event")
	}
}

func (r *RaftListener) SnapshotCreated(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SNAPSHOT_CREATED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot created event")
	}
}

func (r *RaftListener) SnapshotCompacted(info raftio.SnapshotInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_SNAPSHOT,
		Action:    raft.Event_SNAPSHOT_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_Snapshot{
			Snapshot: &raft.RaftSnapshotEvent{
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

	err = r.eventStreamClient.Publish(raftSnapshotSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft snapshot compacted event")
	}
}

func (r *RaftListener) LogCompacted(info raftio.EntryInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_LOG_ENTRY,
		Action:    raft.Event_LOG_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_LogEntry{
			LogEntry: &raft.RaftLogEntryEvent{
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

	err = r.eventStreamClient.Publish(raftLogSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft log compacted event")
	}
}

func (r *RaftListener) LogDBCompacted(info raftio.EntryInfo) {
	payload := &raft.RaftEvent{
		Typ:       raft.EventType_LOG_ENTRY,
		Action:    raft.Event_LOGDB_COMPACTED,
		Timestamp: timestamppb.Now(),
		Event: &raft.RaftEvent_LogEntry{
			LogEntry: &raft.RaftLogEntryEvent{
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

	err = r.eventStreamClient.Publish(raftLogSubject, buf)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't publish raft logdb compacted event")
	}
}
