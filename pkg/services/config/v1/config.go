/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package v1

import (
	"context"

	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/fsm"
	v1 "r3t.io/pleiades/pkg/protocols/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
)

var (
	_ v1.ConfigService_Server = (*ConfigService)(nil)
)

type ConfigService struct {
	manager     *services.StoreManager
	logger      zerolog.Logger
	raftManager *fsm.OldRaftManager[v1.RaftConfiguration]
}

func NewConfigService(manager *services.StoreManager,
	logger zerolog.Logger,
	raftManager *fsm.OldRaftManager[v1.RaftConfiguration]) *ConfigService {
	return &ConfigService{
		manager:     manager,
		logger:      logger,
		raftManager: fsm.NewOldRaftManager(manager, logger),
	}
}

func (c *ConfigService) GetConfig(ctx context.Context, config v1.ConfigService_getConfig) error {
	//TODO implement me
	panic("implement me")
}
