package blaze

import (
	"bytes"
	"context"
	"strings"
	"sync"

	"capnproto.org/go/capnp/v3"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	configv1 "r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/services/v1/config"
)

var (
	readBufferSize int64 = 0
)

func init() {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	svcType, err := configv1.NewRootServiceType(seg)
	if err != nil {
		panic(err)
	}
	svcType.SetType(configv1.ServiceType_Type_configService)

	var buf bytes.Buffer
	err = capnp.NewEncoder(&buf).Encode(msg)
	if err != nil {
		panic(err)
	}

	readBufferSize = int64(buf.Len())
}

type Server struct {
	listener quic.Listener
	logger   zerolog.Logger
	closed   bool
	mu sync.RWMutex
	registry *config.Registry
}

func NewServer(listener quic.Listener, logger zerolog.Logger, registry *config.Registry) *Server {
	l := logger.With().Str("component", "stream-manager").Logger()
	return &Server{listener: listener, logger: l, registry: registry}
}

func (s *Server) Start(ctx context.Context) error {
	rootContext := context.WithValue(ctx, "component", "stream-manager")
	s.logger.Info().Msg("starting listener")
	go s.handleConn(rootContext)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("shutting down listener")
	s.mu.Lock()
	s.closed = false
	s.mu.Unlock()
	return s.listener.Close()
}

func (s *Server) handleConn(ctx context.Context) {
	for {
		s.mu.RLock()
		if s.closed {
			return
		}
		s.mu.RUnlock()

		conn, err := s.listener.Accept(ctx)
		if err != nil {
			errStr := err.Error()

			// see if we're shutting down
			if strings.Contains(errStr, "context cancelled") || strings.Contains(errStr, "server closed") {
				return
			}

			s.logger.Debug().Err(err).Msg("error accepting remote connection")
			continue
		}

		streamCtx := context.WithValue(ctx, "remote-addr", conn.RemoteAddr().String())
		go s.handleStreams(conn, streamCtx)
	}
}

func (s *Server) handleStreams(conn quic.Connection, ctx context.Context) {
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			switch errData := err.(type) {
			case *quic.ApplicationError:

				// todo (sienna): figure out a better way to handle stream closures
				// ref: https://www.rfc-editor.org/rfc/rfc9000.html#section-20.1-2.2.1
				// the connection is closed with no errors, can't handle any more streams
				if errData.ErrorCode == quic.ApplicationErrorCode(quic.NoError) {
					return
				}
			default:
				s.logger.Err(err).Msg("an unknown error has occurred")
				continue
			}
		}

		downstreamLogger := s.logger.
			With().
			Str("remote-addr", conn.RemoteAddr().String()).
			Logger()
		handlerCtx := context.WithValue(ctx, "stream-id", stream.StreamID())

		go s.receiveStream(handlerCtx, downstreamLogger, stream)
	}
}

func (s *Server) receiveStream(ctx context.Context, inheritedLogger zerolog.Logger, stream quic.Stream) {
	logger := inheritedLogger.With().Str("component", "stream").Logger()

	sr, err := NewStreamReceiver(logger, s.registry)
	if err != nil {
		logger.Err(err).Msg("cannot create stream receiver")
		err := stream.Close()
		if err != nil {
			logger.Err(err).Msg("cannot close stream")
		}
		return
	}

	sr.Receive(ctx, logger, stream)
}
