

/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

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

	// this will write but nothing is listening
	n, err := stream.Write([]byte("hello"))
	smt.Require().Equal(5, n, "the number of bytes written must be 5")
}
