package blaze

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"capnproto.org/go/capnp/v3"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	configv1 "r3t.io/pleiades/pkg/protocols/config/v1"
	"storj.io/drpc/drpcmanager"
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
	router   *Router
	manager  *drpcmanager.Manager
	logger   zerolog.Logger
	closed   bool
	mu sync.RWMutex
}

func NewServer(listener quic.Listener, router *Router, logger zerolog.Logger) *Server {
	return &Server{listener: listener, router: router, logger: logger}
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
			Str("stream-id", fmt.Sprintf("%d", uint64(stream.StreamID()))).
			Logger()
		handlerCtx := context.WithValue(ctx, "stream-id", stream.StreamID())

		go s.serveOne(stream, downstreamLogger, handlerCtx)
	}
}

func (s *Server) serveOne(stream quic.Stream, inheritedLogger zerolog.Logger, ctx context.Context) {
	streamManager := drpcmanager.New(stream)

	for {
		dServerStream, rpc, err := streamManager.NewServerStream(ctx)
		if err != nil {
			// todo (sienna): figure out a better way to handle manager closures
			if strings.Contains(err.Error(), "manager closed") {
				if streamErr := stream.Close(); streamErr != nil {
					inheritedLogger.Trace().Err(err).Msg("cannot close stream")
				}
				return
			}
			inheritedLogger.Err(err).Msg("new server stream cannot be created")
			return
		}

		err = s.router.HandleRPC(dServerStream, rpc)
		if err != nil {
			sendErr := dServerStream.SendError(err)
			inheritedLogger.Debug().Err(sendErr).Msg("failure to send error to client")

		}
	}
}
