package blaze

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	"storj.io/drpc"
	"storj.io/drpc/drpcmanager"
)

type ConnectionManager struct {
	listener quic.Listener
	handler  drpc.Handler
	manager  *drpcmanager.Manager
	logger   zerolog.Logger
	closed bool
}

func NewConnectionServer(listener quic.Listener, handler drpc.Handler, logger zerolog.Logger) *ConnectionManager {
	return &ConnectionManager{listener: listener, handler: handler, logger: logger}
}

func (s *ConnectionManager) Start(ctx context.Context) error {
	rootContext := context.WithValue(ctx, "component", "stream-manager")
	s.logger.Info().Msg("starting listener")
	go s.handleConn(rootContext)
	return nil
}

func (s *ConnectionManager) Stop(ctx context.Context) error {
	s.logger.Info().Msg("shutting down listener")
	s.closed = false
	return s.listener.Close()
}

func (s *ConnectionManager) handleConn(ctx context.Context) {
	for {
		if s.closed {
			return
		}

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

func (s *ConnectionManager) handleStreams(conn quic.Connection, ctx context.Context) {
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			switch errData := err.(type) {
			case *quic.ApplicationError:

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
			Str("stream-id", fmt.Sprintf("%d", uint64(stream.StreamID()))).
			Str("remote-addr", conn.RemoteAddr().String()).
			Logger()

		mPlexStream := NewMultiplexedStream(stream, s.handler, downstreamLogger)
		handlerCtx := context.WithValue(ctx, "stream-id", stream.StreamID())
		go func() {
			mPlexStream.Handle(handlerCtx)
		}()

		select {
		case <- ctx.Done():
			return
		default:
			continue
		}
	}
}
