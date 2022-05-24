package servers

import (
	"context"
	"errors"

	dlog "github.com/lni/dragonboat/v3/logger"
	"r3t.io/pleiades/pkg/fsm"
	"r3t.io/pleiades/pkg/pb"
	"r3t.io/pleiades/pkg/servers/services"
)

type ConfigServer struct {
	pb.UnimplementedConfigServiceServer
	manager     *services.StoreManager
	logger      dlog.ILogger
	raftManager *fsm.RaftManager[pb.RaftConfig]
}

func NewConfigServer(manager *services.StoreManager, logger dlog.ILogger) ConfigServer {
	return ConfigServer{
		manager:     manager,
		logger:      logger,
		raftManager: fsm.NewRaftManager(manager, logger)}
}

func (c ConfigServer) GetConfig(ctx context.Context, config *pb.ConfigRequest) (*pb.ConfigResponse, error) {
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

func (c ConfigServer) getRaftConfig(name *string) (*pb.ConfigResponse, error) {
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

func (c ConfigServer) getAllRaftConfigs() (*pb.ConfigResponse, error) {

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

func (c ConfigServer) mustEmbedUnimplementedConfigServiceServer() {}
