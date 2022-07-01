/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package pkg

const (
	Version string = ""
	Sha string = ""
)

//goland:noinspection GoCommentStart
const (
	// the ranges a system cluster (internal to pleiades) can be
	SystemClusterIdRangeStartingValue uint64 = 1
	SystemClusterIdRangeEndingValue   uint64 = 1_000_000_000

	// the ranges an exchange cluster can be
	ExchangeClusterIdRangeStartingValue uint64 = 2_000_000_000
	ExchangeClusterIdRangeEndingValue   uint64 = 3_000_000_000

	// the ranges a customer fsm plugin cluster could be
	CustomerFsmPluginRangeStartingValue uint64 = 5_000_000_000
	CustomerFsmPluginRangeEndingValue   uint64 = 10_000_000_000

	MaxRaftClustersPerHost uint64 = 4096
)
