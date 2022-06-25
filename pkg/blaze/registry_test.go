/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/utils"
)

func TestRegistry(t *testing.T) {
	suite.Run(t, new(RegistryTests))
}

type RegistryTests struct {
	suite.Suite
	logger zerolog.Logger
	registry *Registry
}

// implement BeforeTest interface
func (s *RegistryTests) SetupTest() {
	s.logger = utils.NewTestLogger(s.T())

	var err error
	s.registry,err = NewRegistry(s.logger)
	s.Require().NoError(err, "there must not be an error creating the Registry")
}

func (s *RegistryTests) TestGet() {
	err := s.registry.Put("key", []byte("value"))
	s.Require().NoError(err, "there must not be an error putting the value")

	value, _ := s.registry.Get("key")
	s.Assert().Equal([]byte("value"), value)
}

func (s *RegistryTests) TestPut() {
	err := s.registry.Put("key", []byte("value"))
	s.Require().NoError(err, "there must not be an error putting the value")

	value, _ := s.registry.Get("key")
	s.Assert().Equal([]byte("value"), value)
}
