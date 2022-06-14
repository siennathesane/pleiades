package servers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/fsm"
	configv1 "r3t.io/pleiades/pkg/pb/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
)

var _ configv1.DRPCConfigServiceServer = (*ConfigServiceServer)(nil)

type ConfigServiceServer struct {
	configv1.DRPCConfigServiceUnimplementedServer
	manager     *services.StoreManager
	logger      zerolog.Logger
	raftManager *fsm.RaftManager[configv1.RaftConfig]
}

func NewConfigServiceServer(manager *services.StoreManager, logger zerolog.Logger) *ConfigServiceServer {
	return &ConfigServiceServer{
		manager:     manager,
		logger:      logger,
		raftManager: fsm.NewRaftManager(manager, logger)}
}

func (c *ConfigServiceServer) GetConfig(ctx context.Context, config *configv1.GetConfigRequest) (*configv1.GetConfigResponse, error) {
	switch config.What {
	case configv1.GetConfigRequest_TYPE_ALL:
	case configv1.GetConfigRequest_TYPE_RAFT:
		switch config.Amount {
		case configv1.GetConfigRequest_SPECIFICITY_ONE:
			return c.getRaftConfig(config.Key)
		case configv1.GetConfigRequest_SPECIFICITY_EVERYTHING:
			return c.getAllRaftConfigs()
		}
	}
	return nil, errors.New("cannot determine which type of config to return")
}

func (c *ConfigServiceServer) getRaftConfig(name *string) (*configv1.GetConfigResponse, error) {
	if name == nil {
		return nil, errors.New("cannot request a named record without a key")
	}

	val, err := c.raftManager.Get(*name)
	if err != nil {
		return nil, err
	}

	t := &configv1.GetConfigResponse{
		Type: &configv1.GetConfigResponse_RaftConfig{
			RaftConfig: &configv1.GetRaftConfigurationResponse{
				Configuration: val,
			},
		},
	}

	return t, nil
}

func (c *ConfigServiceServer) getAllRaftConfigs() (*configv1.GetConfigResponse, error) {

	all, err := c.raftManager.GetAll()
	if err != nil {
		return nil, err
	}

	return &configv1.GetConfigResponse{
		Type: &configv1.GetConfigResponse_AllRaftConfigs{
			AllRaftConfigs: &configv1.ListRaftConfigurationResponse{
				AvailableConfigs: all,
			},
		},
	}, nil
}
