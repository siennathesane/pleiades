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
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/api/v1/raft"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftSystemListener(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping raft system listener tests")
	}
	suite.Run(t, new(RaftSystemListenerTestSuite))
}

type RaftSystemListenerTestSuite struct {
	suite.Suite
	logger zerolog.Logger
	e      *EmbeddedEventStream
	client *EmbeddedEventStreamClient
	defaultTimeout time.Duration
}

func (t *RaftSystemListenerTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 500*time.Millisecond

	opts := &EmbeddedEventStreamOpts{
		Options: &server.Options{
			Host: "localhost",
		},
		timeout: utils.Timeout(4000 * time.Millisecond),
	}

	var err error
	t.e, err = NewEmbeddedEventStream(opts)
	t.Require().NoError(err, "there must not be an error creating the event stream")

	t.e.Start()
}

func (t *RaftSystemListenerTestSuite) SetupTest() {
	var err error
	t.client, err = t.e.GetClient()
	t.Require().NoError(err, "there must not be an error creating the eventStreamClient")
	t.Require().NotNil(t.client, "the eventStreamClient must not be nil")
}

func (t *RaftSystemListenerTestSuite) TearDownTest() {

}

func (t *RaftSystemListenerTestSuite) TestNodeHostShuttingDown() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	nodeShuttingDownMsgs := 0

	_, err := t.client.Subscribe(raftHostTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_NODE_HOST_SHUTTING_DOWN, payload.Action, "the actions must match")
		nodeShuttingDownMsgs += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	listener.NodeHostShuttingDown()
	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, nodeShuttingDownMsgs, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestNodeUnloaded() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	nodeUnloadedMessages := 0

	_, err := t.client.Subscribe(raftNodeTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_NODE_UNLOADED, payload.GetAction(), "the actions must match")
		nodeUnloadedMessages += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}

	listener.NodeUnloaded(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, nodeUnloadedMessages, "we must have nodeUnloadedMessages at least one message")
}

func (t *RaftSystemListenerTestSuite) TestNodeReady() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	nodeReadyMessages := 0

	_, err := t.client.Subscribe(raftNodeTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_NODE_READY, payload.GetAction(), "the actions must match")
		nodeReadyMessages += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}

	listener.NodeReady(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, nodeReadyMessages, "we must have nodeReadyMessages at least one message")
}

func (t *RaftSystemListenerTestSuite) TestMembershipChanged() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	membershipChangedMessages := 0

	_, err := t.client.Subscribe(raftNodeTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_MEMBERSHIP_CHANGED, payload.Action, "the actions must match")
		membershipChangedMessages += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}

	listener.MembershipChanged(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, membershipChangedMessages, "we must have membershipChangedMessages at least one message")
}

func (t *RaftSystemListenerTestSuite) TestConnectionEstablished() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	connectionEstablishedMessages := 0

	_, err := t.client.Subscribe(raftConnectionTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_CONNECTION_ESTABLISHED, payload.GetAction(), "the actions must match")
		connectionEstablishedMessages += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.ConnectionInfo{
		Address: "localhost:1000",
		SnapshotConnection: false,
	}

	listener.ConnectionEstablished(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, connectionEstablishedMessages, "we must have connectionEstablishedMessages at least one message")
}

func (t *RaftSystemListenerTestSuite) TestConnectionFailed() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	connectionFailedMessages := 0

	_, err := t.client.Subscribe(raftConnectionTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_CONNECTION_FAILED, payload.GetAction(), "the actions must match")
		connectionFailedMessages += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.ConnectionInfo{
		Address: "localhost:1000",
		SnapshotConnection: false,
	}

	listener.ConnectionFailed(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, connectionFailedMessages, "we must have connectionFailedMessages at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotStarted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SEND_SNAPSHOT_STARTED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SendSnapshotStarted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotCompleted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SEND_SNAPSHOT_COMPLETED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SendSnapshotCompleted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotAborted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SEND_SNAPSHOT_ABORTED, payload.Action, "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SendSnapshotAborted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotReceived() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SNAPSHOT_RECEIVED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SnapshotReceived(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotRecovered() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SNAPSHOT_RECOVERED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SnapshotRecovered(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotCreated() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SNAPSHOT_CREATED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SnapshotCreated(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotCompacted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftSnapshotTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_SNAPSHOT_COMPACTED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID: 100,
		From: 1,
		Index: 5,
	}

	listener.SnapshotCompacted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestLogCompacted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftLogTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_LOG_COMPACTED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.EntryInfo{
		ClusterID: 10,
		NodeID: 100,
		Index: 5,
	}

	listener.LogCompacted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}

func (t *RaftSystemListenerTestSuite) TestLogDBCompacted() {
	listener := NewRaftListener(t.client, t.logger)
	t.Require().NotNil(listener, "the listener must not be nil")

	received := 0

	_, err := t.client.Subscribe(raftLogTopic, func(msg *nats.Msg) {
		payload := &raft.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
		t.Require().Equal(raft.Event_LOGDB_COMPACTED, payload.GetAction(), "the actions must match")
		received += 1
		msg.AckSync()
	})
	t.Require().NoError(err, "there must not be an error when subscribing")

	testMsg := raftio.EntryInfo{
		ClusterID: 10,
		NodeID: 100,
		Index: 5,
	}

	listener.LogDBCompacted(testMsg)

	utils.Wait(t.defaultTimeout)
	t.Require().Equal(1, received, "we must have received at least one message")
}
