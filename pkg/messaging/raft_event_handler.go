/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

// nb (sienna): this feels like an exercise in branch prediction failure.
// todo (sienna): there's gotta be a way to optimize this

package messaging

import (
	"time"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type RaftEventHandler struct {
	logger            zerolog.Logger
	eventStreamClient *EmbeddedMessagingPubSubClient
	queueClient       *EmbeddedMessagingStreamClient
}

func NewRaftEventHandler(eventStreamClient *EmbeddedMessagingPubSubClient, queueClient *EmbeddedMessagingStreamClient, logger zerolog.Logger) *RaftEventHandler {
	return &RaftEventHandler{logger: logger.With().Str("component", "raft-event-handler").Logger(), eventStreamClient: eventStreamClient, queueClient: queueClient}
}

func (r *RaftEventHandler) WaitForLeaderUpdate(shardId uint64, results chan *raftv1.RaftEvent, timeout time.Duration) {
	listener := make(chan *nats.Msg)
	sub, err := r.eventStreamClient.ChanSubscribe(RaftSubject, listener)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't subscribe to raft leader updates")
	}
	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			r.logger.Error().Err(err).Str("subject", RaftSubject).Msg("error unsubscribing")
		}
	}(sub)

	expiry := time.Now().Add(timeout)

	for event := range listener {
		// if we've waited long enough, return
		if time.Now().UnixMilli() > expiry.UnixMilli() {
			r.logger.Debug().Uint64("shard", shardId).Msg("timeout reached listening for leader update")
			return
		}

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

		if payload.GetAction() == raftv1.Event_EVENT_LEADER_UPDATED {
			leaderUpdate := payload.GetLeaderUpdate()
			// this is just defensive, it shouldn't be nil
			if leaderUpdate == nil {
				continue
			}

			// verify and return
			if leaderUpdate.GetShardId() == shardId {
				results <- payload
				return
			}
		}

		err = event.Nak()
		if err != nil {
			r.logger.Error().
				Err(err).
				Uint64("shard", shardId).
				Str("subject", RaftSubject).
				Msg("error negatively acknowledging the message")
		}
	}
	return
}

func (r *RaftEventHandler) WaitForMembershipChange(shardId uint64, results chan *raftv1.RaftEvent, timeout time.Duration) <-chan *raftv1.RaftEvent {
	listener := make(chan *nats.Msg)
	sub, err := r.eventStreamClient.ChanSubscribe(RaftNodeSubject, listener)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't subscribe to raft leader updates")
	}
	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			r.logger.Error().Err(err).Str("subject", RaftNodeSubject).Msg("error unsubscribing")
		}
		close(listener)
	}(sub)

	expiry := time.Now().Add(timeout)

	for event := range listener {
		// if we've waited long enough, return
		if time.Now().UnixMilli() > expiry.UnixMilli() {
			r.logger.Debug().Uint64("shard", shardId).Msg("timeout reached listening for leader update")
			return nil
		}

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
			if nodeEvent.GetShardId() == shardId {
				results <- payload
			}
		}

		err = event.Nak()
		if err != nil {
			if err == nats.ErrMsgNoReply {
				continue
			} else {
				r.logger.Error().
					Err(err).
					Uint64("shard", shardId).
					Str("subject", RaftNodeSubject).
					Msg("error negatively acknowledging the message")
			}
		}
	}

	return nil
}

func (r *RaftEventHandler) WaitForReplicaReady(shardId, replicaId uint64, results chan *raftv1.RaftEvent, timeout time.Duration) {
	listener := make(chan *nats.Msg)
	sub, err := r.eventStreamClient.ChanSubscribe(RaftNodeSubject, listener)
	if err != nil {
		r.logger.Error().Err(err).Msg("can't subscribe to raft leader updates")
	}
	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			r.logger.Error().Err(err).Str("subject", RaftNodeSubject).Msg("error unsubscribing")
		}
	}(sub)

	expiry := time.Now().Add(timeout)

	for event := range listener {
		// if we've waited long enough, return
		if time.Now().UnixMilli() > expiry.UnixMilli() {
			r.logger.Debug().Uint64("shard", shardId).Msg("timeout reached listening for leader update")
			return
		}

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

		if payload.GetAction() == raftv1.Event_EVENT_NODE_READY {
			nodeEvent := payload.GetNode()
			// this is just defensive, it shouldn't be nil
			if nodeEvent == nil {
				continue
			}

			// verify and return
			if nodeEvent.GetShardId() == shardId && nodeEvent.GetReplicaId() == replicaId {
				r.logger.Debug().Uint64("shard-id", shardId).Uint64("replica-id", nodeEvent.GetReplicaId()).Msg("replica is ready")
				results <- payload
				return
			}
		}

		err = event.Nak()
		if err != nil {
			r.logger.Error().
				Err(err).
				Uint64("shard", shardId).
				Str("subject", RaftNodeSubject).
				Msg("error negatively acknowledging the message")
		}
	}
	return
}
