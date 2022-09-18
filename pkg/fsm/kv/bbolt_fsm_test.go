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
	"context"
	"encoding/binary"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

func TestBBoltFsm(t *testing.T) {
	suite.Run(t, new(BBoltFsmTestSuite))
}

type BBoltFsmTestSuite struct {
	suite.Suite
	logger    zerolog.Logger
	shardId   uint64
	replicaId uint64
}

func (t *BBoltFsmTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.shardId = rand.Uint64()
	t.replicaId = rand.Uint64()
}

func (t *BBoltFsmTestSuite) TestNewBBoltStateMachine() {
	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")
}

func (t *BBoltFsmTestSuite) TestBBoltStateMachineOpen() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	err = fsm.store.db.Update(func(tx *bbolt.Tx) error {
		internalBucket, err := tx.CreateBucketIfNotExists([]byte(monotonicLogBucket))
		if err != nil {
			return err
		}

		indexBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(indexBuf, 1)

		return internalBucket.Put([]byte(monotonicLogKey), indexBuf)
	})
	t.Require().NoError(err, "there must not be an error when manually setting the index")

	err = fsm.store.close()
	t.Require().NoError(err, "there must not be an error when manually closing the store")

	idx, err = fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(1), idx, "the index must be one as it's an existing fsm")
}

func (t *BBoltFsmTestSuite) TestBBoltStateMachineClose() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	err = fsm.Close()
	t.Require().NoError(err, "there must not be an error when closing the store")

	fsm = nil
	fsm = newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")
	t.Assert().Panics(func() {
		_ = fsm.Close()
	}, "the fsm should panic when there is no store handle")
}

func (t *BBoltFsmTestSuite) TestBBoltStateMachineUpdate() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	// handle no updates
	entries, err := fsm.Update([]statemachine.Entry{})
	t.Require().NoError(err, "there must not be an error inputting no entries")
	t.Require().Empty(entries, "the entries list must be empty")

	testAccountId := rand.Uint64()
	testBucketId := utils.RandomString(10)
	testOwner := "test@test.com"

	createAccountRequest := &kvstorev1.KVStoreWrapper_CreateAccountRequest{
		CreateAccountRequest: &kvstorev1.CreateAccountRequest{
			AccountId: testAccountId,
			Owner:     testOwner,
		},
	}

	createAccountEntry := &kvstorev1.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_ACCOUNT_REQUEST,
		Payload: createAccountRequest,
	}
	smCreateAccountRequestPayload, err := createAccountEntry.MarshalVT()
	t.Require().NoError(err, "there must not be an error when marshalling the the request")

	entries = append(entries, statemachine.Entry{
		Index: 0,
		Cmd:   smCreateAccountRequestPayload,
	})

	appliedEntries, err := fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error creating an account through the state machine")

	smCmdResponse := appliedEntries[0]
	resp := &kvstorev1.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_ACCOUNT_REPLY, resp.Typ, "the response type must be create account reply")

	createAccountResp := resp.GetCreateAccountReply()
	t.Require().NotNil(createAccountResp, "the account response must not be nil")
	t.Require().NotEmpty(createAccountResp.AccountDescriptor, "the account descriptor must not be nil")
	t.Require().Equal(testAccountId, createAccountResp.GetAccountDescriptor().GetAccountId(), "the account ids must match")

	entries = []statemachine.Entry{}

	createBucketRequest := &kvstorev1.KVStoreWrapper_CreateBucketRequest{
		CreateBucketRequest: &kvstorev1.CreateBucketRequest{
			AccountId: testAccountId,
			Name:      testBucketId,
			Owner:     testOwner,
		},
	}

	createBucketEntry := &kvstorev1.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_BUCKET_REQUEST,
		Payload: createBucketRequest,
	}
	smCreateBucketRequestPayload, err := createBucketEntry.MarshalVT()
	t.Require().NoError(err, "there must not be an error when marshalling the the request")

	entries = append(entries, statemachine.Entry{
		Index: 0,
		Cmd:   smCreateBucketRequestPayload,
	})

	appliedEntries, err = fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error when creating a bucket")
	t.Require().NotEmpty(appliedEntries, "the applied entries must not be empty")

	smCmdResponse = appliedEntries[0]
	resp = &kvstorev1.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_BUCKET_REPLY, resp.Typ, "the response type must be create account reply")

	createBucketResp := resp.GetCreateBucketReply()
	t.Require().NotNil(createBucketResp, "the bucket response must not be nil")
	t.Require().NotEmpty(createBucketResp.BucketDescriptor, "the bucket descriptor must not be nil")
	t.Require().Equal(testOwner, createBucketResp.GetBucketDescriptor().GetOwner(), "the account ids must match")

	entries = []statemachine.Entry{}

	testPutKeyValue, _ := utils.RandomBytes(128)
	putKeyRequest := &kvstorev1.KVStoreWrapper_PutKeyRequest{
		PutKeyRequest: &kvstorev1.PutKeyRequest{
			AccountId:  testAccountId,
			BucketName: testBucketId,
			KeyValuePair: &kvstorev1.KeyValue{
				Key:            "test-key",
				CreateRevision: 0,
				ModRevision:    0,
				Version:        1,
				Value:          testPutKeyValue,
				Lease:          0,
			},
		},
	}

	putKeyEntry := &kvstorev1.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_PUT_KEY_REQUEST,
		Payload: putKeyRequest,
	}
	smCreateBucketRequestPayload, err = putKeyEntry.MarshalVT()
	t.Require().NoError(err, "there must not be an error when marshalling the the request")

	entries = append(entries, statemachine.Entry{
		Index: 0,
		Cmd:   smCreateBucketRequestPayload,
	})

	appliedEntries, err = fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error when creating a bucket")
	t.Require().NotEmpty(appliedEntries, "the applied entries must not be empty")

	smCmdResponse = appliedEntries[0]
	resp = &kvstorev1.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(kvstorev1.KVStoreWrapper_REQUEST_TYPE_PUT_KEY_REPLY, resp.Typ, "the response type must be create account reply")

	putKeyReply := resp.GetPutKeyReply()
	t.Require().NotNil(putKeyReply, "the put key response must not be nil")

	// now we work backwards, to delete everything, but we're executing the commands as an array, so we can batch
	// update things for speed. we're basically undoing everything we just did, but in reverse

	entries = make([]statemachine.Entry, 3)

	// delete the key, to ensure it's gone
	entries[0] = statemachine.Entry{
		Index: 0,
		Cmd: func() []byte {
			req := &kvstorev1.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_DELETE_KEY_REQUEST,
				Payload: &kvstorev1.KVStoreWrapper_DeleteKeyRequest{
					DeleteKeyRequest: &kvstorev1.DeleteKeyRequest{
						AccountId:  testAccountId,
						BucketName: testBucketId,
						Key:        putKeyRequest.PutKeyRequest.GetKeyValuePair().GetKey(),
					},
				},
			}
			serial, _ := req.MarshalVT()
			return serial
		}(),
		Result: statemachine.Result{},
	}

	// delete the bucket
	entries[1] = statemachine.Entry{
		Index: 0,
		Cmd: func() []byte {
			req := &kvstorev1.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_DELETE_BUCKET_REQUEST,
				Payload: &kvstorev1.KVStoreWrapper_DeleteBucketRequest{
					DeleteBucketRequest: &kvstorev1.DeleteBucketRequest{
						AccountId: testAccountId,
						Name:      testBucketId,
					},
				},
			}
			serial, _ := req.MarshalVT()
			return serial
		}(),
		Result: statemachine.Result{},
	}

	// delete the account
	entries[2] = statemachine.Entry{
		Index: 0,
		Cmd: func() []byte {
			req := &kvstorev1.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_DELETE_ACCOUNT_REQUEST,
				Payload: &kvstorev1.KVStoreWrapper_DeleteAccountRequest{
					DeleteAccountRequest: &kvstorev1.DeleteAccountRequest{
						AccountId: testAccountId,
						Owner:     testOwner,
					},
				},
			}
			serial, _ := req.MarshalVT()
			return serial
		}(),
		Result: statemachine.Result{},
	}

	appliedEntries, err = fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error when applying kv store updates")
	t.Require().NotEmpty(appliedEntries, "there must be response entries")
	t.Require().Equal(3, len(appliedEntries), "there must be 4 response entries")
}

func (t *BBoltFsmTestSuite) TestSnapshotLifecycle() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	testAccountId := rand.Uint64()
	testBucketId := utils.RandomString(10)
	testOwner := "test@test.com"

	testPutKeyValue, _ := utils.RandomBytes(128)
	testKvp := &kvstorev1.KeyValue{
		Key:            "test-key",
		CreateRevision: 0,
		ModRevision:    0,
		Version:        1,
		Value:          testPutKeyValue,
		Lease:          0,
	}

	entries := []statemachine.Entry{
		{
			Index: 0,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_ACCOUNT_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_CreateAccountRequest{
						CreateAccountRequest: &kvstorev1.CreateAccountRequest{
							AccountId: testAccountId,
							Owner:     testOwner,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
		{
			Index: 1,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_BUCKET_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_CreateBucketRequest{
						CreateBucketRequest: &kvstorev1.CreateBucketRequest{
							AccountId: testAccountId,
							Name:      testBucketId,
							Owner:     testOwner,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
		{
			Index: 2,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_PUT_KEY_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_PutKeyRequest{
						PutKeyRequest: &kvstorev1.PutKeyRequest{
							AccountId:  testAccountId,
							BucketName: testBucketId,
							KeyValuePair: testKvp,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
	}

	_, err = fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error creating the baseline db")

	empty, err := fsm.PrepareSnapshot()
	t.Require().NoError(err, "there must not be an error when preparing a snapshot")
	t.Require().Nil(empty, "the response to preparing a snapshot much be empty")

	snapshotFile := filepath.Join(t.T().TempDir(), "snapshot.db")
	file, err := os.Create(snapshotFile)
	t.Require().NoError(err, "there must not be an error when opening the temp file")

	err = fsm.SaveSnapshot(context.TODO(), file, make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when saving the snapshot")

	file.Close()

	db, err := bbolt.Open(snapshotFile, os.FileMode(484), nil)
	t.Require().NoError(err, "there must not be an error opening the snapshot file")

	err = db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)

		accountBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(accountBucket)

		bucket := accountBucket.Bucket([]byte(testBucketId))
		t.Require().NotNil(bucket, "the bucket must not be nil")

		val := bucket.Get([]byte("test-key"))
		t.Require().NotNil(val)
		t.Require().NotEmpty(val)

		target := &kvstorev1.KeyValue{}
		err := target.UnmarshalVT(val)
		t.Require().NoError(err, "there must not be an error unmarshalling the kvp")
		t.Require().Equal(testKvp.GetKey(),target.GetKey())
		t.Require().Equal(testKvp.GetValue(),target.GetValue())

		return nil
	})
}

func (t *BBoltFsmTestSuite) TestLookup() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	testAccountId := rand.Uint64()
	testBucketId := utils.RandomString(10)
	testOwner := "test@test.com"

	testPutKeyValue, _ := utils.RandomBytes(128)
	testKvp := &kvstorev1.KeyValue{
		Key:            "test-key",
		CreateRevision: 0,
		ModRevision:    0,
		Version:        1,
		Value:          testPutKeyValue,
		Lease:          0,
	}

	entries := []statemachine.Entry{
		{
			Index: 0,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_ACCOUNT_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_CreateAccountRequest{
						CreateAccountRequest: &kvstorev1.CreateAccountRequest{
							AccountId: testAccountId,
							Owner:     testOwner,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
		{
			Index: 1,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_CREATE_BUCKET_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_CreateBucketRequest{
						CreateBucketRequest: &kvstorev1.CreateBucketRequest{
							AccountId: testAccountId,
							Name:      testBucketId,
							Owner:     testOwner,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
		{
			Index: 2,
			Cmd: func() []byte {
				resp := &kvstorev1.KVStoreWrapper{
					Account: testAccountId,
					Bucket:  testBucketId,
					Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_PUT_KEY_REQUEST,
					Payload: &kvstorev1.KVStoreWrapper_PutKeyRequest{
						PutKeyRequest: &kvstorev1.PutKeyRequest{
							AccountId:  testAccountId,
							BucketName: testBucketId,
							KeyValuePair: testKvp,
						},
					},
				}
				serial, _ := resp.MarshalVT()
				return serial
			}(),
			Result: statemachine.Result{},
		},
	}

	_, err = fsm.Update(entries)
	t.Require().NoError(err, "there must not be an error creating the baseline db")

	req := &kvstorev1.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     kvstorev1.KVStoreWrapper_REQUEST_TYPE_GET_KEY_REQUEST,
		Payload: &kvstorev1.KVStoreWrapper_GetKeyRequest{
			GetKeyRequest: &kvstorev1.GetKeyRequest{
				AccountId:  testAccountId,
				BucketName: testBucketId,
				Key: testKvp.GetKey(),
			},
		},
	}
	reqPayload, _ := req.MarshalVT()

	response, err := fsm.Lookup(reqPayload)
	t.Require().NoError(err, "there must not be an error when looking up the key")

	resp := &kvstorev1.KVStoreWrapper{}
	err = resp.UnmarshalVT(response.([]byte))
	t.Require().NoError(err, "there must not be an error when unmarshalling the lookup value")
	t.Require().NotEmpty(resp)
	t.Require().Equal(kvstorev1.KVStoreWrapper_REQUEST_TYPE_GET_KEY_REPLY, resp.Typ)
	t.Require().NotNil(resp.GetGetKeyReply())
	t.Require().NotNil(resp.GetGetKeyReply().GetKeyValuePair())
	t.Require().Equal(testPutKeyValue, resp.GetGetKeyReply().GetKeyValuePair().GetValue())
}

func (t *BBoltFsmTestSuite) TestSync() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := newBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	err = fsm.Sync()
	t.Require().NoError(err, "there must not be an error syncing the fsm")
}
