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

	"github.com/mxplusb/pleiades/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
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

	clientHost := randomLibp2pTestHost()
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

	clientHost := randomLibp2pTestHost()
	clientHost.Peerstore().AddAddrs(r.rpcHost.ID(), r.rpcHost.Addrs(), peerstore.PermanentAddrTTL)

	stream, err := clientHost.NewStream(context.Background(), r.rpcHost.ID(), RaftControlProtocolVersion)
	r.Require().NoError(err, "there must not be an error when opening a stream")
	r.Require().NotNil(stream, "the stream must not be nil")

	idReqPayload := &database.RaftControlPayload{
		Method: database.RaftControlPayload_GET_ID,
		Types: &database.RaftControlPayload_IdRequest{
			IdRequest: &database.IdRequest{},
		},
	}
	payloadBuf, _ := proto.Marshal(idReqPayload)

	outgoingFrame := NewFrame().WithService(RaftControlServiceByte).WithMethod(GetId).WithPayload(payloadBuf)
	frameBuf, err := outgoingFrame.Marshal()
	r.Require().NoError(err, "there must not be an error when marshalling the outgoingFrame")
	r.Require().NotNil(frameBuf, "the outgoingFrame buffer must not be nil")

	_, err = stream.Write(frameBuf)
	r.Require().NoError(err, "there must not be an error when writing the payload to the stream")

	incomingFrame := NewFrame()
	read, err := incomingFrame.ReadFrom(stream)
	r.Require().NoError(err, "there must not be an error when reading the frame from the stream")
	r.Require().NotEmpty(read, "there must be bytes read")
	r.Require().NotNil(incomingFrame.GetPayload(), "the incoming frame payload must not be nil")

	respPayload := incomingFrame.GetPayload()

	idResponse := &database.IdResponse{}
	err = proto.Unmarshal(respPayload, idResponse)
	r.Require().NoError(err, "there must not be an error when unmarshalling the response payload")
	r.Require().Equal(node.nh.ID(), idResponse.GetId(), "the response ID value must equal the node value")
}
