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
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestFrame(t *testing.T) {
	suite.Run(t, new(FrameTests))
}

type FrameTests struct {
	suite.Suite
}

func (ft *FrameTests) TestNewFrame() {
	known := &Frame{
		state:                 byte(InvalidByte),
		version:               byte(Version1),
		reserved0:             ReservedByte,
		reserved1:             ReservedByte,
		service:               byte(UnspecifiedService),
		method:                0xff,
		authorizationService:  byte(NoAuthorizationService),
		authorizationMethod:   byte(NoAuthorizationMethod),
		authorizationSize:     [4]byte{0x00, 0x00, 0x00, 0x00},
		payloadSize:           [4]byte{0x00, 0x00, 0x00, 0x00},
		authorization:         nil,
		payload:               nil,
		authorizationChecksum: [4]byte{0x00, 0x00, 0x00, 0x00},
		payloadChecksum:       [4]byte{0x00, 0x00, 0x00, 0x00},
	}

	target := NewFrame()
	ft.Require().Equal(known, target, "the created frame must equal the known frame")
}

func (ft *FrameTests) TestWithState() {
	target := NewFrame().WithState(ValidByte)
	ft.Require().Equal(ValidByte, target.GetState(), "the state must be considered valid")
}

func (ft *FrameTests) TestWithService() {
	target := NewFrame().WithService(RaftControlServiceByte)
	ft.Require().Equal(RaftControlServiceByte, target.GetService(), "the service must be the raft control service")
}

func (ft *FrameTests) TestWithMethod() {
	testMethod := MethodByte(0x54)
	target := NewFrame().WithMethod(testMethod)
	ft.Require().Equal(testMethod, target.GetMethod(), "the method value must be set")
}

func (ft *FrameTests) TestWithPayload() {
	testPayload := []byte("12345")
	targetChecksum := crc32.ChecksumIEEE(testPayload)

	var target *Frame
	ft.Require().NotPanics(func() {
		target = NewFrame().WithPayload(testPayload)
	}, "assigning the test payload must not panic")
	payloadSize := binary.LittleEndian.Uint32(target.payloadSize[:])
	ft.Require().Equal(len(testPayload), int(payloadSize), "the payload size must be equal to the known payload size")
	ft.Require().Equal(testPayload, target.GetPayload(), "the target payload must equal the test payload")

	payloadChecksum := binary.LittleEndian.Uint32(target.payloadChecksum[:])
	ft.Require().Equal(targetChecksum, payloadChecksum, "the payload checksums must be equal")
}

func (ft *FrameTests) TestSerialization() {
	testMethod := MethodByte(0x54)
	payload := []byte("12345")
	payloadChecksum := crc32.ChecksumIEEE(payload)

	frameHeaders := 24 // 8 header bytes + 4 uint32s = 24 bytes
	totalFrameLength := len(payload) + frameHeaders

	outgoing := NewFrame().WithState(ValidByte).WithService(RaftControlServiceByte).WithMethod(testMethod).WithPayload(payload)
	ft.Require().Equal(ValidByte, outgoing.GetState())
	ft.Require().Equal(RaftControlServiceByte, outgoing.GetService())
	ft.Require().Equal(testMethod, outgoing.GetMethod())
	ft.Require().Equal(payload, outgoing.GetPayload())
	ft.Require().Equal(payloadChecksum, binary.LittleEndian.Uint32(outgoing.payloadChecksum[:]))

	marshaled, err := outgoing.Marshal()
	ft.Require().NoError(err)

	ft.Require().Equal(totalFrameLength, len(marshaled), "the marshalled frame length must equal the total frame length")

	buf := bytes.NewBuffer(marshaled)

	incoming := NewFrame()
	read, err := incoming.ReadFrom(buf)
	ft.Require().NoError(err, "there must not be an error when trying to serialize the frame")
	ft.Require().Equal(int64(totalFrameLength), read, "the length of bytes read must equal the total frame length")
}
