package blaze

import (
	"context"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	"storj.io/drpc"
	"storj.io/drpc/drpcenc"
	"storj.io/drpc/drpcmanager"
	"storj.io/drpc/drpcmetadata"
	"storj.io/drpc/drpcwire"
)

type ConnectionStream struct {
	stream    quic.Stream
	writer    *drpcwire.Writer
	handler   drpc.Handler
	manager   *drpcmanager.Manager
	logger    zerolog.Logger
	writeBuff []byte
	mu sync.RWMutex
}

func NewConnectionStream(stream quic.Stream, handler drpc.Handler, logger zerolog.Logger) *ConnectionStream {
	writer := drpcwire.NewWriter(stream, 0)
	manager := drpcmanager.New(stream)
	return &ConnectionStream{
		stream: stream,
		manager: manager,
		writer: writer,
		handler: handler,
		logger: logger,
		writeBuff: make([]byte, 0)}
}

func (m *ConnectionStream) Handle(ctx context.Context) {
	dServerStream, rpc, err := m.manager.NewServerStream(ctx)
	if err != nil {
		m.logger.Err(err).Msg("new server stream cannot be created")
	}

	err = m.handler.HandleRPC(dServerStream, rpc)
	if err != nil {
		sendErr := dServerStream.SendError(err)
		if sendErr != nil {
			m.logger.Err(sendErr).Msg("failure to send error to client")
		}
	}
}

func (m *ConnectionStream) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) (err error) {
	var metadata []byte
	if md, ok := drpcmetadata.Get(ctx); ok {
		metadata, err = drpcmetadata.Encode(metadata, md)
		if err != nil {
			return err
		}
	}

	stream, err := m.manager.NewClientStream(ctx)
	if err != nil {
		return err
	}

	// we have to protect m.writeBuff here even though the manager only allows one
	// stream at a time because the stream may async close allowing another
	// concurrent call to Invoke to proceed.
	m.mu.Lock()
	defer m.mu.Unlock()

	m.writeBuff, err = drpcenc.MarshalAppend(in, enc, m.writeBuff[:0])
	if err != nil {
		return err
	}

	if len(metadata) > 0 {
		err = stream.RawWrite(drpcwire.KindInvokeMetadata, metadata)
		if err != nil {
			return err
		}
	}

	err = stream.RawWrite(drpcwire.KindInvoke, []byte(rpc))
	if err != nil {
		return err
	}

	err = stream.RawWrite(drpcwire.KindMessage, m.writeBuff)
	if err != nil {
		return err
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	err = stream.MsgRecv(out, enc)
	if err != nil {
		return err
	}

	return nil
}

func (m *ConnectionStream) NewStream(ctx context.Context, rpc string, enc drpc.Encoding) (drpc.Stream, error) {
	var metadata []byte
	var err error
	if md, ok := drpcmetadata.Get(ctx); ok {
		metadata, err = drpcmetadata.Encode(metadata, md)
		if err != nil {
			return nil, err
		}
	}

	stream, err := m.manager.NewClientStream(ctx)
	if err != nil {
		return nil, err
	}

	if len(metadata) > 0 {
		if err := stream.RawWrite(drpcwire.KindInvokeMetadata, metadata); err != nil {
			return nil, err
		}
	}

	if err := stream.RawWrite(drpcwire.KindInvoke, []byte(rpc)); err != nil {
		return nil, err
	}

	if err := stream.RawFlush(); err != nil {
		return nil, err
	}

	return stream, nil
}

func (m *ConnectionStream) Close() error {
	return m.manager.Close()
}

func (m *ConnectionStream) Closed() <-chan struct{} {
	return m.manager.Closed()
}
