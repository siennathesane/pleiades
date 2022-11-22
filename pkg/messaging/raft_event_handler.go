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
	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type RaftEventHandler struct {
	logger       zerolog.Logger
	pubSubClient *EmbeddedMessagingPubSubClient
	queueClient  *EmbeddedMessagingStreamClient
	cbTable      map[raftv1.Event]map[string]EventCallback
	done         chan struct{}
}

type EventCallback func(event *raftv1.RaftEvent)

func NewRaftEventHandler(eventStreamClient *EmbeddedMessagingPubSubClient,
	queueClient *EmbeddedMessagingStreamClient,
	logger zerolog.Logger) *RaftEventHandler {

	// generate the callback table.
	cbTable := make(map[raftv1.Event]map[string]EventCallback)
	for _, cb := range []raftv1.Event{
		raftv1.Event_EVENT_UNSPECIFIED,
		raftv1.Event_EVENT_CONNECTION_ESTABLISHED,
		raftv1.Event_EVENT_CONNECTION_FAILED,
		raftv1.Event_EVENT_LOG_COMPACTED,
		raftv1.Event_EVENT_LOGDB_COMPACTED,
		raftv1.Event_EVENT_MEMBERSHIP_CHANGED,
		raftv1.Event_EVENT_NODE_HOST_SHUTTING_DOWN,
		raftv1.Event_EVENT_NODE_READY,
		raftv1.Event_EVENT_NODE_UNLOADED,
		raftv1.Event_EVENT_SEND_SNAPSHOT_ABORTED,
		raftv1.Event_EVENT_SEND_SNAPSHOT_COMPLETED,
		raftv1.Event_EVENT_SEND_SNAPSHOT_STARTED,
		raftv1.Event_EVENT_SNAPSHOT_COMPACTED,
		raftv1.Event_EVENT_SNAPSHOT_CREATED,
		raftv1.Event_EVENT_SNAPSHOT_RECEIVED,
		raftv1.Event_EVENT_SNAPSHOT_RECOVERED,
		raftv1.Event_EVENT_LEADER_UPDATED,
	} {
		cbTable[cb] = make(map[string]EventCallback)
	}

	return &RaftEventHandler{
		logger:       logger.With().Str("component", "raft-event-handler").Logger(),
		pubSubClient: eventStreamClient,
		queueClient:  queueClient,
		done:         make(chan struct{}, 1),
		cbTable:      cbTable,
	}
}

// RegisterCallback will add a named callback of a specific name to the callback table. It will overwrite a
// callback of the existing name.
func (r *RaftEventHandler) RegisterCallback(name string, action raftv1.Event, cb EventCallback) {
	r.cbTable[action][name] = cb
	r.logger.Debug().Str("callback", name).Str("action", action.String()).Msg("registered callback")
}

// UnregisterCallback removes a named callback from the callback table.
func (r *RaftEventHandler) UnregisterCallback(name string, action raftv1.Event) {
	delete(r.cbTable[action], name)
}

// Run this with `go Run()`
func (r *RaftEventHandler) Run() {
	listener := make(chan *nats.Msg)
	sub, err := r.pubSubClient.ChanSubscribe("system.raftv1.*", listener)
	if err != nil {
		r.logger.Error().Err(err).Msg("channel subscription failed")
		return
	}
	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			r.logger.Error().Err(err).Msg("unsubscription failed")
		}
		close(listener)
	}(sub)

	for event := range listener {
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

		for k, callback := range r.cbTable[payload.Action] {
			r.logger.Trace().Interface("payload", payload).Str("callback", k).Msg("activated callback")
			go callback(payload)
		}

		// check if we're done or not
		select {
		case <-r.done:
			return
		default:
		}
	}
}

func (r *RaftEventHandler) Stop() {
	r.done <- struct{}{}
}
