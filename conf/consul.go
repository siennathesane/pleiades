package conf

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/fx"
)

// ProvideConsulClient is the main constructor for providing the Consul agent client
func ProvideConsulClient() fx.Option {
	return fx.Provide(NewConsulClient)
}

// NewConsulClient creates a new client for Consul
func NewConsulClient(lifecycle fx.Lifecycle) (*api.Client, error) {
	return api.NewClient(api.DefaultConfig())
}
