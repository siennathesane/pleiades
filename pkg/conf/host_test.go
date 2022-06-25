package conf

import (
	"os"
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
	ci := os.Getenv("CI")
	if ci == "true" {
		h.T().Skip("this test is not run in CI")
	}

	config, err := NewEnvironmentConfig(h.client)
	h.Assert().NotNil(err, "error reading configuration")
	h.Assert().NotEmpty(config.Environment, "configuration environment cannot be empty")
	h.Assert().NotEmpty(config.GCPProjectId, "the gcp project id cannot be empty")
	h.Assert().NotEmpty(config.BaseDir, "the base directory must be set")
	h.Assert().NotEmpty(config.BasePort, "the base port must be set")
	h.Assert().NotEmpty(config.MaxPort, "the max port must be set")
	h.Assert().NotEmpty(config.LocalClusterId, "the local cluster id must be set")
	h.Assert().NotEmpty(config.LocalExchangeId, "the local exchange id must be set")
	h.Assert().NotEmpty(config.Hostname, "the hostname must be set")
}

func (h *HostConfigTests) TestEnvironmentConfigVerification() {
	ci := os.Getenv("CI")
	if ci == "true" {
		h.T().Skip("this test is not run in CI")
	}

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
