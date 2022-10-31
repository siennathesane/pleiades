/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package messaging

import (
	"math/rand"
	"testing"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3/raftio"
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
	logger         zerolog.Logger
	e              *EmbeddedMessaging
	pubSubClient   *EmbeddedMessagingPubSubClient
	queueClient    *EmbeddedMessagingStreamClient
	defaultTimeout time.Duration
}

func (t *RaftSystemListenerTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 500 * time.Millisecond

	var err error
	t.e, err = NewEmbeddedMessagingWithDefaults(t.logger)
	t.Require().NoError(err, "there must not be an error creating the event stream")

	t.e.Start()
	t.Require().NoError(err, "there must not be an error when creating the system stream")
}

func (t *RaftSystemListenerTestSuite) SetupTest() {
	var err error
	t.pubSubClient, err = t.e.GetPubSubClient()
	t.Require().NoError(err, "there must not be an error creating the eventStreamClient")
	t.Require().NotNil(t.pubSubClient, "the eventStreamClient must not be nil")

	t.queueClient, err = t.e.GetStreamClient()
	t.Require().NoError(err, "there must not be an error when getting a stream pubSubClient")
	t.Require().NotNil(t.queueClient, "the queue pubSubClient must not be nil")
}

func (t *RaftSystemListenerTestSuite) TestNewNewRaftSystemListener() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")
}

func (t *RaftSystemListenerTestSuite) TestLeaderUpdated() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.LeaderInfo{
		ClusterID: rand.Uint64(),
		NodeID:    rand.Uint64(),
		Term:      rand.Uint64(),
		LeaderID:  rand.Uint64(),
	}
	listener.LeaderUpdated(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(1000 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_LEADER_UPDATED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestNodeHostShuttingDown() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftHostSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	listener.NodeHostShuttingDown()

	msgs, err := sub.NextMsg(utils.Timeout(1000 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_NODE_HOST_SHUTTING_DOWN, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestNodeUnloaded() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftNodeSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}
	listener.NodeUnloaded(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_NODE_UNLOADED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestNodeReady() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftNodeSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}
	listener.NodeReady(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_NODE_READY, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestMembershipChanged() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftNodeSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.NodeInfo{
		ClusterID: 10,
		NodeID:    100,
	}
	listener.MembershipChanged(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_MEMBERSHIP_CHANGED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestConnectionEstablished() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftConnectionSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.ConnectionInfo{
		Address:            "localhost:1000",
		SnapshotConnection: false,
	}
	listener.ConnectionEstablished(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_CONNECTION_ESTABLISHED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestConnectionFailed() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftConnectionSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.ConnectionInfo{
		Address:            "localhost:1000",
		SnapshotConnection: false,
	}
	listener.ConnectionFailed(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_CONNECTION_FAILED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotStarted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SendSnapshotStarted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SEND_SNAPSHOT_STARTED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotCompleted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SendSnapshotCompleted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SEND_SNAPSHOT_COMPLETED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSendSnapshotAborted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SendSnapshotAborted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SEND_SNAPSHOT_ABORTED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotReceived() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SnapshotReceived(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SNAPSHOT_RECEIVED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotRecovered() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SnapshotRecovered(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SNAPSHOT_RECOVERED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotCreated() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SnapshotCreated(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SNAPSHOT_CREATED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestSnapshotCompacted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftSnapshotSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.SnapshotInfo{
		ClusterID: 10,
		NodeID:    100,
		From:      1,
		Index:     5,
	}
	listener.SnapshotCompacted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_SNAPSHOT_COMPACTED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestLogCompacted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftLogSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.EntryInfo{
		ClusterID: 10,
		NodeID:    100,
		Index:     5,
	}
	listener.LogCompacted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_LOG_COMPACTED, payload.GetAction(), "the actions must match")
}

func (t *RaftSystemListenerTestSuite) TestLogDBCompacted() {
	listener, err := NewRaftSystemListener(t.pubSubClient, t.queueClient, t.logger)
	t.Require().NoError(err, "there must not be an error when creating the listener")
	t.Require().NotNil(listener, "the listener must not be nil")

	sub, err := t.pubSubClient.SubscribeSync(RaftLogSubject)
	t.Require().NoError(err, "there must not be an error when subscribing")
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	testMsg := raftio.EntryInfo{
		ClusterID: 10,
		NodeID:    100,
		Index:     5,
	}
	listener.LogDBCompacted(testMsg)

	msgs, err := sub.NextMsg(utils.Timeout(100 * time.Millisecond))
	t.Require().NoError(err, "there must not be an error when fetching the messages")

	payload := &raftv1.RaftEvent{}
	err = payload.UnmarshalVT(msgs.Data)
	t.Require().NoError(err, "there must not be an error when unmarshalling the payload")
	t.Require().Equal(raftv1.Event_EVENT_LOGDB_COMPACTED, payload.GetAction(), "the actions must match")
}
