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
	"hash/crc32"
	"io"
	"math/rand"
	"testing"
	"time"

	transportv1 "github.com/mxplusb/pleiades/pkg/api/v1"
	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftControlRPCServer(t *testing.T) {
	suite.Run(t, new(RaftControlRPCServerTests))
}

type RaftControlRPCServerTests struct {
	suite.Suite
	rpcHost host.Host
	logger  zerolog.Logger
}

func (r *RaftControlRPCServerTests) SetupSuite() {
	r.logger = conf.NewRootLogger()

	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)
	hostAddr := fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic", port)

	ma, err := multiaddr.NewMultiaddr(hostAddr)
	r.Require().NoError(err, "there must not be an error when parsing the test multiaddr")

	r.rpcHost, err = libp2p.New(libp2p.ListenAddrs(ma))
	r.Require().NoError(err, "there must not be an error when creating the libp2p host")
}

func (r *RaftControlRPCServerTests) TestNewRaftControlRPCServer() {
	nodeHostConfig := buildTestNodeHostConfig(r.T())
	node, err := NewRaftControlNode(nodeHostConfig, r.logger)
	r.Require().NoError(err, "there must not be an error when creating the control node")
	r.Require().NotNil(node, "the node must not be nil")

	testRpcServer := NewRaftControlRPCServer(node, r.rpcHost, r.logger)
	r.Require().NotNil(testRpcServer, "the raft control rpc host must not me nil")

	clientHost := randomTestHost()
	clientHost.Peerstore().AddAddrs(r.rpcHost.ID(), r.rpcHost.Addrs(), peerstore.PermanentAddrTTL)

	stream, err := clientHost.NewStream(context.Background(), r.rpcHost.ID(), RaftControlProtocolVersion)
	r.Require().NoError(err, "there must not be an error when opening a stream")
	r.Require().NotNil(stream, "the stream must not be nil")
}

func (r *RaftControlRPCServerTests) TestGetId() {
	nodeHostConfig := buildTestNodeHostConfig(r.T())
	node, err := NewRaftControlNode(nodeHostConfig, r.logger)
	r.Require().NoError(err, "there must not be an error when creating the control node")
	r.Require().NotNil(node, "the node must not be nil")

	testRpcServer := NewRaftControlRPCServer(node, r.rpcHost, r.logger)
	r.Require().NotNil(testRpcServer, "the raft control rpc host must not me nil")

	clientHost := randomTestHost()
	clientHost.Peerstore().AddAddrs(r.rpcHost.ID(), r.rpcHost.Addrs(), peerstore.PermanentAddrTTL)

	stream, err := clientHost.NewStream(context.Background(), r.rpcHost.ID(), RaftControlProtocolVersion)
	r.Require().NoError(err, "there must not be an error when opening a stream")
	r.Require().NotNil(stream, "the stream must not be nil")

	initialStateMsg := &transportv1.State{State: 0, HeaderToFollow: 1}
	buf, _ := initialStateMsg.MarshalVT()

	_, err = stream.Write(buf)
	r.Require().NoError(err, "there must not be an error when writing the initial state to the stream")

	idReqPayload := &database.RaftControlPayload{
		Method: database.RaftControlPayload_GET_ID,
		Types: &database.RaftControlPayload_IdRequest{
			IdRequest: &database.IdRequest{},
		},
	}
	payloadBuf, _ := idReqPayload.MarshalVT()

	initialHeader := &transportv1.Header{
		Size: uint32(len(payloadBuf)),
		Checksum: crc32.ChecksumIEEE(payloadBuf),
	}
	headerBuf, _ := initialHeader.MarshalVT()

	_, err = stream.Write(headerBuf)
	r.Require().NoError(err, "there must not be an error when writing the header to the stream")

	err = VerifyStreamState(stream)
	r.Require().NoError(err, "there must not be an error when verifying the stream state after sending the header")

	_, err = stream.Write(payloadBuf)
	r.Require().NoError(err, "there must not be an error when writing the payload to the stream")

	err = VerifyStreamState(stream)
	r.Require().NoError(err, "there must not be an error when verifying the stream state after sending the payload")

	responseHeaderBuf := make([]byte, headerSize)
	bytesRead, err := io.ReadFull(stream, responseHeaderBuf)
	r.Require().NoError(err, "there must not be an error when reading the response header")
	r.Require().Equal(len(headerBuf), bytesRead, "the response header must be 10 bytes")

	responseHeader := &transportv1.Header{}
	err = responseHeader.UnmarshalVT(responseHeaderBuf)
	r.Require().NoError(err, "there must not be an error when unmarshalling the header")
	r.Require().NotEmpty(responseHeader.Size, "the length of the response must not be 0")
	r.Require().NotEmpty(responseHeader.Checksum, "the checksum must not be 0")

	msgBuf := make([]byte, responseHeader.Size)
	msgBytes, err := io.ReadFull(stream, msgBuf)
	r.Require().NoError(err, "there must not be an error when reading the response payload")
	r.Require().Equal(int(responseHeader.Size), msgBytes, "the length of the payload must equal the response header size value")

	responsePayload := &database.RaftControlPayload{}
	err = responsePayload.UnmarshalVT(msgBuf)
	r.Require().NoError(err, "there must not be an error when unmarshalling the response payload")

	switch resp := responsePayload.Types.(type) {
	case *database.RaftControlPayload_Error:
		r.Require().Nil(resp.Error.Message, "there should not be an error fetching the raft control node id")
	case *database.RaftControlPayload_IdResponse:
		r.Require().NotEmpty(resp.IdResponse.Id, "the id of the node must not be nil")
	}
}
