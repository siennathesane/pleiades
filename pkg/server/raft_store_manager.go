/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	"context"
	"time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/rs/zerolog"
)

var (
	_ IStore = (*StoreManager)(nil)
)

func newStoreManager(logger zerolog.Logger, nh *dragonboat.NodeHost) *StoreManager {
	l := logger.With().Str("component", "store-manager").Logger()
	return &StoreManager{l, nh}
}

type StoreManager struct {
	logger zerolog.Logger
	nodeHost *dragonboat.NodeHost
}

func (s *StoreManager) SyncPropose(ctx context.Context, session *client.Session, cmd []byte) (statemachine.Result, error) {
	return s.nodeHost.SyncPropose(ctx, session, cmd)
}

func (s *StoreManager) SyncRead(ctx context.Context, clusterID uint64, query interface{}) (interface{}, error) {
	return s.nodeHost.SyncRead(ctx, clusterID, query)
}

func (s *StoreManager) Propose(session *client.Session, cmd []byte, timeout time.Duration) (*dragonboat.RequestState, error) {
	return s.nodeHost.Propose(session, cmd, timeout)
}


