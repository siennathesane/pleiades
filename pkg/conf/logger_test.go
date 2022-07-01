
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package conf

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
)

type LoggerTestSuite struct {
	suite.Suite
	lifecycle *fxtest.Lifecycle
	client    *api.Client
	env       *EnvironmentConfig
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (l *LoggerTestSuite) SetupSuite() {
	var err error
	l.lifecycle = fxtest.NewLifecycle(l.T())
	l.client, err = NewConsulClient(l.lifecycle)
	require.Nil(l.T(), err, "failed to connect to consul")
	require.NotNil(l.T(), l.client, "the consul client can't be empty")

	l.env, err = NewEnvironmentConfig(l.client)
	require.Nil(l.T(), err, "the environment config is needed")
	require.NotNil(l.T(), l.env, "the environment config must be rendered")
}

func (l *LoggerTestSuite) TestNewLogger() {
	ci := os.Getenv("CI")
	if ci == "true" {
		l.T().Skip("this test is not run in CI")
	}

	logger, err := NewLogger(l.lifecycle, l.client, l.env)
	assert.Nil(l.T(), err, "the logger can't be built")
	assert.NotNil(l.T(), logger, "the sugaredLogger can't be nil")

	assert.NotPanics(l.T(), func() {
		logger.SetLevel(dlog.DEBUG)
	})
	assert.Equal(l.T(), dlog.DEBUG, logger.GetLevel(), "the log level should be DEBUG")

	assert.NotPanics(l.T(), func() {
		logger.Debugf("testing debugf functionality")
	})

	for i := 0; i < 3; i++ {
		logger.Debugf("test debug iteration %d", i)
	}
}
