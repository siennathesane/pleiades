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
	"fmt"
	"runtime"
	"time"

	"github.com/mxplusb/pleiades/pkg/messaging/clients"
	"github.com/mxplusb/pleiades/pkg/messaging/raft"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

var (
	singleton *EmbeddedMessaging
)

type EmbeddedMessagingStreamOpts struct {
	*server.Options
	timeout time.Duration
}

func NewEmbeddedMessagingWithDefaults(logger zerolog.Logger) (*EmbeddedMessaging, error) {
	if singleton != nil {
		return singleton, nil
	}

	opts := &server.Options{
		Host:          "localhost",
		JetStream:     true,
		DontListen:    true,
		WriteDeadline: 1_000 * time.Millisecond,
	}
	srv, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	singleton = &EmbeddedMessaging{
		opts: &EmbeddedMessagingStreamOpts{timeout: 4000 * time.Millisecond, Options: opts},
		srv:  srv,
	}

	// set the right logging level
	level := zerolog.GlobalLevel()
	switch level {
	case zerolog.TraceLevel:
		srv.SetLoggerV2(&messagingLogger{
			logger: logger.With().Str("component", "messaging").Logger(),
		}, true, true, true)
		break
	case zerolog.DebugLevel:
		srv.SetLoggerV2(&messagingLogger{
			logger: logger.With().Str("component", "messaging").Logger(),
		}, true, false, false)
		break
	default:
		srv.SetLoggerV2(&messagingLogger{
			logger: logger.With().Str("component", "messaging").Logger(),
		}, false, false, false)
	}

	// ensure it's properly stopped when all references to this are deleted.
	runtime.SetFinalizer(singleton, func(e *EmbeddedMessaging) {
		e.Stop()
	})

	return singleton, err
}

type EmbeddedMessaging struct {
	opts *EmbeddedMessagingStreamOpts
	srv  *server.Server
}

func (ev *EmbeddedMessaging) Start() {
	// verify it's already started
	if ev.srv.ReadyForConnections(100 * time.Millisecond) {
		return
	}

	go ev.srv.Start()
	if !ev.srv.ReadyForConnections(ev.opts.timeout) {
		panic("nats won't start")
	}

	client, err := ev.GetStreamClient()
	if err != nil {
		panic(fmt.Errorf("can't get stream pubSubClient: %s", err))
	}

	_, err = client.AddStream(&nats.StreamConfig{
		Name:        raft.SystemStreamName,
		Description: "All internal system streams",
		Subjects: []string{
			raft.RaftHostSubject,
			raft.RaftLogSubject,
			raft.RaftNodeSubject,
			raft.RaftSnapshotSubject,
			raft.RaftConnectionSubject,
			raft.RaftSubject,
			raft.SystemStreamName,
		},
		Retention: nats.WorkQueuePolicy,
		Discard:   nats.DiscardOld,
		Storage:   nats.MemoryStorage,
	})
	if err != nil {
		panic(fmt.Errorf("can't add system stream: %s", err))
	}
}

func (ev *EmbeddedMessaging) Stop() {
	ev.srv.Shutdown()
}

func (ev *EmbeddedMessaging) GetPubSubClient() (*clients.EmbeddedMessagingPubSubClient, error) {
	conn, err := nats.Connect(ev.srv.ClientURL(), nats.InProcessServer(ev.srv))
	return &clients.EmbeddedMessagingPubSubClient{conn}, err
}

func (ev *EmbeddedMessaging) GetStreamClient() (*clients.EmbeddedMessagingStreamClient, error) {
	conn, err := nats.Connect(ev.srv.ClientURL(), nats.InProcessServer(ev.srv))
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	return &clients.EmbeddedMessagingStreamClient{js}, err
}
