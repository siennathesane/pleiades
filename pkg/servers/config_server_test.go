/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package servers

import (
	"context"
	"encoding/json"
	"net"
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"r3t.io/pleiades/pkg/conf"
	"r3t.io/pleiades/pkg/services"
	"r3t.io/pleiades/pkg/types"
)

type ConfigServerTests struct {
	suite.Suite
	bufSize        int // const bufSize = 1024 * 1024
	bufferListener *bufconn.Listener
	env            *conf.EnvironmentConfig
	lifecycle      *fxtest.Lifecycle
	client         *api.Client
	store          *services.StoreManager
	logger         dlog.ILogger
}

func TestConfigServerSuite(t *testing.T) {
	suite.Run(t, new(ConfigServerTests))
}

func (c *ConfigServerTests) SetupSuite() {
	c.logger = &conf.MockLogger{}

	var err error
	c.lifecycle = fxtest.NewLifecycle(c.T())
	c.client, err = conf.NewConsulClient(c.lifecycle)
	require.Nil(c.T(), err, "failed to connect to consul")
	require.NotNil(c.T(), c.client, "the consul client can't be empty")

	c.env, err = conf.NewEnvironmentConfig(c.client)
	require.Nil(c.T(), err, "the environment config is needed")
	require.NotNil(c.T(), c.env, "the environment config must be rendered")

	c.store = services.NewStoreManager(c.env, c.logger, c.client)
	require.NotNil(c.T(), c.store, "the store manager cannot be nil")
	err = c.store.Start(false)
	require.Nil(c.T(), err, "there must not be an error starting the store")

	c.bufSize = 1024 * 1024 // set the local buffer
}

func (c *ConfigServerTests) TestNewConfigServer() {
	c.bufferListener = bufconn.Listen(c.bufSize)
	s := grpc.NewServer()

	configServer := NewConfigServer(c.store, c.logger)
	require.NotNil(c.T(), configServer, "the config server must not be nil")

	RegisterConfigServiceServer(s, configServer)
	go func() {
		if err := s.Serve(c.bufferListener); err != nil {
			c.T().Fatalf("server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return c.bufferListener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.Nil(c.T(), err, "there must be no error on the bufnet dialer")
	defer func(conn *grpc.ClientConn, t *testing.T) {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}(conn, c.T())
}

func (c *ConfigServerTests) TestConfigServerRaftConfigs() {
	c.bufferListener = bufconn.Listen(c.bufSize)
	s := grpc.NewServer()

	configServer := NewConfigServer(c.store, c.logger)
	require.NotNil(c.T(), configServer, "the config server must not be nil")

	RegisterConfigServiceServer(s, configServer)
	go func() {
		if err := s.Serve(c.bufferListener); err != nil {
			c.T().Fatalf("server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return c.bufferListener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.Nil(c.T(), err, "there must be no error on the bufnet dialer")
	defer func(conn *grpc.ClientConn, t *testing.T) {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}(conn, c.T())

	client := NewConfigServiceClient(conn)
	require.NotNil(c.T(), client, "the config server client must not be nil")

	requestOne := &types.ConfigRequest{What: types.ConfigRequest_RAFT, Amount: types.ConfigRequest_ONE}
	respOne, err := client.GetConfig(ctx, requestOne)
	assert.NotNil(c.T(), err, "there should be an error when requesting a specific record without specifying a record key")
	assert.Nil(c.T(), respOne, "the response should be nil because there is an error")

	testStruct := &types.RaftConfig{ClusterId: 123}
	payload, err := json.Marshal(testStruct)
	require.Nil(c.T(), err, "there shouldn't be an error serializing the test record")

	err = c.store.Put("request-one-test", payload, reflect.TypeOf(&types.RaftConfig{}))
	require.Nil(c.T(), err, "there shouldn't be an error storing a record for the test")

	storageKey := "request-one-test"
	requestTwo := &types.ConfigRequest{
		What:   types.ConfigRequest_RAFT,
		Amount: types.ConfigRequest_ONE,
		Key:    &storageKey,
	}

	respTwo, err := client.GetConfig(ctx, requestTwo)
	require.Nil(c.T(), err, "fetching a named record mustn't throw an error")
	require.NotNil(c.T(), respTwo, "the named record mustn't be nil")
	switch r := respTwo.Type.(type) {
	case *types.ConfigResponse_RaftConfig:
		assert.Equal(c.T(), testStruct, r.RaftConfig.Configuration, "the fetched raft config should be equal")
	}
}
