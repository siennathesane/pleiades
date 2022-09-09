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

	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/transport"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	multiplex "github.com/libp2p/go-mplex"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

func init() {
	dlog.SetLoggerFactory(conf.DragonboatLoggerFactory)
}

type Runtime struct {
	addrs    []multiaddr.Multiaddr
	ctx      context.Context
	host     host.Host
	listener transport.Listener
	logger   zerolog.Logger
	mp       multiplex.Multiplex
	peerId   peer.ID
	privKey  crypto.PrivKey
	quicTr   transport.Transport
}

func (r *Runtime) Run() error {

	transports := libp2p.ChainOptions(libp2p.Transport(r.quicTr),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(websocket.New))

	muxer := libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport)

	var err error
	r.host, err = libp2p.New(transports,
		muxer,
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
		libp2p.ListenAddrs(r.addrs...),
		libp2p.Identity(r.privKey))
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to create host")
		return err
	}

	//r.host.SetStreamHandler(RaftControlProtocolVersion, r.handleStream)

	return nil
}

func (r *Runtime) Stop() {
	if err := r.host.Close(); err != nil {
		r.logger.Fatal().Err(err).Msg("cannot cleanly shut down networking")
	}
}

