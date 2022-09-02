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
	"github.com/cockroachdb/errors"
)

type StateMachineType uint64

const (
	testStateMachineType StateMachineType = 0
	BBoltStateMachineType StateMachineType = 1
)

var (
	ErrUnsupportedStateMachine = errors.New("state machine type is unsupported")
)