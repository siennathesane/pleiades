package main

import (
	"fmt"

	"go.uber.org/fx"
	conf2 "r3t.io/pleiades/pkg/conf"
)

const (
	defaultClusterId uint64 = 1
)

func main() {
	fmt.Println("hello from boulder.")

	app := fx.New(conf2.ProvideConsulClient(), conf2.ProvideEnvironmentConfig(), conf2.ProvideLogger())
	app.Run()

	//consulClient, err := etcd.NewClient(etcd.DefaultConfig())
	//if err != nil {
	//	panic(fmt.Errorf("cannot reach consul: %s", err))
	//}
	//
	//env, err := conf.NewEnvironmentConfig(consulClient)
	//if err != nil {
	//	panic(fmt.Errorf("cannot load environment config: %s", err))
	//}
	//
	//logger := conf.NewLogger(env.Environment)
	//
	//dirManager := services.NewDirectoryManager(env, logger.LoggerFactory("directory-manager"), consulClient)
	//localClusterDirectory, err := dirManager.NewDirectory(services.WriteAheadLogDirectory, 1)
	//if err != nil {
	//	logger.Errorf("can't create wal directory: %s", err)
	//	panic(err)
	//}
	//
	//nodeConfig := config.NodeHostConfig{WALDir: dataDir, NodeHostDir: dataDir, RTTMillisecond: 200, RaftAddress: "localhost:6000"}
	//
	//logger.Infof("configuring local raft node")
	//conf, err := conf.NewRaftConfig(consulClient, logger.LoggerFactory("configuration"))
	//if err != nil {
	//	logger.Panicf("error fetching raft configuration: %s", err)
	//}
	//
	//host, err := dragonboat.NewNodeHost(nodeConfig)
	//if err != nil {
	//	panic(err)
	//}
}
