/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"context"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/rs/zerolog"
)

var (
	_ ITransactionManager = (*SessionManager)(nil)
)

func newSessionManager(nh *dragonboat.NodeHost, logger zerolog.Logger) *SessionManager {
	l := logger.With().Str("component", "session-manager").Logger()
	return &SessionManager{l, nh}
}

type SessionManager struct {
	logger zerolog.Logger
	nh     *dragonboat.NodeHost
}

func (s *SessionManager) GetNoOpSession(clusterID uint64) *client.Session {
	return s.nh.GetNoOPSession(clusterID)
}

func (s *SessionManager) GetSession(ctx context.Context, clusterID uint64) (*client.Session, error) {
	return s.nh.SyncGetSession(ctx, clusterID)
}

func (s *SessionManager) CloseSession(ctx context.Context, cs *client.Session) error {
	return s.nh.SyncCloseSession(ctx, cs)
}
