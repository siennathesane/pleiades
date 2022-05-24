package managers

import (
	"encoding/json"
	"fmt"
	"reflect"

	dlog "github.com/lni/dragonboat/v3/logger"
	"r3t.io/pleiades/pkg/services"
	"r3t.io/pleiades/pkg/types"
)

var _ services.IStore[types.RaftConfig] = &RaftManager[types.RaftConfig]{}

type RaftManager struct {
	logger dlog.ILogger
	store  *services.StoreManager
}

func NewRaftManager(store *services.StoreManager, logger dlog.ILogger) *RaftManager {
	return &RaftManager{logger: logger, store: store}
}

func (rm *RaftManager) Get(key string) (*types.RaftConfig, error) {
	payload, err := rm.store.Get(key, reflect.TypeOf(&types.RaftConfig{}))

	if err != nil {
		rm.logger.Errorf(fmt.Errorf("error fetching %s from raft store: %w", key, err).Error())
		return &types.RaftConfig{}, err
	}

	config := &types.RaftConfig{}
	if err := config.Unmarshal(payload); err != nil {
		return &types.RaftConfig{}, err
	}
	return config, nil
}

func (rm *RaftManager) GetAll() (map[string]*types.RaftConfig, error) {
	respMap, err := rm.store.GetAll(reflect.TypeOf(&types.RaftConfig{}))
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*types.RaftConfig)
	for k, _ := range respMap {
		var c *types.RaftConfig
		if err := json.Unmarshal(respMap[k], &c); err != nil {
			return nil, err
		}
		configs[k] = c
	}

	return configs, nil
}

func (rm *RaftManager) Put(key string, payload *types.RaftConfig) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return rm.store.Put(key, encoded, reflect.TypeOf(&types.RaftConfig{}))
}
