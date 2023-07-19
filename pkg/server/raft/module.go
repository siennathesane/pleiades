/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package raft

import (
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"go.uber.org/fx"
)

var RaftModule = fx.Module("raftpb",
	fx.Provide(NewHost),
	fx.Provide(runtime.AsRoute(NewRaftHostConnectAdapter)),
)
