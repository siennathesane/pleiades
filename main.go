package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
)

const (
	defaultClusterId uint64 = 1
)

func main() {
	fmt.Println("hello from boulder.")
	nodeId := flag.Int("node-id", 1, "node id")

	conf := config.Config{
		NodeID:             uint64(*nodeId),
		ClusterID:          defaultClusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    100,
		CompactionOverhead: 5}
	dataDir := filepath.Join("example-data", "hello-world", fmt.Sprintln("node-%d", *nodeId))

	nodeConfig := config.NodeHostConfig{WALDir: dataDir, NodeHostDir: dataDir, RTTMillisecond: 200, RaftAddress: "localhost:6000"}

	host, err := dragonboat.NewNodeHost(nodeConfig)
	if err != nil {
		panic(err)
	}
}
