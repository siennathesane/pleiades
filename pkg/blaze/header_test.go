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
	"hash/crc32"
	"testing"

	transportv1 "github.com/mxplusb/pleiades/pkg/api/v1"
)

func TestHeaderSize(t *testing.T) {
	header := &transportv1.Header{
		Size:     123,
		Checksum: crc32.ChecksumIEEE([]byte("really actually test")),
	}
	t.Logf("header size: %d", header.SizeVT())

	marshalled, _ := header.MarshalVT()
	t.Logf("marshalled header size: %d", len(marshalled))
}
