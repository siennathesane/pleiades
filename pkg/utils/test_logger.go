/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package utils

import (
	"testing"

	"github.com/rs/zerolog"
)

func NewTestLogger(t *testing.T) zerolog.Logger {
	testWriter := zerolog.NewTestWriter(t)
	return zerolog.New(testWriter)
}
