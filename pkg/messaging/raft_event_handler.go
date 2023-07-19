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
	"github.com/mxplusb/pleiades/pkg/raftpb"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type RaftEventHandler struct {
	logger       zerolog.Logger
	pubSubClient *EmbeddedMessagingPubSubClient
	queueClient  *EmbeddedMessagingStreamClient
	cbTable      map[raftpb.Event]map[string]EventCallback
	sub          *nats.Subscription
	done         chan struct{}
}

type EventCallback func(event *raftpb.RaftEvent)

func NewRaftEventHandler(eventStreamClient *EmbeddedMessagingPubSubClient,
	queueClient *EmbeddedMessagingStreamClient,
	logger zerolog.Logger) *RaftEventHandler {

	// generate the callback table.
	cbTable := make(map[raftpb.Event]map[string]EventCallback)
	for _, cb := range []raftpb.Event{
		raftpb.Event_EVENT_UNSPECIFIED,
		raftpb.Event_EVENT_CONNECTION_ESTABLISHED,
		raftpb.Event_EVENT_CONNECTION_FAILED,
		raftpb.Event_EVENT_LOG_COMPACTED,
		raftpb.Event_EVENT_LOGDB_COMPACTED,
		raftpb.Event_EVENT_MEMBERSHIP_CHANGED,
		raftpb.Event_EVENT_NODE_HOST_SHUTTING_DOWN,
		raftpb.Event_EVENT_NODE_READY,
		raftpb.Event_EVENT_NODE_UNLOADED,
		raftpb.Event_EVENT_SEND_SNAPSHOT_ABORTED,
		raftpb.Event_EVENT_SEND_SNAPSHOT_COMPLETED,
		raftpb.Event_EVENT_SEND_SNAPSHOT_STARTED,
		raftpb.Event_EVENT_SNAPSHOT_COMPACTED,
		raftpb.Event_EVENT_SNAPSHOT_CREATED,
		raftpb.Event_EVENT_SNAPSHOT_RECEIVED,
		raftpb.Event_EVENT_SNAPSHOT_RECOVERED,
		raftpb.Event_EVENT_LEADER_UPDATED,
	} {
		cbTable[cb] = make(map[string]EventCallback)
	}

	return &RaftEventHandler{
		logger:       logger.With().Str("component", "raftpb-event-handler").Logger(),
		pubSubClient: eventStreamClient,
		queueClient:  queueClient,
		done:         make(chan struct{}, 1),
		cbTable:      cbTable,
	}
}

// RegisterCallback will add a named callback of a specific name to the callback table. It will overwrite a
// callback of the existing name.
func (r *RaftEventHandler) RegisterCallback(name string, action raftpb.Event, cb EventCallback) {
	r.cbTable[action][name] = cb
	r.logger.Debug().Str("callback", name).Str("action", action.String()).Msg("registered callback")
}

// UnregisterCallback removes a named callback from the callback table.
func (r *RaftEventHandler) UnregisterCallback(name string, action raftpb.Event) {
	delete(r.cbTable[action], name)
}

// Run this with `go Run()`
func (r *RaftEventHandler) Run() {
	var err error
	r.sub, err = r.pubSubClient.Subscribe("system.raftpb.*", func(msg *nats.Msg) {
		payload := &raftpb.RaftEvent{}
		err := payload.UnmarshalVT(msg.Data)
		if err != nil {
			r.logger.Error().
				Str("subject", msg.Subject).
				Str("reply", msg.Reply).
				Int("size", len(msg.Data)).
				Err(err).
				Msg("can't unmarshal event data")
			return
		}

		for k, callback := range r.cbTable[payload.Action] {
			r.logger.Trace().Interface("payload", payload).Str("callback", k).Msg("activated callback")
			go callback(payload)
		}
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("channel subscription failed")
		return
	}
	//defer func(sub *nats.Subscription) {
	//	err := sub.Unsubscribe()
	//	if err != nil {
	//		r.logger.Error().Err(err).Msg("unsubscription failed")
	//	}
	//	close(listener)
	//}(sub)

	r.logger.Debug().Msg("raftpb event handler callback defined")

	//for event := range listener {
	//	payload := &raftpb.RaftEvent{}
	//	err := payload.UnmarshalVT(event.Data)
	//	if err != nil {
	//		r.logger.Error().
	//			Str("subject", event.Subject).
	//			Str("reply", event.Reply).
	//			Int("size", len(event.Data)).
	//			Err(err).
	//			Msg("can't unmarshal event data")
	//		continue
	//	}
	//
	//	for k, callback := range r.cbTable[payload.Action] {
	//		r.logger.Trace().Interface("payload", payload).Str("callback", k).Msg("activated callback")
	//		go callback(payload)
	//	}
	//
	//	// check if we're done or not
	//	select {
	//	case <-r.done:
	//		return
	//	default:
	//		break
	//	}
	//}
}

func (r *RaftEventHandler) Stop() {
	if err := r.sub.Unsubscribe(); err != nil {
		r.logger.Error().Err(err).Msg("error unsubscribing from nats")
	}
}
