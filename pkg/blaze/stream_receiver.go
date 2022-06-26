/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/server"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	configv1 "r3t.io/pleiades/pkg/protocols/config/v1"
)

// StreamReceiver takes a new stream and returns a *server.Server reference from the registry
type StreamReceiver struct {
	logger zerolog.Logger
	registry *Registry
}

func NewStreamReceiver(logger zerolog.Logger, registry *Registry) (*StreamReceiver, error) {
	l := logger.With().Str("component", "stream_receiver").Logger()

	return &StreamReceiver{logger: l, registry: registry}, nil
}

func (sr *StreamReceiver) Receive(stream quic.Stream) (*server.Server, error) {
	localLogger := sr.logger.With().Int64("stream_id", int64(stream.StreamID())).Logger()

	var buf bytes.Buffer
	n, err := io.CopyN(&buf, stream, readBufferSize)
	if err != nil {
		localLogger.Err(err).Msg("error reading stream")
		return nil, err
	}

	if n != readBufferSize {
		if n < readBufferSize {
			err := errors.New("read length is less than read buffer size")
			localLogger.Err(err).Int64("read_length", n).Msg("can't determine server when the read buffer isn't full")
			return nil, err
		}
	}

	if n >= readBufferSize {
		sr.logger.Trace().Int64("read_length", n).Msg("read length might be greater than read buffer size, will try to determine service")

		msg,err := capnp.Unmarshal(buf.Bytes())
		if err != nil {
			localLogger.Err(err).Msg("error unmarshalling capnp message")
			return nil, err
		}

		svcType, err := configv1.ReadRootServiceType(msg)
		if err != nil {
			localLogger.Err(err).Msg("error reading service type")
			return nil, err
		}

		t := svcType.Type().String()
		localLogger.Trace().Str("service_type", t).Msg("discovered service type")

		switch t {
		case "configService":
			localLogger.Trace().Msg("config service was determined")
			return sr.checkRegistry(t)
		case "test":
			localLogger.Trace().Msg("test service was determined")
			return sr.checkRegistry(t)
		default:
			localLogger.Trace().Str("service_type", t).Msg("unknown service type")
			return nil, fmt.Errorf("unknown service type of %s", t)
		}
	}

	return nil, errors.New("cannot determine service type")
}

func (sr *StreamReceiver) checkRegistry(svc string) (*server.Server, error) {
	target, err := sr.registry.Get(svc)
	if err != nil {
		sr.logger.Err(err).Str("service", svc).Msg("error getting server from registry")
		return nil, err
	}
	if target == nil {
		err := errors.New("server is nil")
		sr.logger.Err(err).Str("service", svc).Msg("error getting server from registry")
		return nil, err
	}

	return target, nil
}
