package pkg

//goland:noinspection GoCommentStart
const (
	// the ranges a system cluster (internal to pleiades) can be
	SystemClusterIdRangeStartingValue uint64 = 1
	SystemClusterIdRangeEndingValue uint64 = 1_000_000_000

	// the ranges an exchange cluster can be
	ExchangeClusterIdRangeStartingValue uint64 = 2_000_000_000
	ExchangeClusterIdRangeEndingValue uint64 = 3_000_000_000

	// the ranges a customer fsm plugin cluster could be
	CustomerFsmPluginRangeStartingValue uint64 = 5_000_000_000
	CustomerFsmPluginRangeEndingValue uint64 = 10_000_000_000

	MaxRaftClustersPerHost uint64 = 4096
)
