package blaze

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/services/v1/config"
	"r3t.io/pleiades/pkg/utils"
)

func TestStreamManager(t *testing.T) {
	suite.Run(t, new(StreamServerTests))
}

type StreamServerTests struct {
	suite.Suite
	logger   zerolog.Logger
	qtk      *QuicTestKit
	registry *config.Registry
}

func (smt *StreamServerTests) SetupSuite() {
	smt.logger = utils.NewTestLogger(smt.T())
	smt.qtk = NewQuicTestKit(smt.T())
	smt.registry, _ = config.NewRegistry(smt.logger)
}

func (smt *StreamServerTests) BeforeTest(suiteName, testName string) {
	smt.qtk.Start()
}

//func (smt *StreamServerTests) AfterTest(suiteName, testName string) {
//	smt.qtk.Stop()
//}

func (smt *StreamServerTests) TestHandleConnection() {
	testServer := NewServer(smt.qtk.listener, smt.logger, smt.registry)
	smt.Require().NotNil(testServer, "the server must not be nil")

	ctx := context.Background()
	err := testServer.Start(ctx)
	smt.Require().NoError(err, "there must not be an error starting the test server")

	conn := smt.qtk.GetConnection()

	stream, err := conn.OpenStream()
	smt.Require().NoError(err, "there must not be an error opening a stream")

	// this will fail
	n, err := stream.Write([]byte("hello"))
	smt.Require().Error(err, "there must be an error writing to the stream")
	smt.Require().Equal(5, n, "the number of bytes written must be 5")
}
