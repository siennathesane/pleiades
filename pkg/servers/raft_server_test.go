package servers

import (
	"testing"

	"github.com/hashicorp/consul/api"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/test/bufconn"
	"r3t.io/pleiades/pkg/conf"
	"r3t.io/pleiades/pkg/services"
)

type RaftConfigServerTests struct {
	suite.Suite
	bufSize        int // const bufSize = 1024 * 1024
	bufferListener *bufconn.Listener
	env            *conf.EnvironmentConfig
	lifecycle      *fxtest.Lifecycle
	client         *api.Client
	store          *services.StoreManager
	logger         dlog.ILogger
}

func TestRaftConfigServer(t *testing.T) {
	suite.Run(t, new(RaftConfigServerTests))
}

func (rcs *RaftConfigServerTests) SetupSuite() {
	rcs.logger = conf.MockLogger{}

	var err error
	rcs.lifecycle = fxtest.NewLifecycle(rcs.T())
	rcs.client, err = conf.NewConsulClient(rcs.lifecycle)
	require.Nil(rcs.T(), err, "failed to connect to consul")
	require.NotNil(rcs.T(), rcs.client, "the consul client can't be empty")

	rcs.env, err = conf.NewEnvironmentConfig(rcs.client)
	require.Nil(rcs.T(), err, "the environment config is needed")
	require.NotNil(rcs.T(), rcs.env, "the environment config must be rendered")

	rcs.store = services.NewStoreManager(rcs.env, rcs.logger, rcs.client)
	require.NotNil(rcs.T(), rcs.store, "the store manager cannot be nil")
	err = rcs.store.Start(false)
	require.Nil(rcs.T(), err, "there must not be an error starting the store")

	rcs.bufSize = 1024 * 1024 // set the local buffer
}