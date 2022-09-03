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
	"encoding/binary"
	"hash/crc32"
	"io"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const (
	bytesInInt32 = 4
	ReservedByte = 0xbc
)

var (
	ErrInvalidHeaders  = errors.New("invalid header length")
	ErrInvalidChecksum = errors.New("invalid checksum")
)

type VersionByte byte

const (
	CurrentVersion VersionByte = 0x01
	Version1       VersionByte = 0x01
)

var (
	ErrUnsupportedVersion = errors.New("unsupported version")
)

type StateByte byte

const (
	ValidByte   StateByte = 0x00
	InvalidByte StateByte = 0x01

	// 10-49 Reserved - Future Use
	// 50-100 Stream Markers
	StreamStartByte    StateByte = 0x50
	StreamContinueByte StateByte = 0x51
	StreamErrorByte    StateByte = 0x52
	StreamEndByte      StateByte = 0x53

	// 250-254 System Usage
	SystemErrorByte StateByte = 0xfa
)

var (
	ErrInvalidState = errors.New("unsupported state")
)

type ServiceByte byte

const (
	// 0-9: Reserved
	// 10-29: Raft Services
	RaftControlServiceByte ServiceByte = 0x10
	RaftClusterServiceByte ServiceByte = 0x11
	RaftTransportService   ServiceByte = 0x12

	// 30-49: Session Services
	SessionServiceByte ServiceByte = 0x30

	// 50-69: Authentication Services
	BasicAuthServiceByte ServiceByte = 0x50
	UnspecifiedService   ServiceByte = 0xff
)

var (
	ErrUnsupportedService = errors.New("unsupported service")
)

// MethodByte represents a type to be implemented by each service
type MethodByte byte

type AuthorizationServiceByte byte

const (
	// 0-253: Reserved for future use
	// 254: No auth
	NoAuthorizationService AuthorizationServiceByte = 0xff
)

type AuthorizationMethodByte byte

const (
	// 0-253: Reserved for future use
	// 254: No auth
	NoAuthorizationMethod AuthorizationMethodByte = 0xff
)

// NewFrame generates a new *Frame with zeroed out fields.
func NewFrame() *Frame {
	return &Frame{
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
}

// Frame represents a single packet in a stream.
type Frame struct {
	// The overall state of the stream
	state byte

	// reserved0 is a reserved byte
	version byte
	// reserved0 is a reserved byte
	reserved0 byte
	// reserved1 is a reserved byte
	reserved1 byte

	// The service byte represents which service to route the payload to
	service byte

	// The method byte represents which method the service will call
	method byte

	// The authorizationService byte represents which authorization service to route authorization to
	authorizationService byte
	// The authorizationMethod byte represents a specific authorization method to call
	authorizationMethod byte

	// The authorizationSize field is an int32 representing the size of the authorization payload
	authorizationSize [4]byte

	// The payloadSize field is an int32 representing the size of the payload
	payloadSize [4]byte

	// The authorization field contains authorization information for the request
	authorization []byte

	// The payload field contains the application-level information
	payload []byte

	// the authorizationChecksum is a 32-bit CRC for the authorization field
	authorizationChecksum [4]byte

	// The payloadChecksum field is a 32-bit CRC for the payload
	payloadChecksum [4]byte
}

// ReadFrom reads from io.Reader r to fill out a frame. This method overwrites the frame
func (f *Frame) ReadFrom(r io.Reader) (n int64, err error) {
	readThusFar := int64(0)

	// header is 16 bytes, so read the first 16
	headerBuf := make([]byte, 16)
	read, err := io.ReadFull(r, headerBuf)
	if err != nil {
		return 0, ErrInvalidHeaders
	}
	if read < 16 {
		return int64(read), ErrInvalidHeaders
	}

	readThusFar += 16

	if err := validateHeader(headerBuf[0]); err != nil {
		return int64(read), err
	}
	f.state = headerBuf[0]

	if err := validateVersion(headerBuf[1]); err != nil {
		return int64(read), err
	}
	f.version = headerBuf[1]

	if err := validateService(headerBuf[4]); err != nil {
		return int64(read), err
	}
	f.service = headerBuf[4]

	// the method bytes don't really matter since they're all service specific anyways
	f.method = headerBuf[5]
	f.authorizationService = headerBuf[6]
	f.authorizationMethod = headerBuf[7]

	f.authorizationSize = [4]byte{headerBuf[8], headerBuf[9], headerBuf[10], headerBuf[11]}
	f.payloadSize = [4]byte{headerBuf[12], headerBuf[13], headerBuf[14], headerBuf[15]}

	// to handle the unimplemented auth case
	authSize := binary.LittleEndian.Uint32(f.authorizationSize[:])
	if authSize > 0 {
		authBuf := make([]byte, authSize)
		read, err = io.ReadFull(r, authBuf)
		if err != nil {
			return int64(read), err
		}
		f.authorization = authBuf
		readThusFar += int64(read)
	}

	// prepare and fetch the payload
	payloadSize := binary.LittleEndian.Uint32(f.payloadSize[:])
	if payloadSize == 0 {
		return int64(read), errors.New("payload length is zero")
	}

	f.payload = make([]byte, payloadSize)
	read, err = io.ReadFull(r, f.payload)
	if err != nil {
		return int64(read), err
	}
	readThusFar += int64(payloadSize)

	// grab the auth checksum
	read, err = io.ReadFull(r, f.authorizationChecksum[:])
	if err != nil {
		return int64(read), err
	}
	authChecksum := binary.LittleEndian.Uint32(f.authorizationChecksum[:])
	if err := validateAuthChecksum(int32(authChecksum)); err != nil {
		return int64(read), err
	}
	readThusFar += 4

	// grab the payload checksum
	read, err = io.ReadFull(r, f.payloadChecksum[:])
	if err != nil {
		return int64(read), err
	}
	payloadChecksum := binary.LittleEndian.Uint32(f.payloadChecksum[:])
	if err := validatePayloadChecksum(payloadChecksum, f.payload); err != nil {
		return int64(read), err
	}
	readThusFar += 4

	return readThusFar, nil
}

func validatePayloadChecksum(checksum uint32, payload []byte) error {
	incomingSum := crc32.ChecksumIEEE(payload)
	if incomingSum != checksum {
		return ErrInvalidChecksum
	}
	return nil
}

// auth is not currently implemented
func validateAuthChecksum(checksum int32) error {
	return nil
}

// auth is currently not implemented
func validateAuthService(svc byte) error {
	return nil
}

func validateService(svc byte) error {
	switch ServiceByte(svc) {
	case RaftControlServiceByte:
		return nil
	case RaftClusterServiceByte:
		return nil
	case SessionServiceByte:
		return nil
	case BasicAuthServiceByte:
		return nil
	case UnspecifiedService:
		return nil
	default:
		return ErrUnsupportedService
	}
}

func validateVersion(version byte) error {
	switch VersionByte(version) {
	case Version1:
		return nil
	default:
		return ErrUnsupportedVersion
	}
}

func validateHeader(state byte) error {
	switch StateByte(state) {
	case ValidByte:
		return nil
	case InvalidByte:
		return nil
	case StreamStartByte:
		return nil
	case StreamContinueByte:
		return nil
	case StreamErrorByte:
		return nil
	case StreamEndByte:
		return nil
	case SystemErrorByte:
		return nil
	default:
		return ErrInvalidState
	}
}

func (f *Frame) Marshal() ([]byte, error) {
	target := make([]byte, 16)

	target[0] = f.state
	target[1] = f.version
	target[2] = f.reserved0
	target[3] = f.reserved1
	target[4] = f.service
	target[5] = f.method
	target[6] = f.authorizationService
	target[7] = f.authorizationMethod

	// while authorization is unsupported
	f.authorizationSize[0] = 0x00
	f.authorizationSize[1] = 0x00
	f.authorizationSize[2] = 0x00
	f.authorizationSize[3] = 0x00

	target[8] = f.authorizationSize[0]
	target[9] = f.authorizationSize[1]
	target[10] = f.authorizationSize[2]
	target[11] = f.authorizationSize[3]

	target[12] = f.payloadSize[0]
	target[13] = f.payloadSize[1]
	target[14] = f.payloadSize[2]
	target[15] = f.payloadSize[3]

	target = append(target, f.authorization...)
	target = append(target, f.payload...)

	// this is while auth is still unsupported
	f.authorizationChecksum[0] = 0x00
	f.authorizationChecksum[1] = 0x00
	f.authorizationChecksum[2] = 0x00
	f.authorizationChecksum[3] = 0x00

	target = append(target,
		f.authorizationChecksum[0],
		f.authorizationChecksum[1],
		f.authorizationChecksum[2],
		f.authorizationChecksum[3])

	// just as an in case type thing
	f.setPayloadChecksum()

	target = append(target,
		f.payloadChecksum[0],
		f.payloadChecksum[1],
		f.payloadChecksum[2],
		f.payloadChecksum[3])

	return target, nil
}

func (f *Frame) WithState(state StateByte) *Frame {
	f.state = byte(state)
	return f
}

func (f *Frame) GetState() StateByte {
	return StateByte(f.state)
}

func (f *Frame) WithService(service ServiceByte) *Frame {
	f.service = byte(service)
	return f
}

func (f *Frame) GetService() ServiceByte {
	return ServiceByte(f.service)
}

func (f *Frame) WithMethod(method MethodByte) *Frame {
	f.method = byte(method)
	return f
}

func (f *Frame) GetMethod() MethodByte {
	return MethodByte(f.method)
}

func (f *Frame) WithPayload(payload []byte) *Frame {
	f.payload = payload
	size := uint32ToBytes(uint32(len(payload)))
	f.payloadSize[0] = size[0]
	f.payloadSize[1] = size[1]
	f.payloadSize[2] = size[2]
	f.payloadSize[3] = size[3]

	f.setPayloadChecksum()

	return f
}

func (f *Frame) setPayloadChecksum() {
	checksumBytes := uint32ToBytes(crc32.ChecksumIEEE(f.payload))
	f.payloadChecksum[0] = checksumBytes[0]
	f.payloadChecksum[1] = checksumBytes[1]
	f.payloadChecksum[2] = checksumBytes[2]
	f.payloadChecksum[3] = checksumBytes[3]
}

func (f *Frame) GetPayload() []byte {
	return f.payload
}

func uint32ToBytes(val uint32) []byte {
	target := make([]byte, 4)
	binary.LittleEndian.PutUint32(target, val)
	return target
}

func unmarshal[T proto.Message](payload []byte) (T, error) {
	var t T
	err := proto.Unmarshal(payload, t)
	return t, err
}

func sendFrame(frame *Frame, logger zerolog.Logger, stream network.Stream) {
	respBuf, err := frame.Marshal()
	if err != nil {
		// todo (sienna): add error handling
		logger.Error().Err(err).Msg("error marshaling frame")
	}

	// set the write deadline
	deadline := time.Now().Add(RaftControlRPCWriteTimeout)
	if err := stream.SetWriteDeadline(deadline); err != nil {
		_ = stream.Reset()
	}

	_, err = stream.Write(respBuf)
	if err != nil {
		// todo (sienna): add error handling
		logger.Error().Err(err).Msg("error sending frame")
	}
}
