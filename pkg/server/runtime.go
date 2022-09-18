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
	"net/http"

	"github.com/mxplusb/api/kvstore/v1/kvstorev1connect"
	"github.com/mxplusb/api/raft/v1/raftv1connect"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/rs/zerolog"
)

func init() {
	dlog.SetLoggerFactory(configuration.DragonboatLoggerFactory)
}

type Options struct {
	GRPCPort int
	EmbeddedEventStreamPort int
	RaftPort int
}

type Server struct {
	logger                 zerolog.Logger
	nh                     *dragonboat.NodeHost
	raftHost               IHost
	raftShard              IShardManager
	raftTransactionManager ITransactionManager
	bboltStoreManager      IKVStore
}

func New(nhc dconfig.NodeHostConfig, mux *http.ServeMux, logger zerolog.Logger) (*Server, error) {
	srv := &Server{
		logger: logger.With().Str("component", "server").Logger(),
	}

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		return nil, errors.Wrap(err, "can't start node host")
	}

	rh := newRaftHost(nh, logger)
	rhAdapter := &raftHostConnectAdapter{
		logger: logger,
		host:   rh,
	}
	path, handler := raftv1connect.NewHostServiceHandler(rhAdapter)
	mux.Handle(path, handler)
	srv.raftHost = rh

	sm := newShardManager(nh, logger)
	smAdapter := &raftShardConnectAdapter{
		logger:       logger,
		shardManager: sm,
	}
	path, handler = raftv1connect.NewShardServiceHandler(smAdapter)
	mux.Handle(path, handler)
	srv.raftShard = sm

	tm := newTransactionManager(nh, logger)
	tmAdapter := &kvstoreTransactionConnectAdapter{
		logger:             logger,
		transactionManager: tm,
	}
	path, handler = kvstorev1connect.NewTransactionsServiceHandler(tmAdapter)
	mux.Handle(path, handler)
	srv.raftTransactionManager = tm

	store := newBboltStoreManager(tm, nh, logger)
	storeAdapter := &kvstoreBboltConnectAdapter{
		logger:       logger,
		storeManager: store,
	}
	path, handler = kvstorev1connect.NewKvStoreServiceHandler(storeAdapter)
	mux.Handle(path, handler)
	srv.bboltStoreManager = store

	srv.nh = nh

	return srv, nil
}

func (s *Server) GetRaftHost() IHost {
	return s.raftHost
}

func (s *Server) GetRaftTransactionManager() ITransactionManager {
return s.raftTransactionManager
}

func (s *Server) GetRaftKVStore() IKVStore{
return s.bboltStoreManager
}

func (s *Server) GetRaftShardManager() IShardManager {
	return s.raftShard
}

func (s *Server) Stop() {
	s.nh.Stop()
}