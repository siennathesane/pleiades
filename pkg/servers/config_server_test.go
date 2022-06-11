/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package servers

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/test/bufconn"
	"r3t.io/pleiades/pkg/blaze"
	"r3t.io/pleiades/pkg/conf"
	"r3t.io/pleiades/pkg/pb"
	"r3t.io/pleiades/pkg/servers/services"
	"r3t.io/pleiades/pkg/utils"
)

func TestConfigServerSuite(t *testing.T) {
	suite.Run(t, new(ConfigServerTests))
}

type ConfigServerTests struct {
	suite.Suite
	bufSize        int // const bufSize = 1024 * 1024
	bufferListener *bufconn.Listener
	env            *conf.EnvironmentConfig
	lifecycle      *fxtest.Lifecycle
	client         *api.Client
	store          *services.StoreManager
	logger         zerolog.Logger
	mux            *blaze.Router
}

func (cst *ConfigServerTests) SetupSuite() {
	cst.logger = utils.NewTestLogger(cst.T())

	var err error
	cst.lifecycle = fxtest.NewLifecycle(cst.T())
	cst.client, err = conf.NewConsulClient(cst.lifecycle)
	cst.Require().NoError(err, "failed to connect to consul")
	cst.Require().NotNil(cst.client, "the consul client can't be empty")

	cst.env, err = conf.NewEnvironmentConfig(cst.client)
	cst.Require().NoError(err, "the environment config is needed")
	cst.Require().NotNil(cst.env, "the environment config must be rendered")

	cst.store = services.NewStoreManager(cst.env, cst.logger, cst.client)
	cst.Require().NotNil(cst.store, "the store manager cannot be nil")

	err = cst.store.Start(false)
	cst.Require().NoError(err, "there must not be an error starting the store")

	cst.Require().NotPanics(func() {
		cst.mux = blaze.NewRouter()
	}, "there must not be a panic when building a new muxer")
	cst.Require().NoError(err, "there must not be an error when creating a new muxer")
	cst.Require().NotNil(cst.mux, "the muxer must not be nil")
}

func (cst *ConfigServerTests) TestNewConfigServer() {

	configServer := NewConfigServiceServer(cst.store, cst.logger)
	cst.Require().NotNil(configServer, "the config server must not be nil")

	err := pb.DRPCRegisterConfigService(cst.mux, configServer)
	cst.Require().NoError(err, "there must not be an error when registering the the config service")
}

func (cst *ConfigServerTests) TestConfigServerRaftConfigs() {
	// build the testkit
	testKit := blaze.NewTestKit(cst.T())

	// build the config service impl
	configServer := NewConfigServiceServer(cst.store, cst.logger)
	cst.Require().NotNil(configServer, "the config server must not be nil")

	// register it
	err := pb.DRPCRegisterConfigService(cst.mux, configServer)
	cst.Require().NoError(err, "there must not be an error when registering the the config service")

	// generate the test server
	testKit.NewServer(&blaze.TestKitServerArgs{AutoStart: true, Muxer: cst.mux})

	// generate a new connection stream
	configServiceTransportStream := testKit.NewConnectionStream()
	configServiceStream := blaze.NewConnectionStream(configServiceTransportStream, cst.mux, cst.logger)

	// build a client
	client := pb.NewDRPCConfigServiceClient(configServiceStream)
	cst.Require().NotNil(client, "the config server client must not be nil")

	// top level context
	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)

	// build a new request which fetches nothing
	requestOne := &pb.ConfigRequest{
		What:   pb.ConfigRequest_RAFT,
		Amount: pb.ConfigRequest_ONE,
	}
	respOne, err := client.GetConfig(ctx, requestOne)

	// verify nothing is there
	cst.Assert().Error(err, "there should be an error when requesting a specific record without specifying a record key")
	cst.Assert().Nil(respOne, "the response should be nil because there is an error")

	// build and marshal a generic config
	testStruct := &pb.RaftConfig{ClusterId: 123}
	payload, err := testStruct.MarshalVT()
	cst.Require().NoError(err, "there must not be an error serializing the test record")

	// manually store the test config in the underlying store
	testStorageKey := "request-one-test"
	err = cst.store.Put(testStorageKey, payload, reflect.TypeOf(&pb.RaftConfig{}))
	cst.Require().NoError(err, "there must not be an error storing a record for the test")

	// build the request to get the storage key we just manually stored
	requestTwo := &pb.ConfigRequest{
		What:   pb.ConfigRequest_RAFT,
		Amount: pb.ConfigRequest_ONE,
		Key:    &testStorageKey,
	}

	// fetch the storage key
	respTwo, err := client.GetConfig(ctx, requestTwo)
	cst.Require().NoError(err, "fetching a named record must not throw an error")
	cst.Require().NotNil(respTwo, "the named record mustn't be nil")

	// verify it's the proper type
	switch r := respTwo.Type.(type) {
	case *pb.ConfigResponse_RaftConfig:
		cst.Require().Equal(testStruct, r.RaftConfig.Configuration, "the fetched raft config should be equal")
	}

	// close the connection
	err = client.DRPCConn().Close()
	cst.Require().NoError(err, "there must not be an error when closing the client")

	// close the underlying listener
	testKit.CloseListener()
}
