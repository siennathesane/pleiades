/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package kvstore

import (
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"go.uber.org/fx"
)

var KvStoreModule = fx.Module("kv",
	fx.Provide(runtime.AsRoute(NewKvStoreBboltConnectAdapter)),
	fx.Provide(runtime.AsRoute(NewKvstoreTransactionConnectAdapter)),
	fx.Provide(NewBboltStoreManager),
)
