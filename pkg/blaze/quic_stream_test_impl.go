/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"bytes"
	"context"
	"time"

	"github.com/lucas-clemente/quic-go"
)

type TestStream struct {
	buf *bytes.Buffer
}

func NewTestStream(bufSize int) *TestStream {
	return &TestStream{
		buf: bytes.NewBuffer(make([]byte, bufSize)),
	}
}

func (t *TestStream) Read(p []byte) (n int, err error) {
	return t.buf.Read(p)
}

func (t *TestStream) CancelRead(code quic.StreamErrorCode) {
	// does nothing
}

func (t *TestStream) SetReadDeadline(time time.Time) error {
	return nil
}

func (t *TestStream) StreamID() quic.StreamID {
	return 100
}

func (t *TestStream) Write(p []byte) (n int, err error) {
	return t.buf.Write(p)
}

func (t *TestStream) Close() error {
	t.buf.Reset()
	return nil
}

func (t *TestStream) CancelWrite(code quic.StreamErrorCode) {
	// do nothing
}

func (t *TestStream) Context() context.Context {
	return context.Background()
}

func (t *TestStream) SetWriteDeadline(time time.Time) error {
	return nil
}

func (t *TestStream) SetDeadline(time time.Time) error {
	return nil
}
