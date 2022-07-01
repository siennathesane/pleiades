
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package config

import (
	"testing"

	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/protocols/v1/config"
	"r3t.io/pleiades/pkg/utils"
)

func TestRegistry(t *testing.T) {
	suite.Run(t, new(RegistryTests))
}

type RegistryTests struct {
	suite.Suite
	logger   zerolog.Logger
	registry *Registry
}

// implement BeforeTest interface
func (s *RegistryTests) SetupTest() {
	s.logger = utils.NewTestLogger(s.T())

	var err error
	s.registry, err = NewRegistry(s.logger)
	s.Require().NoError(err, "there must not be an error creating the Registry")
}

func (s *RegistryTests) TestGet() {
	srv := &server.Server{}
	err := s.registry.PutServer(config.ServiceType_Type_test, srv)
	s.Require().NoError(err, "there must not be an error putting the value")

	value, _ := s.registry.GetServer(config.ServiceType_Type_test)
	s.Assert().Equal(srv, value)
}

func (s *RegistryTests) TestPut() {
	srv := &server.Server{}
	err := s.registry.PutServer(config.ServiceType_Type_test, srv)
	s.Require().NoError(err, "there must not be an error putting the value")

	value, _ := s.registry.GetServer(config.ServiceType_Type_test)
	s.Assert().Equal(srv, value)
}
