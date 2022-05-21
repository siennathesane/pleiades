package servers

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	dlog "github.com/lni/dragonboat/v3/logger"
	"r3t.io/pleiades/pkg/services"
	"r3t.io/pleiades/pkg/types"
)

var _ ConfigServiceServer = ConfigServer{}

type ConfigServer struct {
	UnimplementedConfigServiceServer
	manager *services.StoreManager
	logger  *dlog.ILogger
}

func NewConfigServer(manager *services.StoreManager, logger *dlog.ILogger) *ConfigServer {
	return &ConfigServer{manager: manager, logger: logger}
}

func (c ConfigServer) GetConfig(ctx context.Context, config *types.ConfigRequest) (*types.ConfigResponse, error) {
	switch config.What {
	case types.ConfigRequest_ALL:
	case types.ConfigRequest_RAFT:
		switch config.Amount {
		case types.ConfigRequest_ONE:
			return c.getRaftConfig(config.Name)
		case types.ConfigRequest_EVERYTHING:
			return c.getAllRaftConfigs()
		}
	}
	return nil, errors.New("cannot determine which type of config to return")
}

func (c ConfigServer) getRaftConfig(name *string) (*types.ConfigResponse, error) {
	if name == nil {
		return nil, errors.New("cannot request a named record without a key")
	}

	val, err := c.manager.Get(*name, reflect.TypeOf(&types.RaftConfig{}))
	if err != nil {
		return nil, err
	}

	var config types.RaftConfig
	if err := json.Unmarshal(val, &config); err != nil {
		return nil, err
	}

	t := &types.ConfigResponse{
		Type: &types.ConfigResponse_RaftConfig{
			RaftConfig: &types.GetRaftConfigResponse{
				Configuration: &config,
			},
		},
	}

	return t, nil
}

func (c ConfigServer) getAllRaftConfigs() (*types.ConfigResponse, error) {
	resp, err := c.manager.GetAll(reflect.TypeOf(&types.RaftConfig{}))
	if err != nil {
		return nil, err
	}

	var configs map[string]*types.RaftConfig
	for k, _ := range resp {
		var val *types.RaftConfig
		if err := json.Unmarshal(resp[k], &val); err != nil {
			return nil, err
		}
		configs[k] = val
	}

	return &types.ConfigResponse{
		Type: &types.ConfigResponse_AllRaftConfigs{
			AllRaftConfigs: &types.ListRaftConfigsResponse{
				AvailableConfigs: configs,
			},
		},
	}, nil
}
