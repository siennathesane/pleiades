/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"bytes"
	"context"
	"errors"
	"io"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	configv1 "r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/services/v1/config"
)

const (
	ConfigService configv1.ServiceType_Type = configv1.ServiceType_Type_configService
	TestService configv1.ServiceType_Type = configv1.ServiceType_Type_test
)

// StreamReceiver takes a new stream and returns a *server.Server reference from the registry
type StreamReceiver struct {
	logger zerolog.Logger
	registry *config.Registry
}

func NewStreamReceiver(logger zerolog.Logger, registry *config.Registry) (*StreamReceiver, error) {
	l := logger.With().Str("component", "stream_receiver").Logger()

	return &StreamReceiver{logger: l, registry: registry}, nil
}

func (sr *StreamReceiver) Receive(ctx context.Context, logger zerolog.Logger, stream quic.Stream) {
	localLogger := logger.With().Int64("stream_id", int64(stream.StreamID())).Logger()

	var buf bytes.Buffer
	n, err := io.CopyN(&buf, stream, readBufferSize)
	if err != nil {
		localLogger.Err(err).Msg("error reading stream")
		return
	}

	if n != readBufferSize {
		if n < readBufferSize {
			err := errors.New("read length is less than read buffer size")
			localLogger.Err(err).Int64("read-length", n).Msg("can't determine server when the read buffer isn't full")
			return
		}
	}

	if n >= readBufferSize {
		sr.logger.Trace().Int64("read-length", n).Msg("read length might be greater than read buffer size, will try to determine service")

		msg,err := capnp.Unmarshal(buf.Bytes())
		if err != nil {
			localLogger.Err(err).Msg("error unmarshalling capnp message")
			return
		}

		svcType, err := configv1.ReadRootServiceType(msg)
		if err != nil {
			localLogger.Err(err).Msg("error reading service type")
			return
		}

		t := svcType.Type()
		localLogger.Trace().Str("service-type", t.String()).Msg("discovered service type")

		switch t {
		case ConfigService:
			localLogger.Trace().Msg("config service was determined")
			sr.handleConfigService(stream.Context(), localLogger, stream)
			return
		case TestService:
			localLogger.Trace().Msg("test service was determined")
			return
		default:
			localLogger.Trace().Str("service_type", t.String()).Msg("unknown service type")
			return
		}
	}

	return
}

func (sr *StreamReceiver) checkRegistry(svc configv1.ServiceType_Type) (any, error) {
	target, err := sr.registry.GetServer(svc)
	if err != nil {
		sr.logger.Err(err).Uint16("service", uint16(svc)).Msg("error getting server from registry")
		return 0, err
	}
	if target == nil {
		err := errors.New("server is nil")
		sr.logger.Err(err).Uint16("service", uint16(svc)).Msg("error getting server from registry")
		return 0, err
	}

	return target, nil
}

func (sr *StreamReceiver) handleConfigService(ctx context.Context, logger zerolog.Logger, stream quic.Stream) {
	l := logger.With().Str("aspect", "config-service").Logger()
	main, err := sr.registry.GetServer(configv1.ServiceType_Type_configService)
	if err != nil {
		l.Err(err).Msg("cannot get config service")
		return
	}

	configService := main.(*config.ConfigServer)
	if configService == nil {
		l.Err(errors.New("config service is nil")).Msg("cannot build config service")
		return
	}

	configServiceClientBuilder := configv1.ConfigService_ServerToClient(configService, &server.Policy{MaxConcurrentCalls: 250})

	conn := rpc.NewConn(rpc.NewStreamTransport(stream), &rpc.Options{BootstrapClient: configServiceClientBuilder.Client})
	defer func() {
		err := conn.Close()
		if err != nil {
			l.Err(err).Msg("error closing connection")
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-ctx.Done():
		err = conn.Close()
		if err != nil {
			l.Err(err).Msg("error closing connection")
		}
		return
	}
}