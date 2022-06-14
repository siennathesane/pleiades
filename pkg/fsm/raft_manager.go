package fsm

import (
	"reflect"

	"github.com/rs/zerolog"
	proto2 "google.golang.org/protobuf/proto"
	configv1 "r3t.io/pleiades/pkg/pb/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
)

var _ services.IStore[configv1.RaftConfig] = &RaftManager[configv1.RaftConfig]{}

type RaftManager[T configv1.RaftConfig] struct {
	logger zerolog.Logger
	store  *services.StoreManager
}

func NewRaftManager(store *services.StoreManager, logger zerolog.Logger) *RaftManager[configv1.RaftConfig] {
	return &RaftManager[configv1.RaftConfig]{logger: logger, store: store}
}

func (rm *RaftManager[T]) Get(key string) (*configv1.RaftConfig, error) {
	payload, err := rm.store.Get(key, reflect.TypeOf(&configv1.RaftConfig{}))

	if err != nil {
		rm.logger.Err(err).Str("key", key).Msg("error fetching key from raft store")
		return &configv1.RaftConfig{}, err
	}

	var config configv1.RaftConfig
	if err := proto2.Unmarshal(payload, &config); err != nil {
		return &configv1.RaftConfig{}, err
	}
	return &config, nil
}

func (rm *RaftManager[T]) GetAll() (map[string]*configv1.RaftConfig, error) {
	respMap, err := rm.store.GetAll(reflect.TypeOf(&configv1.RaftConfig{}))
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*configv1.RaftConfig)
	for k, _ := range respMap {
		var c configv1.RaftConfig
		if err := proto2.Unmarshal(respMap[k], &c); err != nil {
			return nil, err
		}
		configs[k] = &c
	}

	return configs, nil
}

func (rm *RaftManager[T]) Put(key string, payload *configv1.RaftConfig) error {
	encoded, err := proto2.Marshal(payload)
	if err != nil {
		return err
	}
	return rm.store.Put(key, encoded, reflect.TypeOf(&configv1.RaftConfig{}))
}
