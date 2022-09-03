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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/raftpb"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftTransport(t *testing.T) {
	suite.Run(t, new(RaftTransportTestSuite))
}

type RaftTransportTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func (r *RaftTransportTestSuite) SetupSuite() {
	r.logger = utils.NewTestLogger(r.T())
}

func (r *RaftTransportTestSuite) TestNewRaftTransport() {
	testHost := randomLibp2pTestHost()

	rtf := NewRaftTransportFactory(testHost, r.logger)

	transport := rtf.Create(dconfig.NodeHostConfig{}, r.testMessageBatchHandler, r.testChunkHandler)
	r.Assert().NotNil(transport)
}

func (r *RaftTransportTestSuite) TestGetConnection() {
	testLocalHost := randomLibp2pTestHost()

	rtf := NewRaftTransportFactory(testLocalHost, r.logger)

	transport := rtf.Create(dconfig.NodeHostConfig{}, r.testMessageBatchHandler, r.testChunkHandler)

	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)
	validHostAddr := fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic", port)

	ma, _ := multiaddr.NewMultiaddr(validHostAddr)
	testRemoteHost, _ := libp2p.New(libp2p.ListenAddrs(ma))
	defer func(lhost host.Host) {
		err := lhost.Close()
		r.Require().NoError(err, "failed to close host")
	}(testRemoteHost)

	peerId := testRemoteHost.ID()
	r.Assert().NotNil(peerId)

	conn, err := transport.GetConnection(context.Background(), validHostAddr)
	r.Require().Nil(err, "there must not be an error when getting a new transport connection")
	r.Require().NotNil(conn ,"the transport connection must not be nil")


}

func (r *RaftTransportTestSuite) testMessageBatchHandler(batch raftpb.MessageBatch) {
	r.Assert().NotNil(batch)
}

func (r *RaftTransportTestSuite) testChunkHandler(chunk raftpb.Chunk) bool {
	r.Assert().NotNil(chunk)
	return false
}