package fsm

import (
	"reflect"

	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/pb"
	"r3t.io/pleiades/pkg/servers/services"
)

var _ services.IStore[pb.RaftConfig] = &RaftManager[pb.RaftConfig]{}

type RaftManager[T pb.RaftConfig] struct {
	logger zerolog.Logger
	store  *services.StoreManager
}

func NewRaftManager(store *services.StoreManager, logger zerolog.Logger) *RaftManager[pb.RaftConfig] {
	return &RaftManager[pb.RaftConfig]{logger: logger, store: store}
}

func (rm *RaftManager[T]) Get(key string) (*pb.RaftConfig, error) {
	payload, err := rm.store.Get(key, reflect.TypeOf(&pb.RaftConfig{}))

	if err != nil {
		rm.logger.Err(err).Str("key", key).Msg("error fetching key from raft store")
		return &pb.RaftConfig{}, err
	}

	config := &pb.RaftConfig{}
	if err := config.UnmarshalVT(payload); err != nil {
		return &pb.RaftConfig{}, err
	}
	return config, nil
}

func (rm *RaftManager[T]) GetAll() (map[string]*pb.RaftConfig, error) {
	respMap, err := rm.store.GetAll(reflect.TypeOf(&pb.RaftConfig{}))
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*pb.RaftConfig)
	for k, _ := range respMap {
		c := &pb.RaftConfig{}
		if err := c.UnmarshalVT(respMap[k]); err != nil {
			return nil, err
		}
		configs[k] = c
	}

	return configs, nil
}

func (rm *RaftManager[T]) Put(key string, payload *pb.RaftConfig) error {
	encoded, err := payload.MarshalVT()
	if err != nil {
		return err
	}
	return rm.store.Put(key, encoded, reflect.TypeOf(&pb.RaftConfig{}))
}
