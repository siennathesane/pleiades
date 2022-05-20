package servers

import (
	"context"

	"r3t.io/pleiades/pkg/protos"
)

var _ RaftConfigServiceServer = RaftConfigServer{}

type RaftConfigServer struct {
	UnimplementedRaftConfigServiceServer
}

func (r RaftConfigServer) AddConfiguration(ctx context.Context, request *protos.NewRaftConfigRequest) (*protos.RaftConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) GetConfiguration(ctx context.Context, request *protos.GetRaftConfigRequest) (*protos.GetRaftConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) ListConfigurations(ctx context.Context, configs *protos.ListRaftConfigs) (*protos.ListRaftConfigsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r RaftConfigServer) mustEmbedUnimplementedRaftConfigServiceServer() {
	//TODO implement me
	panic("implement me")
}
