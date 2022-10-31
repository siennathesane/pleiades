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
	"testing"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
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
	e              *EmbeddedMessaging
	client         *EmbeddedMessagingPubSubClient
	queueClient    *EmbeddedMessagingStreamClient
	defaultTimeout time.Duration
}

func (t *RaftEventHandlerTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 500 * time.Millisecond

	var err error
	t.e, err = NewEmbeddedMessagingWithDefaults(t.logger)
	t.Require().NoError(err, "there must not be an error creating the event stream")

	t.e.Start()
}

func (t *RaftEventHandlerTestSuite) TearDownSuite() {
	//t.e = nil
	//t.pubSubClient = nil
	//t.queueClient = nil
	//runtime.GC()
}

func (t *RaftEventHandlerTestSuite) SetupTest() {
	var err error
	t.client, err = t.e.GetPubSubClient()
	t.Require().NoError(err, "there must not be an error creating the eventStreamClient")
	t.Require().NotNil(t.client, "the eventStreamClient must not be nil")

	t.queueClient, err = t.e.GetStreamClient()
	t.Require().NoError(err, "there must not be an error when getting a stream pubSubClient")
	t.Require().NotNil(t.queueClient, "the queue pubSubClient must not be nil")
}

func (t *RaftEventHandlerTestSuite) TestWaitForMembershipChange() {
	testShardId := uint64(10)

	eh := NewRaftEventHandler(t.client, t.queueClient, t.logger)

	results := make(chan *raftv1.RaftEvent, 1)
	go eh.WaitForMembershipChange(testShardId, results, utils.Timeout(300000*time.Millisecond))

	for i := uint64(0); i < 500; i++ {
		payload := &raftv1.RaftEvent{
			Typ:       raftv1.EventType_EVENT_TYPE_NODE,
			Action:    raftv1.Event_EVENT_MEMBERSHIP_CHANGED,
			Timestamp: timestamppb.Now(),
			Event: &raftv1.RaftEvent_Node{Node: &raftv1.RaftNodeEvent{
				ShardId:   i,
				ReplicaId: i * 10,
			}},
		}
		buf, _ := payload.MarshalVT()
		err := t.client.Publish(RaftNodeSubject, buf)
		t.Require().NoError(err, "there must not be an error when publishing an event")
		utils.Wait(1 * time.Millisecond)
	}

	result := <-results
	t.Require().NotNil(result, "the result of waiting for a membership change must not be nil")
}
