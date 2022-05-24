package servers

import (
	"context"
	"errors"

	dlog "github.com/lni/dragonboat/v3/logger"
	"r3t.io/pleiades/pkg/managers"
	"r3t.io/pleiades/pkg/services"
	"r3t.io/pleiades/pkg/types"
)

var _ ConfigServiceServer = &ConfigServer{}

type ConfigServer struct {
	UnimplementedConfigServiceServer
	manager     *services.StoreManager
	logger      dlog.ILogger
	raftManager *managers.RaftManager[types.RaftConfig]
}

func NewConfigServer(manager *services.StoreManager, logger dlog.ILogger) *ConfigServer {
	return &ConfigServer{
		manager:     manager,
		logger:      logger,
		raftManager: managers.NewRaftManager(manager, logger)}
}

func (c *ConfigServer) GetConfig(ctx context.Context, config *types.ConfigRequest) (*types.ConfigResponse, error) {
	switch config.What {
	case types.ConfigRequest_ALL:
	case types.ConfigRequest_RAFT:
		switch config.Amount {
		case types.ConfigRequest_ONE:
			return c.getRaftConfig(config.Key)
		case types.ConfigRequest_EVERYTHING:
			return c.getAllRaftConfigs()
		}
	}
	return nil, errors.New("cannot determine which type of config to return")
}

func (c *ConfigServer) getRaftConfig(name *string) (*types.ConfigResponse, error) {
	if name == nil {
		return nil, errors.New("cannot request a named record without a key")
	}

	val, err := c.raftManager.Get(*name)
	if err != nil {
		return nil, err
	}

	t := &types.ConfigResponse{
		Type: &types.ConfigResponse_RaftConfig{
			RaftConfig: &types.GetRaftConfigResponse{
				Configuration: val,
			},
		},
	}

	return t, nil
}

func (c *ConfigServer) getAllRaftConfigs() (*types.ConfigResponse, error) {

	all, err := c.raftManager.GetAll()
	if err != nil {
		return nil, err
	}

	return &types.ConfigResponse{
		Type: &types.ConfigResponse_AllRaftConfigs{
			AllRaftConfigs: &types.ListRaftConfigsResponse{
				AvailableConfigs: all,
			},
		},
	}, nil
}
