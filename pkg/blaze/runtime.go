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
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/aperturerobotics/starpc/srpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	multiplex "github.com/libp2p/go-mplex"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

func NewRuntime(ctx context.Context, conf *conf.NodeHostConfig, clogger conf.Logger) (*Runtime, error) {
	logger := clogger.GetLogger()
	l := logger.With().Str("component", "runtime").Logger()
	mux := srpc.NewMux()

	run := &Runtime{
		mux:    mux,
		logger: l,
		ctx:    ctx,
		addrs:  make([]multiaddr.Multiaddr, 0),
	}

	err := RegisterNodeHostRpcServer(mux, conf, clogger)
	if err != nil {
		l.Error().Err(err).Msg("failed to register node host rpc server")
		return nil, err
	}

	listenSplit := strings.Split(conf.ListenAddress, ":")
	if len(listenSplit) != 2 {
		l.Error().Msgf("invalid listen address: %s", conf.ListenAddress)
		return run, err
	}

	quicAddr := fmt.Sprintf("/ip4/%s/udp/%s/quic", listenSplit[0], listenSplit[1])
	websocketAddr := fmt.Sprintf("/ip4/%s/tcp/%s/ws", listenSplit[0], listenSplit[1])
	multiAddrs := []string{quicAddr, websocketAddr}
	l.Info().Msgf("will listen on %v", multiAddrs)

	for idx := range multiAddrs {
		ma, err := multiaddr.NewMultiaddr(multiAddrs[idx])
		if err != nil {
			l.Error().Err(err).Msgf("failed to create multiaddr: %s", multiAddrs[idx])
			return run, err
		}
		run.addrs = append(run.addrs, ma)
	}

	run.privKey, _, err = crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		l.Error().Err(err).Msg("failed to generate key pair")
		return run, err
	}

	run.peerId, err = peer.IDFromPrivateKey(run.privKey)
	if err != nil {
		l.Error().Err(err).Msg("failed to generate peer id")
		return run, err
	}

	run.quicTr, err = libp2pquic.NewTransport(run.privKey, nil, nil, nil)
	if err != nil {
		l.Error().Err(err).Msg("failed to create transport")
		return run, err
	}

	return run, nil
}

type Runtime struct {
	addrs    []multiaddr.Multiaddr
	ctx      context.Context
	host     host.Host
	listener transport.Listener
	logger   zerolog.Logger
	mp       multiplex.Multiplex
	mux      srpc.Mux
	peerId   peer.ID
	privKey  crypto.PrivKey
	srv      *srpc.Server
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

	r.host.SetStreamHandler(NodeHostProtocolVersion, r.handleStream)

	return nil
}

func (r *Runtime) Stop() {
	if err := r.host.Close(); err != nil {
		r.logger.Fatal().Err(err).Msg("cannot cleanly shut down networking")
	}
}

func (r *Runtime) handleStream(stream network.Stream) {
	if err := r.srv.HandleStream(r.ctx, stream); err != nil {
		r.logger.Error().Err(err).Str("stream-id", stream.ID()).Msg("cannot handle stream")
		if err := stream.Reset(); err != nil {
			r.logger.Error().Err(err).Str("stream-id", stream.ID()).Msg("cannot reset stream")
			return
		}
	}
}
