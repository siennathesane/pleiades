/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package server

import (
	"context"
	"fmt"
	"net/http"

	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/lni/dragonboat/v3"
	dconfig "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/mxplusb/pleiades/pkg/api/raft/v1/raftv1connect"
	"github.com/mxplusb/pleiades/pkg/server/eventing"
	"github.com/mxplusb/pleiades/pkg/server/kvstore"
	"github.com/mxplusb/pleiades/pkg/server/raft"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/mxplusb/pleiades/pkg/server/serverutils"
	"github.com/mxplusb/pleiades/pkg/server/shard"
	"github.com/mxplusb/pleiades/pkg/server/transactions"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func init() {
	dlog.SetLoggerFactory(serverutils.DragonboatLoggerFactory)
}

// singletons
var (
	ServerModule = fx.Module("server",
		kvstore.KvStoreModule,
		raft.RaftModule,
		shard.ShardModule,
		transactions.TransactionsModule,
		eventing.EventingModule,
		fx.Provide(NewNodeHost),
		fx.Provide(NewHttpServeMux),
		fx.Invoke(NewHttpServer),
	)

	httpServer *http.Server
	nodeHost   *dragonboat.NodeHost
)

type HttpServeMuxBuilderParams struct {
	fx.In

	Logger   zerolog.Logger
	Handlers []runtime.ServiceHandler `group:"routes"`
}

type HttpServeMuxBuilderResults struct {
	fx.Out

	Mux *http.ServeMux
}

func NewHttpServeMux(params HttpServeMuxBuilderParams) HttpServeMuxBuilderResults {
	mux := http.NewServeMux()

	for _, route := range params.Handlers {
		params.Logger.Debug().Str("path", route.Path()).Msg("registering handler")
		mux.Handle(route.Path(), route)
	}

	// add grpc health checking
	checker := grpchealth.NewStaticChecker(
		kvstorev1connect.KvStoreServiceName,
		raftv1connect.HostServiceName)
	mux.Handle(grpchealth.NewHandler(checker))

	// add grpc reflection for grpcurl and other tools
	reflector := grpcreflect.NewStaticReflector(
		kvstorev1connect.KvStoreServiceName,
		raftv1connect.HostServiceName)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return HttpServeMuxBuilderResults{Mux: mux}
}

type HttpServerBuilderParams struct {
	fx.In

	Logger zerolog.Logger
	Config *viper.Viper
	Mux    *http.ServeMux
}

type HttpServerBuilderResults struct {
	fx.Out

	Server *http.Server
}

func NewHttpServer(lc fx.Lifecycle, params HttpServerBuilderParams) HttpServerBuilderResults {
	port := params.Config.GetUint("server.host.httpListenPort")
	if port == 0 {
		params.Logger.Fatal().Msg("http port cannot be 0!")
	}
	addr := fmt.Sprintf("%s:%d", params.Config.GetString("server.host.listenAddr"), params.Config.GetUint("server.host.httpListenPort"))
	httpServer = &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(params.Mux, &http2.Server{}),
	}

	params.Logger.Debug().Str("http-addr", addr).Msg("http listen address")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//goland:noinspection GoUnhandledErrorResult
			go httpServer.ListenAndServe()
			params.Logger.Info().Msg("started http server")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})
	return HttpServerBuilderResults{Server: httpServer}
}

type NodeHostBuilderParams struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Logger         zerolog.Logger
	NodeHostConfig dconfig.NodeHostConfig
	Server         *eventing.EventServer
}

func NewNodeHost(params NodeHostBuilderParams) (*dragonboat.NodeHost, error) {
	handler, err := params.Server.GetRaftSystemEventListener()
	if err != nil {
		params.Logger.Error().Err(err).Msg("can't build raft system listeners")
		return nil, err
	}

	params.NodeHostConfig.SystemEventListener = handler
	params.NodeHostConfig.RaftEventListener = handler

	nodeHost, err = dragonboat.NewNodeHost(params.NodeHostConfig)
	if err != nil {
		params.Logger.Error().Err(err).Msg("can't build node host")
	}

	// dragonboat starts itself when New() is created, this is purely for the startup sequence
	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			nodeHost.Stop()
			return nil
		},
	})

	return nodeHost, err
}
