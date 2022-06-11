package servers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/fsm"
	"r3t.io/pleiades/pkg/pb"
	"r3t.io/pleiades/pkg/servers/services"
)

var _ pb.DRPCConfigServiceServer = (*ConfigServiceServer)(nil)

type ConfigServiceServer struct {
	pb.DRPCConfigServiceUnimplementedServer
	manager     *services.StoreManager
	logger      zerolog.Logger
	raftManager *fsm.RaftManager[pb.RaftConfig]
}

func NewConfigServiceServer(manager *services.StoreManager, logger zerolog.Logger) *ConfigServiceServer {
	return &ConfigServiceServer{
		manager:     manager,
		logger:      logger,
		raftManager: fsm.NewRaftManager(manager, logger)}
}

func (c *ConfigServiceServer) GetConfig(ctx context.Context, config *pb.ConfigRequest) (*pb.ConfigResponse, error) {
	switch config.What {
	case pb.ConfigRequest_ALL:
	case pb.ConfigRequest_RAFT:
		switch config.Amount {
		case pb.ConfigRequest_ONE:
			return c.getRaftConfig(config.Key)
		case pb.ConfigRequest_EVERYTHING:
			return c.getAllRaftConfigs()
		}
	}
	return nil, errors.New("cannot determine which type of config to return")
}

func (c *ConfigServiceServer) getRaftConfig(name *string) (*pb.ConfigResponse, error) {
	if name == nil {
		return nil, errors.New("cannot request a named record without a key")
	}

	val, err := c.raftManager.Get(*name)
	if err != nil {
		return nil, err
	}

	t := &pb.ConfigResponse{
		Type: &pb.ConfigResponse_RaftConfig{
			RaftConfig: &pb.GetRaftConfigResponse{
				Configuration: val,
			},
		},
	}

	return t, nil
}

func (c *ConfigServiceServer) getAllRaftConfigs() (*pb.ConfigResponse, error) {

	all, err := c.raftManager.GetAll()
	if err != nil {
		return nil, err
	}

	return &pb.ConfigResponse{
		Type: &pb.ConfigResponse_AllRaftConfigs{
			AllRaftConfigs: &pb.ListRaftConfigsResponse{
				AvailableConfigs: all,
			},
		},
	}, nil
}
