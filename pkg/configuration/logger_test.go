/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (l *LoggerTestSuite) TestNewLogger() {
	//logger, err := NewLogger(zerolog.NewTestWriter(l.T()))
	//assert.Nil(l.T(), err, "the logger can't be built")
	//assert.NotNil(l.T(), logger, "the sugaredLogger can't be nil")
	//
	//assert.NotPanics(l.T(), func() {
	//	logger.SetLevel(dlog.DEBUG)
	//})
	//
	//assert.NotPanics(l.T(), func() {
	//	logger.Debugf("testing debugf functionality")
	//})
	//
	//for i := 0; i < 3; i++ {
	//	logger.Debugf("test debug iteration %d", i)
	//}

	//assert.Equal(l.T(), dlog.DEBUG, logger.GetLevel(), "the log level should be DEBUG")
}
