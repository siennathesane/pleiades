/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"bytes"
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	configv1 "r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/services/v1/config"
	"r3t.io/pleiades/pkg/utils"
)

// StreamReceiverTest tests the StreamReceiver class.
func TestStreamReceiver(t *testing.T) {
	suite.Run(t, new(StreamReceiverTest))
}

type StreamReceiverTest struct {
	suite.Suite
	tk       *QuicTestKit
	logger   zerolog.Logger
	registry *config.Registry
}

func (s *StreamReceiverTest) SetupSuite() {
	s.logger = utils.NewTestLogger(s.T())
	s.tk = NewQuicTestKit(s.T())

	var err error
	s.registry, err = config.NewRegistry(s.logger)
	s.Require().NoError(err, "failed to create registry")

	err = s.registry.PutServer(
		configv1.ServiceType_Type_test,
		&server.Server{},
	)
	s.Require().NoError(err, "failed to register test server")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiverServerRouter() {
	sr, err := NewStreamReceiver(s.logger, s.registry)
	s.Require().NoError(err, "there must not be an error when creating a stream receiver")

	msg,seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	s.Require().NoError(err, "there must not be an error when creating a message")

	svcType, err := configv1.NewRootServiceType(seg)
	s.Require().NoError(err, "there must not be an error when creating a service type")
	svcType.SetType(configv1.ServiceType_Type_test)

	var buf bytes.Buffer
	err = capnp.NewEncoder(&buf).Encode(msg)
	s.Require().NoError(err, "there must not be an error when creating a capnp message")

	testStream := s.tk.NewConnectionStream()

	go func() {
		time.Sleep(time.Second * 1)
		n, err := testStream.Write(buf.Bytes())
		s.Require().NoError(err, "there must not be an error when writing to a stream")
		s.Require().Equal(buf.Len(), n, "the number of bytes written must match the number of bytes read")
	}()

	conn, err := s.tk.GetListener().Accept(context.Background())
	s.Require().NoError(err, "there must not be an error when accepting a connection")
	s.Require().NotNil(conn, "the connection must not be nil")

	receivingTestStream, err := conn.AcceptStream(context.Background())
	s.Require().NoError(err, "there must not be an error when accepting a stream")
	s.Require().NotNil(receivingTestStream, "the stream must not be nil")

	sr.Receive(context.Background(), s.logger, receivingTestStream)
}
