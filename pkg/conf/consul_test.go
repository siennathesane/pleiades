package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewConsulClient(t *testing.T) {
	// verify consul is available
	assert.Nil(t, portChecker("localhost", "8500"), "verify consul is available. make sure it's running on your local machine and with the proper configs")

	client, err := NewConsulClient(fxtest.NewLifecycle(t))
	if err != nil {
		require.FailNow(t, "failed to connect to consul")
	}

	require.NotNil(t, client, "consul api instance is returning nil")
}
