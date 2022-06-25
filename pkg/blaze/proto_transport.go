/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"io"

	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
)

var (
	_ io.ReadWriteCloser = (*ProtoTransport)(nil)
)

type ProtoTransport struct {
	logger zerolog.Logger
	stream quic.Stream
}

func NewProtoTransport(logger zerolog.Logger, stream quic.Stream) *ProtoTransport {
	l := logger.With().Int64("stream_id", int64(stream.StreamID())).Str("aspect", "stream").Logger()
	return &ProtoTransport{logger: l, stream: stream}
}

func (pt *ProtoTransport) Read(p []byte) (n int, err error) {
	pt.logger.Trace().Int64("stream_id", int64(pt.stream.StreamID())).Msg("reading from stream")
	return pt.stream.Read(p)
}

func (pt *ProtoTransport) Write(p []byte) (n int, err error) {
	pt.logger.Trace().Int64("stream_id", int64(pt.stream.StreamID())).Msg("writing to stream")
	return pt.stream.Write(p)
}

func (pt *ProtoTransport) Close() error {
	pt.logger.Trace().Int64("stream_id", int64(pt.stream.StreamID())).Msg("closing stream")
	return pt.stream.Close()
}
