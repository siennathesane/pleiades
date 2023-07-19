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
	"context"
	"fmt"
	"time"

	nserv "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var EmbeddedMessagingModule = fx.Module("embedded_messaging", fx.Invoke(NewEmbeddedMessagingWithDefaults))

var (
	singleton *EmbeddedMessaging
)

type embeddedMessagingStreamOpts struct {
	*nserv.Options
	timeout time.Duration
}

type EmbeddedMessagingStreamClient struct {
	nats.JetStreamContext
}

type EmbeddedMessagingPubSubClient struct {
	*nats.Conn
}

type EmbeddedMessagingWithDefaultsParams struct {
	fx.In

	Logger    zerolog.Logger
	Lifecycle fx.Lifecycle
}

type EmbeddedMessagingWithDefaultsResults struct {
	fx.Out

	EmbeddedMessaging *EmbeddedMessaging
}

func NewEmbeddedMessagingWithDefaults(params EmbeddedMessagingWithDefaultsParams) (EmbeddedMessagingWithDefaultsResults, error) {
	if singleton != nil {
		return EmbeddedMessagingWithDefaultsResults{
			EmbeddedMessaging: singleton,
		}, nil
	}

	opts := &nserv.Options{
		Host:          "localhost",
		JetStream:     true,
		DontListen:    true,
		WriteDeadline: 1_000 * time.Millisecond,
	}
	srv, err := nserv.NewServer(opts)
	if err != nil {
		return EmbeddedMessagingWithDefaultsResults{}, err
	}

	singleton = &EmbeddedMessaging{
		opts: &embeddedMessagingStreamOpts{timeout: 4000 * time.Millisecond, Options: opts},
		srv:  srv,
	}

	// set the right logging level
	level := zerolog.GlobalLevel()
	switch level {
	case zerolog.TraceLevel:
		srv.SetLoggerV2(&messagingLogger{
			logger: params.Logger.With().Str("component", "messaging").Logger(),
		}, true, true, true)
		break
	case zerolog.DebugLevel:
		srv.SetLoggerV2(&messagingLogger{
			logger: params.Logger.With().Str("component", "messaging").Logger(),
		}, true, false, false)
		break
	default:
		srv.SetLoggerV2(&messagingLogger{
			logger: params.Logger.With().Str("component", "messaging").Logger(),
		}, false, false, false)
	}

	singleton.Start()

	// register the stop function
	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if singleton == nil {
				return nil
			}
			if singleton.srv == nil {
				return nil
			}
			// todo (sienna): figure out why this is panicking
			singleton.Stop()
			return nil
		},
	})

	return EmbeddedMessagingWithDefaultsResults{
		EmbeddedMessaging: singleton,
	}, err
}

type EmbeddedMessaging struct {
	opts *embeddedMessagingStreamOpts
	srv  *nserv.Server
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
		Name:        SystemStreamName,
		Description: "All internal system streams",
		Subjects: []string{
			RaftHostSubject,
			RaftLogSubject,
			RaftNodeSubject,
			RaftSnapshotSubject,
			RaftConnectionSubject,
			RaftSubject,
			SystemStreamName,
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

func (ev *EmbeddedMessaging) GetPubSubClient() (*EmbeddedMessagingPubSubClient, error) {
	conn, err := nats.Connect(ev.srv.ClientURL(), nats.InProcessServer(ev.srv))
	return &EmbeddedMessagingPubSubClient{conn}, err
}

func (ev *EmbeddedMessaging) GetStreamClient() (*EmbeddedMessagingStreamClient, error) {
	conn, err := nats.Connect(ev.srv.ClientURL(), nats.InProcessServer(ev.srv))
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	return &EmbeddedMessagingStreamClient{js}, err
}
