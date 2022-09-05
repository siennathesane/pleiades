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
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

func randomLibp2pTestHost() host.Host {
	rand.Seed(time.Now().UTC().UnixNano())
	port := 1024 + rand.Intn(65535-1024)
	hostAddr := fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic", port)

	ma, _ := multiaddr.NewMultiaddr(hostAddr)

	lhost, _ := libp2p.New(libp2p.ListenAddrs(ma))
	return lhost
}
