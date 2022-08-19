/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"bytes"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var _ network.Stream = (*testStream)(nil)

type testStream struct {
	buf bytes.Buffer
	allowRead, allowWrite bool
	id string
}

func (t *testStream) Read(p []byte) (n int, err error) {
	if !t.allowRead {
		return 0, errors.New("buffer is closed for reading")
	}
	return t.buf.Read(p)
}

func (t *testStream) Write(p []byte) (n int, err error) {
	if !t.allowWrite {
		return 0, errors.New("buffer is closed for writing")
	}
	return t.buf.Write(p)
}

func (t *testStream) Close() error {
	t.buf.Reset()
	t.allowRead = false
	t.allowWrite = false
	return nil
}

func (t *testStream) CloseWrite() error {
	t.allowWrite = false
	return nil
}

func (t *testStream) CloseRead() error {
	t.allowRead = false
	return nil
}

func (t *testStream) Reset() error {
	t.buf.Reset()
	return nil
}

func (t *testStream) SetDeadline(time time.Time) error {
	return nil
}

func (t *testStream) SetReadDeadline(time time.Time) error {
	return nil
}

func (t *testStream) SetWriteDeadline(time time.Time) error {
	return nil
}

func (t *testStream) ID() string {
	return t.id
}

func (t *testStream) Protocol() protocol.ID {
	return protocol.ID(t.id)
}

func (t *testStream) SetProtocol(id protocol.ID) error {
	t.id = string(id)
	return nil
}

func (t *testStream) Stat() network.Stats {
	return network.Stats{}
}

func (t *testStream) Conn() network.Conn {
	return nil
}

func (t *testStream) Scope() network.StreamScope {
	return nil
}
