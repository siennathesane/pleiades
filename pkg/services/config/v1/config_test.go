/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package v1

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/utils"
)

func TestConfigService(t *testing.T) {
	suite.Run(t, new(ConfigServiceTests))
}

type ConfigServiceTests struct {
	suite.Suite
	logger zerolog.Logger
}

func (cst *ConfigServiceTests) SetupTest() {
	cst.logger = utils.NewTestLogger(cst.T())
}

func (cst *ConfigServiceTests) TestConfigService() {
	cst.T().Skip("TODO")
}
