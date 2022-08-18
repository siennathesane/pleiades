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
	"testing"

	"github.com/lni/dragonboat/v3"
	"github.com/stretchr/testify/suite"
)

func TestRaftControl(t *testing.T) {
	suite.Run(t, new(RaftControlTests))
}

type RaftControlTests struct {
	suite.Suite

	nh *dragonboat.NodeHost
}

func (rct *RaftControlTests) SetupSuite() {
	rct.nh = buildTestNodeHost(rct.T())
}

func (rct *RaftControlTests) TestGetId() {
	rct.T().Log("instantiated")
}
