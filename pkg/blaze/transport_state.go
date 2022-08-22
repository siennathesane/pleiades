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
	"io"
	"time"

	v1 "github.com/mxplusb/pleiades/pkg/api/v1"
	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p-core/network"
	"google.golang.org/protobuf/proto"
)

var (
	StateDuration      = 1 * time.Second
	HeaderDuration     = 2 * time.Second
	ReadDuration       = 5 * time.Second
	WriteDuration      = 5 * time.Second
	stateSize      int = 0
)

func VerifyStreamState(stream network.Stream) error {
	readDeadline := time.Now().Add(StateDuration)
	if err := stream.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	stateBuf := make([]byte, stateSize)
	if _, err := io.ReadFull(stream, stateBuf); err != nil {
		return err
	}

	state := &v1.State{}
	if err := proto.Unmarshal(stateBuf, state); err != nil {
		return err
	}

	if state.State > 0 {
		return errors.New("invalid state")
	}

	return nil
}

func SendStreamState(stream network.Stream, streamState StreamState, followthrough bool) error {
	stateDuration := time.Now().Add(StateDuration)
	if err := stream.SetWriteDeadline(stateDuration); err != nil {
		return err
	}

	state := &v1.State{State: uint32(streamState)}

	// make sure we set the bits
	if followthrough {
		state.HeaderToFollow = 1
	} else {
		state.HeaderToFollow = 0
	}

	buf := make([]byte, state.SizeVT())
	if err := proto.Unmarshal(buf, state); err != nil {
		return err
	}

	if _, err := stream.Write(buf); err != nil {
		return err
	}

	return nil
}

