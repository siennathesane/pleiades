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
	"reflect"
	"testing"
)

func TestGetAccountAndBucket(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name        string
		args        args
		wantAccount []byte
		wantBucket  []byte
		wantErr     bool
	}{
		{"valid",
			args{key: []byte("/1234/bucket/key")},
			[]byte("1234"),
			[]byte("bucket"),
			false,
		},
		{"valid-single-account-number",
			args{key: []byte("/1/bucket/key")},
			[]byte("1"),
			[]byte("bucket"),
			false,
		},
		{"valid-single-bucket-name",
			args{key: []byte("/1/b/key")},
			[]byte("1"),
			[]byte("b"),
			false,
		},
		{"missing-first-slash",
			args{key: []byte("1234/bucket/key")},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccount, gotBucket, err := getAccountAndBucket(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountAndBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAccount, tt.wantAccount) {
				t.Errorf("GetAccountAndBucket() gotAccount = %v, want %v", gotAccount, tt.wantAccount)
			}
			if !reflect.DeepEqual(gotBucket, tt.wantBucket) {
				t.Errorf("GetAccountAndBucket() gotBucket = %v, want %v", gotBucket, tt.wantBucket)
			}
		})
	}
}
