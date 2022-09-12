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
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
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
	//viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := NewBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")
}

func (t *BBoltFsmTestSuite) TestBBoltStateMachineOpen() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := NewBBoltStateMachine(t.shardId, t.replicaId)
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

	fsm := NewBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")

	idx, err := fsm.Open(make(chan struct{}))
	t.Require().NoError(err, "there must not be an error when opening the fsm")
	t.Require().Equal(uint64(0), idx, "the index must be zero as it's a non-existent fsm")

	err = fsm.Close()
	t.Require().NoError(err, "there must not be an error when closing the store")

	fsm = nil
	fsm = NewBBoltStateMachine(t.shardId, t.replicaId)
	t.Require().NotNil(fsm, "the fsm must not be nil")
	t.Assert().Panics(func() {
		_ = fsm.Close()
	}, "the fsm should panic when there is no store handle")
}

func (t *BBoltFsmTestSuite) TestBBoltStateMachineUpdate() {
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	fsm := NewBBoltStateMachine(t.shardId, t.replicaId)
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

	createAccountRequest := &database.KVStoreWrapper_CreateAccountRequest{
		CreateAccountRequest: &database.CreateAccountRequest{
			AccountId: testAccountId,
			Owner:     testOwner,
		},
	}

	createAccountEntry := &database.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     database.KVStoreWrapper_CREATE_ACCOUNT_REQUEST,
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
	resp := &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(database.KVStoreWrapper_CREATE_ACCOUNT_REPLY, resp.Typ, "the response type must be create account reply")

	createAccountResp := resp.GetCreateAccountReply()
	t.Require().NotNil(createAccountResp, "the account response must not be nil")
	t.Require().NotEmpty(createAccountResp.AccountDescriptor, "the account descriptor must not be nil")
	t.Require().Equal(testAccountId, createAccountResp.GetAccountDescriptor().GetAccountId(), "the account ids must match")

	entries = []statemachine.Entry{}

	createBucketRequest := &database.KVStoreWrapper_CreateBucketRequest{
		CreateBucketRequest: &database.CreateBucketRequest{
			AccountId: testAccountId,
			Name:      testBucketId,
			Owner:     testOwner,
		},
	}

	createBucketEntry := &database.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     database.KVStoreWrapper_CREATE_BUCKET_REQUEST,
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
	resp = &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(database.KVStoreWrapper_CREATE_BUCKET_REPLY, resp.Typ, "the response type must be create account reply")

	createBucketResp := resp.GetCreateBucketReply()
	t.Require().NotNil(createBucketResp, "the bucket response must not be nil")
	t.Require().NotEmpty(createBucketResp.BucketDescriptor, "the bucket descriptor must not be nil")
	t.Require().Equal(testOwner, createBucketResp.GetBucketDescriptor().GetOwner(), "the account ids must match")

	entries = []statemachine.Entry{}

	testPutKeyValue, _ := utils.RandomBytes(128)
	putKeyRequest := &database.KVStoreWrapper_PutKeyRequest{
		PutKeyRequest: &database.PutKeyRequest{
			AccountId:  testAccountId,
			BucketName: testBucketId,
			KeyValuePair: &database.KeyValue{
				Key:            "test-key",
				CreateRevision: 0,
				ModRevision:    0,
				Version:        1,
				Value:          testPutKeyValue,
				Lease:          0,
			},
		},
	}

	putKeyEntry := &database.KVStoreWrapper{
		Account: testAccountId,
		Bucket:  testBucketId,
		Typ:     database.KVStoreWrapper_PUT_KEY_REQUEST,
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
	resp = &database.KVStoreWrapper{}
	err = resp.UnmarshalVT(smCmdResponse.Result.Data)
	t.Require().NoError(err, "there must not be an error when unmarshaling the cmd response")
	t.Require().Equal(database.KVStoreWrapper_PUT_KEY_REPLY, resp.Typ, "the response type must be create account reply")

	putKeyReply := resp.GetPutKeyReply()
	t.Require().NotNil(putKeyReply, "the put key response must not be nil")

	// now we work backwards, to delete everything, but we're executing the commands as an array, so we can batch
	// update things for speed. we're basically undoing everything we just did, but in reverse

	entries = make([]statemachine.Entry, 3)

	// delete the key, to ensure it's gone
	entries[0] = statemachine.Entry{
		Index: 0,
		Cmd: func() []byte {
			req := &database.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     database.KVStoreWrapper_DELETE_KEY_REQUEST,
				Payload: &database.KVStoreWrapper_DeleteKeyRequest{
					DeleteKeyRequest: &database.DeleteKeyRequest{
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
			req := &database.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     database.KVStoreWrapper_DELETE_BUCKET_REQUEST,
				Payload: &database.KVStoreWrapper_DeleteBucketRequest{
					DeleteBucketRequest: &database.DeleteBucketRequest{
						AccountId:  testAccountId,
						Name: testBucketId,
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
			req := &database.KVStoreWrapper{
				Account: testAccountId,
				Bucket:  testBucketId,
				Typ:     database.KVStoreWrapper_DELETE_ACCOUNT_REQUEST,
				Payload: &database.KVStoreWrapper_DeleteAccountRequest{
					DeleteAccountRequest: &database.DeleteAccountRequest{
						AccountId:  testAccountId,
						Owner: testOwner,
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

//func (t *BBoltFsmTestSuite) TestPrepareSnapshot() {
//	testOpts := &bbolt.Options{
//		Timeout:         0,
//		NoGrowSync:      false,
//		NoFreelistSync:  false,
//		FreelistType:    bbolt.FreelistMapType,
//		ReadOnly:        false,
//		InitialMmapSize: 0,
//		PageSize:        0,
//		NoSync:          false,
//		OpenFile:        nil,
//		Mlock:           false,
//	}
//
//	fsm := NewBBoltStateMachine(1, 1, t.T().TempDir(), testOpts)
//	empty, err := fsm.PrepareSnapshot()
//	t.Require().Empty(empty, "this must be a noop")
//	t.Require().NoError(err, "there must be no error")
//}
//
//func (t *TestBBoltFsm) TestSnapshotLifecycle() {
//	testOpts := &bbolt.Options{
//		Timeout:         0,
//		NoGrowSync:      false,
//		NoFreelistSync:  false,
//		FreelistType:    bbolt.FreelistMapType,
//		ReadOnly:        false,
//		InitialMmapSize: 0,
//		PageSize:        0,
//		NoSync:          false,
//		OpenFile:        nil,
//		Mlock:           false,
//	}
//
//	fsm := NewBBoltStateMachine(1, 1, t.T().TempDir(), testOpts)
//	index, err := fsm.Open(make(<-chan struct{}))
//	t.Require().NoError(err, "there must not be an error when opening the database")
//	t.Require().Equal(uint64(0), index, "the index must equal as there are no records")
//
//	rootPrn := &PleiadesResourceName{
//		Partition:    GlobalPartition,
//		Service:      Pleiades,
//		Region:       GlobalRegion,
//		AccountId:    fsm2.testAccountKey,
//		ResourceType: Bucket,
//		ResourceId:   "test-bucket",
//	}
//
//	var testKvps []db.KeyValue
//
//	for i := 0; i < 3; i++ {
//		kvp := db.KeyValue{
//			Key:            []byte(fmt.Sprintf("%s/test-key-%d", rootPrn.ToFsmRootPath("test-bucket"), i)),
//			Value:          []byte(fmt.Sprintf("test-value-%d", i)),
//			CreateRevision: 0,
//			ModRevision:    0,
//			Version:        1,
//			Lease:          0,
//		}
//
//		testKvps = append(testKvps, kvp)
//	}
//
//	testUpdates := make([]statemachine.Entry, 0)
//	for idx := range testKvps {
//		payload, err := testKvps[idx].MarshalVT()
//		t.Require().NoError(err, "there must not be an error encoding the message")
//
//		testUpdates = append(testUpdates, statemachine.Entry{
//			Index:  uint64(idx),
//			Cmd:    payload,
//			Result: statemachine.Result{},
//		})
//	}
//
//	var endingIndex []statemachine.Entry
//	t.Require().NotPanics(func() {
//		endingIndex, err = fsm.Update(testUpdates)
//	})
//	t.Require().NoError(err, "there must not be an error delivering updates")
//	t.Require().Equal(
//		testUpdates[len(testUpdates)-1].Index,
//		endingIndex[len(endingIndex)-1].Index,
//		fmt.Sprintf("the ending index must be %d", testUpdates[len(testUpdates)-1].Index))
//
//	// todo (sienna): replace local buffer with net.Pipe() for better test
//	var buf bytes.Buffer
//	t.Require().NotPanics(func() {
//		err = fsm.SaveSnapshot(context.Background(), &buf, make(<-chan struct{}))
//	}, "saving the snapshot to the buffer must not panic")
//	t.Require().NoError(err, "there must not be an error when saving the snapshot")
//
//	t.Require().NotPanics(func() {
//		err = fsm.RecoverFromSnapshot(&buf, make(<-chan struct{}))
//	}, "recovering the snapshot must not panic")
//	t.Require().NoError(err, "there must not be an error when recovering from a snapshot")
//
//	dbPath := fsm.dbPath(true)
//	bdb, err := bbolt.Open(dbPath, os.FileMode(fsm2.dbFileModeVal), nil)
//	t.Require().NoError(err)
//
//	// todo (sienna): fix this to ensure the finalKvp is stored, not just the value
//	var target []byte
//	t.Require().NoError(bdb.Update(func(tx *bbolt.Tx) error {
//		// rootPrn.ToFsmRootPath("test-bucket") + "/" + "test-key-2"
//		bucketHierarchy := strings.Split(rootPrn.ToFsmRootPath("test-bucket")+"/"+"test-key-2", "/")[1:]
//		parentBucketName := bucketHierarchy[0]
//		childBucketNames := bucketHierarchy[1:]
//		parentBucket := tx.Bucket([]byte(parentBucketName))
//		resChan := make(chan []byte, 1)
//		if err := keyOp(parentBucket, childBucketNames, make([]byte, 0), get, &resChan); err != nil {
//			return err
//		}
//		target = <-resChan
//		close(resChan)
//		return nil
//	}))
//
//	finalKvp := db.KeyValue{}
//	err = finalKvp.UnmarshalVT(target)
//	t.Require().NoError(err, "there must not be an error unmarshalling the key")
//
//	testKey := testKvps[len(testKvps)-1].Key
//	returnedKey := finalKvp.Key
//
//	t.Require().Equal(testKey, returnedKey, "the serialized result must match the initial value")
//	t.Require().Equal(testKvps[len(testKvps)-1].Version, finalKvp.Version, "the serialized result must match the initial value")
//	t.Require().Equal(testKvps[len(testKvps)-1].ModRevision, finalKvp.ModRevision, "the serialized result must match the initial value")
//	t.Require().Equal(testKvps[len(testKvps)-1].Lease, finalKvp.Lease, "the serialized result must match the initial value")
//	t.Require().Equal(testKvps[len(testKvps)-1].CreateRevision, finalKvp.CreateRevision, "the serialized result must match the initial value")
//
//	testValue := testKvps[len(testKvps)-1].Value
//	returnedValue := finalKvp.Value
//
//	t.Require().Equal(testValue, returnedValue, "the serialized result must match the initial value")
//}
//
//func (t *BBoltFsmTestSuite) TestLookup() {
//	testOpts := &bbolt.Options{
//		Timeout:         0,
//		NoGrowSync:      false,
//		NoFreelistSync:  false,
//		FreelistType:    bbolt.FreelistMapType,
//		ReadOnly:        false,
//		InitialMmapSize: 0,
//		PageSize:        0,
//		NoSync:          false,
//		OpenFile:        nil,
//		Mlock:           false,
//	}
//
//	fsm := NewBBoltStateMachine(1, 1, t.T().TempDir(), testOpts)
//	index, err := fsm.Open(make(<-chan struct{}))
//	t.Require().NoError(err, "there must not be an error when opening the database")
//	t.Require().Equal(uint64(0), index, "the index must equal as there are no records")
//
//	rootPrn := &PleiadesResourceName{
//		Partition:    GlobalPartition,
//		Service:      Pleiades,
//		Region:       GlobalRegion,
//		AccountId:    fsm2.testAccountKey,
//		ResourceType: Bucket,
//		ResourceId:   "test-bucket",
//	}
//
//	var testKvps []db.KeyValue
//
//	for i := 0; i < 3; i++ {
//		kvp := db.KeyValue{
//			Key:            []byte(fmt.Sprintf("%s/test-key-%d", rootPrn.ToFsmRootPath("test-bucket"), i)),
//			Value:          []byte(fmt.Sprintf("test-value-%d", i)),
//			CreateRevision: 0,
//			ModRevision:    0,
//			Version:        1,
//			Lease:          0,
//		}
//
//		testKvps = append(testKvps, kvp)
//	}
//
//	testUpdates := make([]statemachine.Entry, 0)
//	for idx := range testKvps {
//		marshalled, err := testKvps[idx].MarshalVT()
//		t.Require().NoError(err, "there must not be an error marshalling the message")
//
//		testUpdates = append(testUpdates, statemachine.Entry{
//			Index:  uint64(idx),
//			Cmd:    marshalled,
//			Result: statemachine.Result{},
//		})
//	}
//
//	var endingIndex []statemachine.Entry
//	t.Require().NotPanics(func() {
//		endingIndex, err = fsm.Update(testUpdates)
//	})
//	t.Require().NoError(err, "there must not be an error delivering updates")
//	t.Require().Equal(
//		testUpdates[len(testUpdates)-1].Index,
//		endingIndex[len(endingIndex)-1].Index,
//		fmt.Sprintf("the ending index must be %d", testUpdates[len(testUpdates)-1].Index))
//
//	val, err := fsm.Lookup(testUpdates[len(testUpdates)-1].Cmd)
//	t.Require().NoError(err, "there must not be an error when calling lookup")
//
//	var casted db.KeyValue
//	t.Require().NotPanics(func() {
//		casted = val.(db.KeyValue)
//	}, "casting the lookup value must not panic")
//
//	expectedValue := testKvps[len(testUpdates)-1].Value
//	castedValue := casted.Value
//	t.Require().Equal(expectedValue, castedValue, "the found value must be identical")
//
//	expectedKey := testKvps[len(testUpdates)-1].Key
//	castedKey := casted.Key
//	t.Require().Equal(expectedKey, castedKey, "the found value must be identical")
//
//	t.Require().Equal(testKvps[len(testUpdates)-1].Lease, casted.Lease, "the found value must be identical")
//	t.Require().Equal(testKvps[len(testUpdates)-1].CreateRevision, casted.CreateRevision, "the found value must be identical")
//	t.Require().Equal(testKvps[len(testUpdates)-1].Version, casted.Version, "the found value must be identical")
//	t.Require().Equal(testKvps[len(testUpdates)-1].ModRevision, casted.ModRevision, "the found value must be identical")
//}
//
//func (t *BBoltFsmTestSuite) TestSync() {
//	testOpts := &bbolt.Options{
//		Timeout:         0,
//		NoGrowSync:      false,
//		NoFreelistSync:  false,
//		FreelistType:    bbolt.FreelistMapType,
//		ReadOnly:        false,
//		InitialMmapSize: 0,
//		PageSize:        0,
//		NoSync:          false,
//		OpenFile:        nil,
//		Mlock:           false,
//	}
//
//	fsm := NewBBoltStateMachine(1, 1, t.T().TempDir(), testOpts)
//	index, err := fsm.Open(make(<-chan struct{}))
//	t.Require().NoError(err, "there must not be an error when opening the database")
//	t.Require().Equal(uint64(0), index, "the index must equal as there are no records")
//
//	rootPrn := &PleiadesResourceName{
//		Partition:    GlobalPartition,
//		Service:      Pleiades,
//		Region:       GlobalRegion,
//		AccountId:    fsm2.testAccountKey,
//		ResourceType: Bucket,
//		ResourceId:   "test-bucket",
//	}
//
//	var testKvps []db.KeyValue
//
//	for i := 0; i < 3; i++ {
//		kvp := db.KeyValue{
//			Key:            []byte(fmt.Sprintf("%s/test-key-%d", rootPrn.ToFsmRootPath("test-bucket"), i)),
//			Value:          []byte(fmt.Sprintf("test-value-%d", i)),
//			CreateRevision: 0,
//			ModRevision:    0,
//			Version:        1,
//			Lease:          0,
//		}
//
//		testKvps = append(testKvps, kvp)
//	}
//
//	testUpdates := make([]statemachine.Entry, 0)
//	for idx := range testKvps {
//		marshalled, err := testKvps[idx].MarshalVT()
//		t.Require().NoError(err, "there must not be an error marshalling the message")
//
//		testUpdates = append(testUpdates, statemachine.Entry{
//			Index:  uint64(idx),
//			Cmd:    marshalled,
//			Result: statemachine.Result{},
//		})
//	}
//
//	var endingIndex []statemachine.Entry
//	t.Require().NotPanics(func() {
//		endingIndex, err = fsm.Update(testUpdates)
//	})
//	t.Require().NoError(err, "there must not be an error delivering updates")
//	t.Require().Equal(
//		testUpdates[len(testUpdates)-1].Index,
//		endingIndex[len(endingIndex)-1].Index,
//		fmt.Sprintf("the ending index must be %d", testUpdates[len(testUpdates)-1].Index))
//
//	t.Require().NotPanics(func() {
//		err = fsm.Sync()
//	})
//	t.Require().NoError(err, "there must not be an error when syncing bbolt to disk")
//}
