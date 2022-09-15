/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package pubsub

import (
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type EmbeddedEventStreamOpts struct {
	*server.Options
	timeout time.Duration
}

type EmbeddedQueueClient struct {
	nats.JetStreamContext
}

type EmbeddedEventStreamClient struct {
	*nats.Conn
}

func NewEmbeddedEventStream(opts *EmbeddedEventStreamOpts) (*EmbeddedEventStream, error) {
	// todo (sienna): ensure that StoreDir is set
	srv, err := server.NewServer(opts.Options)
	return &EmbeddedEventStream{
		opts: opts,
		srv:  srv,
	}, err
}

type EmbeddedEventStream struct {
	opts *EmbeddedEventStreamOpts
	srv  *server.Server
}

func (ev *EmbeddedEventStream) Start() {
	go ev.srv.Start()
	if !ev.srv.ReadyForConnections(ev.opts.timeout) {
		panic("event server won't start")
	}
}

func (ev *EmbeddedEventStream) Stop() {
	ev.srv.Shutdown()
}

func (ev *EmbeddedEventStream) GetClient() (*EmbeddedEventStreamClient, error) {
	conn, err := nats.Connect(ev.srv.ClientURL(), nats.InProcessServer(ev.srv))
	return &EmbeddedEventStreamClient{conn}, err
}
