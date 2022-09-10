/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package kv

import (
	"bufio"
	"bytes"

	"github.com/cockroachdb/errors"
)

const (
	splitByte uint8 = 47
)

func getAccountAndBucket(key []byte) ([]byte, []byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(key))

	// Call Split to specify that we want to Scan each individual byte.
	scanner.Split(bufio.ScanBytes)

	account := make([]byte, 0)
	bucket := make([]byte, 0)

	sepsFound := 0
	idx := 0
	for scanner.Scan() {
		b := scanner.Bytes()

		// if the first character isn't a `/`, throw an error
		if idx == 0 && b[0] != splitByte {
			return nil, nil, errors.New("key must start with `/`")
		}

		// we found the correct char, now skip ahead
		if idx == 0 && b[0] == splitByte {
			idx++
			sepsFound++
			continue
		}

		// found the bucket separator, have to account for single digit accounts numbers
		if idx >= 2 && b[0] == splitByte {
			sepsFound++
			continue
		}

		if idx >= 5 && b[0] == splitByte {
			sepsFound++
		}

		// check to see if we're still in the accountId section
		if idx >= 0 && sepsFound < 2 {
			account = append(account, b[0])
			idx++
			continue
		}

		// now we're in the bucket key
		if sepsFound == 2 && b[0] != splitByte {
			bucket = append(bucket, b[0])
			idx++
			continue
		}

		// we're out of the bucket key, so break the loop
		if b[0] == splitByte {
			idx++
			break
		}

		// if we still haven't seen the closing bucket key, fuck it
		if idx > 256 {
			return nil, nil, errors.New("bucket name is too long")
		}
	}
	return account, bucket, nil
}
