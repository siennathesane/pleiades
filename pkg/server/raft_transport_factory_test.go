/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	"testing"

	"github.com/mxplusb/pleiades/pkg/utils"
	dconfig "github.com/lni/dragonboat/v3/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestRaftTransportFactory(t *testing.T) {
	suite.Run(t, new(RaftTransportFactoryTestSuite))
}

type RaftTransportFactoryTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func (r *RaftTransportFactoryTestSuite) SetupSuite() {
	r.logger = utils.NewTestLogger(r.T())
}

func (r *RaftTransportFactoryTestSuite) TestNewRaftTransportFactory() {
	firstTestHost := randomLibp2pTestHost()

	rtf := NewRaftTransportFactory(firstTestHost, r.logger)
	r.Require().NotNil(rtf, "raft transport factory must not be nil")
}

func (r *RaftTransportFactoryTestSuite) TestNewRaftTransport() {
	firstTestHost := randomLibp2pTestHost()

	rtf := NewRaftTransportFactory(firstTestHost, r.logger)
	r.Require().NotNil(rtf, "raft transport factory must not be nil")

	transport := rtf.Create(dconfig.NodeHostConfig{}, nil, nil)
	r.Require().NotNil(transport, "raft transport must not be nil")
	r.Require().NotNil(transport.(*raftTransport).host, "raft transport host must not be nil")
}

func (r *RaftTransportFactoryTestSuite) TestValidate() {
	firstTestHost := randomLibp2pTestHost()

	rtf := NewRaftTransportFactory(firstTestHost, r.logger)
	r.Require().NotNil(rtf, "raft transport factory must not be nil")

	invalid := rtf.Validate("256.256.256.256:1234")
	r.Require().False(invalid, "invalid address must be detected")

	valid := rtf.Validate("/ip4/1.2.3.4/tcp/1234")
	r.Require().True(valid, "valid address must be detected")
}