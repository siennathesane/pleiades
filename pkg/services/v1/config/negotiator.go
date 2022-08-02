
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package config

import (
	"context"

	hostv1 "gitlab.com/anthropos-labs/pleiades/pkg/protocols/v1/host"
	"capnproto.org/go/capnp/v3/server"
	"github.com/rs/zerolog"
)

var (
	_ hostv1.Negotiator_Server = (*NegotiatorServer)(nil)
)

type NegotiatorServer struct {
	logger   zerolog.Logger
	registry *Registry
}

func NewNegotiator(logger zerolog.Logger, registry *Registry) *NegotiatorServer {
	return &NegotiatorServer{logger: logger, registry: registry}
}

func (n *NegotiatorServer) Register(t hostv1.ServiceType_Type, srv any) error {
	return n.registry.PutServer(t, srv)
}

func (n *NegotiatorServer) ConfigService(ctx context.Context, call hostv1.Negotiator_configService) error {
	results, err := call.AllocResults()
	if err != nil {
		n.logger.Error().Err(err).Msg("failed to allocate results")
		return err
	}

	val, err := n.registry.GetServer(hostv1.ServiceType_Type_configService)
	if err != nil {
		n.logger.Error().Err(err).Msg("failed to get host service")
		return err
	}

	target := val.(*ConfigServer)
	svc := hostv1.ConfigService_ServerToClient(target, &server.Policy{
		MaxConcurrentCalls: 250,
	})

	return results.SetSvc(svc)
}
