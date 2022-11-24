/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

import (
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/execution"
)

type Runtime interface {
	Init() error
	CreateExecutor(*Request) (execution.IExecutor, error)
}
