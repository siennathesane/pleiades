/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package v1

import (
	"context"

	v1 "r3t.io/pleiades/pkg/protocols/config/v1"
)

var (
	_ v1.ConfigService_Server = (*ConfigService)(nil)
)

type ConfigService struct {

}

func (c *ConfigService) GetConfig(ctx context.Context, config v1.ConfigService_getConfig) error {
	//TODO implement me
	panic("implement me")
}

