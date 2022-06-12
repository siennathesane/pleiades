package conf

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewConsulClient(t *testing.T) {
	// verify consul is available
	consulHost := os.Getenv("CONSUL_HTTP_ADDR")
	require.NotEmpty(t, consulHost, "$CONSUL_HTTP_ADDR must be set")

	urls, err := url.Parse(consulHost)
	require.NoError(t, err, "there must not be an error when parsing $CONSUL_HTTP_ADDR")

	assert.NoError(t, portChecker(urls.Host, "443"), "verify consul is available")

	client, err := NewConsulClient(fxtest.NewLifecycle(t))
	if err != nil {
		require.FailNow(t, "failed to connect to consul")
	}

	require.NotNil(t, client, "consul etcd instance is returning nil")
}
