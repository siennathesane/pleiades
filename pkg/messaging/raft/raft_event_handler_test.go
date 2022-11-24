/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package raft

import (
	"testing"
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/messaging/clients"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRaftEventHandler(t *testing.T) {
	suite.Run(t, new(RaftEventHandlerTestSuite))
}

type RaftEventHandlerTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	e              *messaging.EmbeddedMessaging
	pubSubClient   *clients.EmbeddedMessagingPubSubClient
	queueClient    *clients.EmbeddedMessagingStreamClient
	defaultTimeout time.Duration
}

func (t *RaftEventHandlerTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 500 * time.Millisecond

	var err error
	t.e, err = messaging.NewEmbeddedMessagingWithDefaults(t.logger)
	t.Require().NoError(err, "there must not be an error creating the event stream")

	t.e.Start()
}

func (t *RaftEventHandlerTestSuite) SetupTest() {
	var err error
	t.pubSubClient, err = t.e.GetPubSubClient()
	t.Require().NoError(err, "there must not be an error creating the pubSubClient")
	t.Require().NotNil(t.pubSubClient, "the pubSubClient must not be nil")

	t.queueClient, err = t.e.GetStreamClient()
	t.Require().NoError(err, "there must not be an error when getting a stream pubSubClient")
	t.Require().NotNil(t.queueClient, "the queue pubSubClient must not be nil")
}

func (t *RaftEventHandlerTestSuite) TestRegisterCallback() {
	ev := NewRaftEventHandler(t.pubSubClient, t.queueClient, t.logger)

	t.Require().NotPanics(func() {
		ev.RegisterCallback("test", raftv1.Event_EVENT_NODE_READY, func(event *raftv1.RaftEvent) {
			t.Require().NotNil(event, "the event payload must not be nil")
			t.Require().Equal(raftv1.Event_EVENT_NODE_READY, event.Event, "the event must match the expected type")
		})
	})

	t.Require().NotNil(ev.cbTable, "the callback table must not be nil")
	t.Require().NotEmpty(ev.cbTable[raftv1.Event_EVENT_NODE_READY]["test"], "there must be a function named 'test'")
}

func (t *RaftEventHandlerTestSuite) TestUnregisterCallback() {
	ev := NewRaftEventHandler(t.pubSubClient, t.queueClient, t.logger)

	t.Require().NotPanics(func() {
		ev.RegisterCallback("test", raftv1.Event_EVENT_NODE_READY, func(event *raftv1.RaftEvent) {
			t.Require().NotNil(event, "the event payload must not be nil")
			t.Require().Equal(raftv1.Event_EVENT_NODE_READY, event.Event, "the event must match the expected type")
		})
	}, "registering a call back must not panic")

	t.Require().NotNil(ev.cbTable, "the callback table must not be nil")
	t.Require().NotEmpty(ev.cbTable[raftv1.Event_EVENT_NODE_READY]["test"], "there must be a function named 'test'")

	t.Require().NotPanics(func() {
		ev.UnregisterCallback("test", raftv1.Event_EVENT_NODE_READY)
	}, "unregistering a callback must not panic")

	t.Require().Empty(ev.cbTable[raftv1.Event_EVENT_NODE_READY]["test"], "there must not be a function under 'test")
}

func (t *RaftEventHandlerTestSuite) TestCallback() {
	ev := NewRaftEventHandler(t.pubSubClient, t.queueClient, t.logger)

	called := 0
	t.Require().NotPanics(func() {
		ev.RegisterCallback("test", raftv1.Event_EVENT_NODE_READY, func(event *raftv1.RaftEvent) {
			t.Require().NotNil(event, "the event payload must not be nil")
			t.Require().Equal(raftv1.Event_EVENT_NODE_READY, event.Action, "the event must match the expected type")
			called += 1
		})
	}, "registering a call back must not panic")

	payload := &raftv1.RaftEvent{
		Typ:       raftv1.EventType_EVENT_TYPE_NODE,
		Action:    raftv1.Event_EVENT_NODE_READY,
		Timestamp: timestamppb.Now(),
		Event: &raftv1.RaftEvent_Node{
			Node: &raftv1.RaftNodeEvent{
				ShardId:   10,
				ReplicaId: 10,
			},
		},
	}

	go ev.Run()
	utils.Wait(100 * time.Millisecond)

	buf, _ := payload.MarshalVT()
	err := t.pubSubClient.Publish(RaftNodeSubject, buf)
	t.Require().NoError(err, "there must not be a publishing error")

	utils.Wait(100 * time.Millisecond)
	t.Require().Equal(1, called, "the callback must have been triggered at most once")

	ev.Stop()
}
