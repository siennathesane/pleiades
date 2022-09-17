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

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type EmbeddedMessagingStreamOpts struct {
	*server.Options
	timeout time.Duration
}

type EmbeddedMessagingStreamClient struct {
	nats.JetStreamContext
}

type EmbeddedMessagingPubSubClient struct {
	*nats.Conn
}

func NewEmbeddedMessaging(opts *EmbeddedMessagingStreamOpts) (*EmbeddedMessaging, error) {
	// todo (sienna): ensure that StoreDir is set
	srv, err := server.NewServer(opts.Options)
	return &EmbeddedMessaging{
		opts: opts,
		srv:  srv,
	}, err
}

type EmbeddedMessaging struct {
	opts *EmbeddedMessagingStreamOpts
	srv  *server.Server
}

func (ev *EmbeddedMessaging) Start() {
	go ev.srv.Start()
	if !ev.srv.ReadyForConnections(ev.opts.timeout) {
		panic("event server won't start")
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