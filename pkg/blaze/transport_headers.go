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

	"github.com/cockroachdb/errors"
)

const (
	requestHeaderSize        = 18
	raftType          uint16 = 100
	snapshotType      uint16 = 200
)

type requestHeader struct {
	method uint16
	size   uint64
	crc    uint32
}

func (h *requestHeader) encode(buf []byte) []byte {
	if len(buf) < requestHeaderSize {
		panic("input buf too small")
	}

	// set the method type and size of payload
	binary.LittleEndian.PutUint16(buf, h.method)
	binary.LittleEndian.PutUint64(buf[2:], h.size)

	binary.LittleEndian.PutUint32(buf[10:], 0)
	binary.LittleEndian.PutUint32(buf[14:], h.crc)

	v := crc32.ChecksumIEEE(buf[:requestHeaderSize])
	binary.LittleEndian.PutUint32(buf[10:], v)

	return buf[:requestHeaderSize]
}

func (h *requestHeader) decode(buf []byte) error {
	if len(buf) < requestHeaderSize {
		return errors.New("input buffer too small")
	}

	incoming := binary.LittleEndian.Uint32(buf[10:])
	binary.LittleEndian.PutUint32(buf[10:], 0)

	expected := crc32.ChecksumIEEE(buf[:requestHeaderSize])
	if incoming != expected {
		return errors.New("invalid crc checksum")
	}

	binary.LittleEndian.PutUint32(buf[10:], incoming)
	method := binary.LittleEndian.Uint16(buf)

	if method != raftType && method != snapshotType {
		return errors.New("invalid method type")
	}

	h.method = method
	h.size = binary.LittleEndian.Uint64(buf[2:])
	h.crc = binary.LittleEndian.Uint32(buf[14:])

	return nil
}
