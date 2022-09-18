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
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type RaftEventHandler struct {
	logger            zerolog.Logger
	eventStreamClient *EmbeddedMessagingPubSubClient
	queueClient       *EmbeddedMessagingStreamClient
}

func NewRaftEventHandler(eventStreamClient *EmbeddedMessagingPubSubClient, queueClient *EmbeddedMessagingStreamClient,logger zerolog.Logger) *RaftEventHandler{
	return &RaftEventHandler{logger: logger.With().Str("aspect", "raft-event-handler").Logger(), eventStreamClient: eventStreamClient, queueClient: queueClient}
}

func (r *RaftEventHandler) WaitForMembershipChange(shardId, replicaId uint64, timeout time.Duration) error {
	//expiration := time.Now().Add(timeout)

	listener := make(chan *nats.Msg)
	sub, err := r.eventStreamClient.ChanSubscribe(raftNodeSubject, listener)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't subscribe to raft node updates")
		return errors.Wrap(err, "can't subscribe to raft node updates")
	}
	defer sub.Unsubscribe()

	for event := range listener {
		//if expiration.UnixMilli() <= time.Now().UnixMilli() {
		//	return errors.New("listener timeout")
		//}

		payload := &raftv1.RaftEvent{}
		err := payload.UnmarshalVT(event.Data)
		if err != nil {
			r.logger.Error().
				Str("subject", event.Subject).
				Str("reply", event.Reply).
				Int("size", len(event.Data)).
				Err(err).
				Msg("can't unmarshal event data")
			continue
		}

		if payload.GetAction() == raftv1.Event_EVENT_MEMBERSHIP_CHANGED {
			nodeEvent := payload.GetNode()
			// this is just defensive, it shouldn't be nil
			if nodeEvent == nil {
				continue
			}

			// verify and return
			if nodeEvent.GetShardId() == shardId &&
				nodeEvent.GetReplicaId() == replicaId {
				return nil
			}
		}
	}

	return nil
}
