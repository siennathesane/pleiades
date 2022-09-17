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

	"github.com/mxplusb/pleiades/api/v1/raft"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/nats-io/nats-server/v2/server"
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

	opts := &EmbeddedMessagingStreamOpts{
		Options: &server.Options{
			Host:       "localhost",
			JetStream:  true,
			DontListen: true,
		},
		timeout: utils.Timeout(4000 * time.Millisecond),
	}

	var err error
	t.e, err = NewEmbeddedMessaging(opts)
	t.Require().NoError(err, "there must not be an error creating the event stream")

	t.e.Start()
}

func (t *RaftEventHandlerTestSuite) SetupTest() {
	var err error
	t.client, err = t.e.GetPubSubClient()
	t.Require().NoError(err, "there must not be an error creating the eventStreamClient")
	t.Require().NotNil(t.client, "the eventStreamClient must not be nil")

	t.queueClient, err = t.e.GetStreamClient()
	t.Require().NoError(err, "there must not be an error when getting a stream client")
	t.Require().NotNil(t.queueClient, "the queue client must not be nil")
}

func (t *RaftEventHandlerTestSuite) TestWaitForMembershipChange() {
	testShardId := uint64(10)
	testReplicaId := uint64(100)

	eh := NewRaftEventHandler(t.client, t.queueClient, t.logger)

	go func() {
		for i := uint64(0); i < 500; i++ {
			payload := &raft.RaftEvent{
				Typ:       raft.EventType_NODE,
				Action:    raft.Event_MEMBERSHIP_CHANGED,
				Timestamp: timestamppb.Now(),
				Event:     &raft.RaftEvent_Node{Node: &raft.RaftNodeEvent{
					ShardId:   i,
					ReplicaId: i*10,
				}},
			}
			buf, _ := payload.MarshalVT()
			err := t.client.Publish(raftNodeSubject, buf)
			t.Require().NoError(err, "there must not be an error when publishing an event")
			utils.Wait(1 * time.Millisecond)
		}
	}()

	err := eh.WaitForMembershipChange(testShardId, testReplicaId, utils.Timeout(100*time.Millisecond))
	t.Require().NoError(err, "there must not be an error when waiting for a specific membership change message")


}
