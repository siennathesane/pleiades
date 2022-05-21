package servers

import (
	"context"

	"r3t.io/pleiades/pkg/types"
)

var _ RaftConfigServiceServer = RaftConfigServer{}

type RaftConfigServer struct {
	UnimplementedRaftConfigServiceServer
}

func (r RaftConfigServer) AddConfiguration(ctx context.Context, request *types.NewRaftConfigRequest) (*types.RaftConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) GetConfiguration(ctx context.Context, request *types.GetRaftConfigRequest) (*types.GetRaftConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) ListConfigurations(ctx context.Context, configs *types.ListRaftConfigs) (*types.ListRaftConfigsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) mustEmbedUnimplementedRaftConfigServiceServer() {
	//TODO implement me
	panic("implement me")
}
