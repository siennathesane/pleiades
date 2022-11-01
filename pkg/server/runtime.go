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

	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/mxplusb/pleiades/pkg/server/eventing"
	"github.com/mxplusb/pleiades/pkg/server/kvstore"
	"github.com/mxplusb/pleiades/pkg/server/raft"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/mxplusb/pleiades/pkg/server/shard"
	"github.com/mxplusb/pleiades/pkg/server/transactions"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/rs/zerolog"
)

func init() {
	dlog.SetLoggerFactory(serverutils.DragonboatLoggerFactory)
}

type Options struct {
	GRPCPort int
	RaftPort int
}

type Server struct {
	logger                 zerolog.Logger
	nh                     *dragonboat.NodeHost
	raftHost               runtime.IHost
	raftShard              runtime.IShardManager
	raftTransactionManager runtime.ITransactionManager
	bboltStoreManager      runtime.IKVStore
	eventServer            *eventing.Server
}

func New(nhc dconfig.NodeHostConfig, mux *http.ServeMux, logger zerolog.Logger) (*Server, error) {
	srv := &Server{
		logger: logger.With().Str("component", "server").Logger(),
	}

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		return nil, errors.Wrap(err, "can't start node host")
	}

	srv.eventServer, err = eventing.NewServer(logger)
	if err != nil {
		return nil, errors.Wrap(err, "can't instantiate nats")
	}

	rh := raft.NewRaftHost(nh, logger)
	rhAdapter := raft.NewRaftHostConnectAdapter(rh, logger)
	path, handler := raftv1connect.NewHostServiceHandler(rhAdapter)
	mux.Handle(path, handler)
	srv.raftHost = rh

	shardManagerClient, err := srv.eventServer.GetStreamClient()
	if err != nil {
		return nil, errors.Wrap(err, "can't create stream client")
	}
	sm := shard.NewShardManager(nh, shardManagerClient, nil, logger)
	smAdapter := shard.NewRaftShardConnectAdapter(sm, logger)
	path, handler = raftv1connect.NewShardServiceHandler(smAdapter)
	mux.Handle(path, handler)
	srv.raftShard = sm

	tm := transactions.NewTransactionManager(nh, logger)
	tmAdapter := kvstore.NewKvstoreTransactionConnectAdapter(tm, logger)
	path, handler = kvstorev1connect.NewTransactionsServiceHandler(tmAdapter)
	mux.Handle(path, handler)
	srv.raftTransactionManager = tm

	store := kvstore.NewBboltStoreManager(tm, nh, logger)
	storeAdapter := kvstore.NewKvstoreBboltConnectAdapter(store, logger)
	path, handler = kvstorev1connect.NewKvStoreServiceHandler(storeAdapter)
	mux.Handle(path, handler)
	srv.bboltStoreManager = store

	srv.nh = nh

	return srv, nil
}

func (s *Server) GetRaftHost() runtime.IHost {
	return s.raftHost
}

func (s *Server) GetRaftTransactionManager() runtime.ITransactionManager {
	return s.raftTransactionManager
}

func (s *Server) GetRaftKVStore() runtime.IKVStore {
	return s.bboltStoreManager
}

func (s *Server) GetRaftShardManager() runtime.IShardManager {
	return s.raftShard
}

func (s *Server) Stop() {
	s.nh.Stop()
}
