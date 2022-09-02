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
