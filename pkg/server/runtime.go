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
	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
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

func New(nhc dconfig.NodeHostConfig, gServer *grpc.Server, logger zerolog.Logger) (*Server, error) {
	srv := &Server{
		logger: logger.With().Str("component", "server").Logger(),
	}

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		return nil, errors.Wrap(err, "can't start node host")
	}

	rh := newRaftHost(nh, logger)
	rhAdapter := &raftHostGrpcAdapter{
		logger: logger,
		host:   rh,
	}
	raftv1.RegisterHostServiceServer(gServer, rhAdapter)
	srv.raftHost = rh

	sm := newShardManager(nh, logger)
	smAdapter := &raftShardGrpcAdapter{
		logger:       logger,
		shardManager: sm,
	}
	raftv1.RegisterShardServiceServer(gServer, smAdapter)
	srv.raftShard = sm

	tm := newTransactionManager(nh, logger)
	tmAdapter := &raftTransactionGrpcAdapter{
		logger:             logger,
		transactionManager: tm,
	}
	RegisterTransactionsServer(gServer, tmAdapter)
	srv.raftTransactionManager = tm

	store := newBboltStoreManager(tm, nh, logger)
	storeAdapter := &raftBBoltStoreManagerGrpcAdapter{
		logger:       logger,
		storeManager: store,
	}
	RegisterKVStoreServiceServer(gServer, storeAdapter)
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