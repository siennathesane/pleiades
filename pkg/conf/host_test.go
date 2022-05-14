package conf

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
)

type HostConfigTests struct {
	suite.Suite
	client *api.Client
}

func TestHostConfiguration(t *testing.T) {
	suite.Run(t, new(HostConfigTests))
}

func (h *HostConfigTests) SetupSuite() {
	var err error
	h.client, err = NewConsulClient(fxtest.NewLifecycle(h.T()))
	require.Nil(h.T(), err, "failed to connect to consul")
}

func (h *HostConfigTests) TestNewEnvironmentConfigLoad() {
	config, err := NewEnvironmentConfig(h.client)
	assert.Nil(h.T(), err, "error reading configuration")
	assert.NotEmpty(h.T(), config.Environment, "configuration environment cannot be empty")
	assert.NotEmpty(h.T(), config.GCPProjectId, "the gcp project id cannot be empty")
	assert.NotEmpty(h.T(), config.BaseDir, "the base directory must be set")
	assert.NotEmpty(h.T(), config.BasePort, "the base port must be set")
	assert.NotEmpty(h.T(), config.MaxPort, "the max port must be set")
	assert.NotEmpty(h.T(), config.LocalClusterId, "the local cluster id must be set")
	assert.NotEmpty(h.T(), config.LocalExchangeId, "the local exchange id must be set")
	assert.NotEmpty(h.T(), config.Hostname, "the hostname must be set")
}

func (h *HostConfigTests) TestEnvironmentConfigVerification() {
	validateHostConfig = false
	config, err := NewEnvironmentConfig(h.client)
	assert.Nil(h.T(), err, "error reading configuration")

	priorPort := config.BasePort
	config.BasePort = 150
	assert.Error(h.T(), config.validate(), "the port validation is shouldn't accept ports lower than 1024")
	config.BasePort = priorPort // reset so it doesn't throw another error

	config.MaxPort = config.BasePort + 10
	assert.Error(h.T(), config.validate(), "the port range should be at least 4000")
}
