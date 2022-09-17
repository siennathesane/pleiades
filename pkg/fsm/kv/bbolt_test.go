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
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

const (
	// fuzzRounds just sets how many times the fuzzer will run in test mode.
	fuzzRounds = 256
)

func TestBBolt(t *testing.T) {
	suite.Run(t, new(BBoltTestSuite))
}

type BBoltTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func (t *BBoltTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
}

func (t *BBoltTestSuite) TestNewBBoltStore() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	b, err := newBboltStore(shardId, replicaId, t.T().TempDir(), t.logger)
	t.Require().Error(err, "there must be an error when passing a bad directory")
	t.Require().Nil(b, "the bbolt store must be nil")

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err = newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")
}

func (t *BBoltTestSuite) TestCreateAccountBucket() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testOwner := "test@test.com"

	req := &kvstorev1.CreateAccountRequest{}
	resp, err := b.CreateAccountBucket(req)
	t.Require().Error(err, "there must be an error when sending an empty request")
	t.Require().Empty(resp.GetAccountDescriptor(), "the response must be empty")

	req.AccountId = testAccountId
	resp, err = b.CreateAccountBucket(req)
	t.Require().Error(err, "there must be an error when sending a partial request")
	t.Require().Empty(resp.GetAccountDescriptor(), "the response must be empty")

	req.Owner = testOwner
	resp, err = b.CreateAccountBucket(req)
	t.Require().NoError(err, "there must not be an error when sending a request")
	t.Require().NotEmpty(resp.GetAccountDescriptor(), "the response must not be empty")

	desc := resp.GetAccountDescriptor()
	t.Require().Equal(testAccountId, desc.GetAccountId(), "the account ids must equal")
	t.Require().Equal(testOwner, desc.GetOwner(), "the owner ids must equal")
	t.Require().Equal(uint64(0), desc.GetBucketCount(), "the bucket count must be zero")
	t.Require().NotNil(desc.GetCreated(), "the created timestamp must not be nil")
	t.Require().NotNil(desc.GetLastUpdated(), "the last updated timestamp must not be nil")

	expectedBucketList := make([][]byte, 0)
	t.Require().Empty(expectedBucketList, desc.GetBuckets(), "the list of buckets must be empty")

	// reach into the db and verify it was stored correctly
	foundAcctDescriptor := &kvstorev1.AccountDescriptor{}
	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)
		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		descriptor := acctBucket.Get([]byte(descriptorKey))
		t.Require().NotNil(descriptor, "the descriptor value must not be nil")

		if err := proto.Unmarshal(descriptor, foundAcctDescriptor); err != nil {
			return err
		}
		return nil
	})
	t.Require().NoError(err, "there must not be an error when peeking the account descriptor")

	// verify what was stored in the db
	t.Require().Equal(testAccountId, foundAcctDescriptor.GetAccountId(), "the account ids must equal")
	t.Require().Equal(testOwner, foundAcctDescriptor.GetOwner(), "the owner ids must equal")
	t.Require().Equal(uint64(0), foundAcctDescriptor.GetBucketCount(), "the bucket count must be zero")
	t.Require().NotNil(foundAcctDescriptor.GetCreated(), "the created timestamp must not be nil")
	t.Require().NotNil(foundAcctDescriptor.GetLastUpdated(), "the last updated timestamp must not be nil")
	t.Require().Empty(expectedBucketList, foundAcctDescriptor.GetBuckets(), "the list of buckets must be empty")
}

func (t *BBoltTestSuite) TestDeleteAccountBucket() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testOwner := "test@test.com"

	// prep the test
	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	req := &kvstorev1.DeleteAccountRequest{}

	req.AccountId = testAccountId
	resp, err := b.DeleteAccountBucket(req)
	t.Require().Error(err, "there must be an error when sending a partial request")
	t.Require().False(resp.GetOk(), "the response must not be ok")

	req.Owner = testOwner
	resp, err = b.DeleteAccountBucket(req)
	t.Require().NoError(err, "there must not be an error when sending a request")
	t.Require().True(resp.GetOk(), "the response must be ok")

	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)

		bucket := tx.Bucket(accountBuf)
		t.Require().Nil(bucket)
		return nil
	})
	t.Require().NoError(err, "there must not be an error when peeking for the deleted account bucket")

	resp, err = b.DeleteAccountBucket(&kvstorev1.DeleteAccountRequest{
		AccountId: 1234,
		Owner:     "empty",
	})
	t.Require().Error(err, "there must be an error when trying to delete an account which doesn't exist")
	t.Require().False(resp.GetOk(), "the response must not be ok")
}

func (t *BBoltTestSuite) TestCreateBucket() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testBucketName := "test-bucket"
	testOwner := "test@test.com"

	// bad request
	req := &kvstorev1.CreateBucketRequest{
		AccountId: 0,
		Name:      "",
		Owner:     "",
	}

	resp, err := b.CreateBucket(req)
	t.Require().Error(err, "there must be an error when sending an empty request")
	t.Require().Empty(resp, "the response must be blank")

	req.AccountId = testAccountId
	resp, err = b.CreateBucket(req)
	t.Require().Error(err, "there must be an error when sending a partial request")
	t.Require().Empty(resp, "the response must be blank")

	req.Name = testBucketName
	resp, err = b.CreateBucket(req)
	t.Require().Error(err, "there must be an error when sending a partial request")
	t.Require().Empty(resp, "the response must be blank")

	req.Owner = testOwner

	resp, err = b.CreateBucket(req)
	t.Require().Error(err, "there must be an error when the account bucket doesn't exist")

	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error when creating the account bucket")

	resp, err = b.CreateBucket(req)
	t.Require().NoError(err, "there must not be an error when sending a valid request")
	t.Require().NotEmpty(resp, "the response must not be blank")
	t.Require().NotEmpty(resp.GetBucketDescriptor(), "the descriptor must not be empty")

	desc := resp.GetBucketDescriptor()
	t.Require().Equal(testOwner, desc.GetOwner(), "the owners must be equal")

	acctDescriptor := &kvstorev1.AccountDescriptor{}
	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)
		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		descriptor := acctBucket.Get([]byte(descriptorKey))
		t.Require().NotNil(descriptor, "the descriptor value must not be nil")

		if err := proto.Unmarshal(descriptor, acctDescriptor); err != nil {
			return err
		}
		return nil
	})
	t.Require().NoError(err, "there must not be an error when peeking into the database")

	t.Require().Equal(uint64(1), acctDescriptor.GetBucketCount(), "the number of buckets must match")
	t.Require().Equal(testBucketName, acctDescriptor.GetBuckets()[0], "the bucket names must match")
}

func (t *BBoltTestSuite) TestDeleteBucket() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testBucketName := "test-bucket"
	testOwner := "test@test.com"

	// prep the test
	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	_, err = b.CreateBucket(&kvstorev1.CreateBucketRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
		Name:      testBucketName,
	})
	t.Require().NoError(err, "there must not be an error creating the test bucket")

	// bad request
	req := &kvstorev1.DeleteBucketRequest{
		AccountId: 0,
		Name:      "",
	}

	resp, err := b.DeleteBucket(req)
	t.Require().Error(err, "there must be an error when sending an empty request")
	t.Require().Empty(resp, "the response must be blank")

	req.AccountId = testAccountId
	resp, err = b.DeleteBucket(req)
	t.Require().Error(err, "there must be an error when sending a partial request")
	t.Require().Empty(resp, "the response must be blank")

	req.Name = testBucketName
	resp, err = b.DeleteBucket(req)
	t.Require().NoError(err, "there must not be an error when sending a partial request")
	t.Require().NotEmpty(resp, "the response must be blank")
	t.Require().True(resp.GetOk(), "the request must be ok")

	acctDescriptor := &kvstorev1.AccountDescriptor{}
	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)
		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		descriptor := acctBucket.Get([]byte(descriptorKey))
		t.Require().NotNil(descriptor, "the descriptor value must not be nil")

		if err := proto.Unmarshal(descriptor, acctDescriptor); err != nil {
			return err
		}
		return nil
	})
	t.Require().NoError(err, "there must not be an error when peeking into the database")

	t.Require().Equal(uint64(0), acctDescriptor.GetBucketCount(), "the number of buckets must match")
	t.Require().Empty(acctDescriptor.GetBuckets(), "the bucket names must empty")
}

func (t *BBoltTestSuite) TestGetKey() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testBucketName := "test-bucket"
	testOwner := "test@test.com"

	// prep the test
	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	_, err = b.CreateBucket(&kvstorev1.CreateBucketRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
		Name:      testBucketName,
	})
	t.Require().NoError(err, "there must not be an error creating the test bucket")

	now := time.Now().UnixMilli()
	expectedKvp := &kvstorev1.KeyValue{
		Key:            "test-key",
		CreateRevision: now,
		ModRevision:    now,
		Version:        1,
		Value:          []byte("test-value"),
		Lease:          0,
	}

	err = b.db.Update(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)

		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		bucket := acctBucket.Bucket([]byte(testBucketName))
		t.Require().NotNil(bucket, "the target bucket must not be nil")

		payload, err := expectedKvp.MarshalVT()
		t.Require().NoError(err, "there must not be an error marshalling the test key")

		err = bucket.Put([]byte(expectedKvp.Key), payload)
		t.Require().NoError(err, "there must not be an error when storing the test key")

		return nil
	})
	t.Require().NoError(err, "there must not be an error when setting the test key")

	// search for a non-existent key in a non-existent account in a non-existent bucket
	resp, err := b.GetKey(&kvstorev1.GetKeyRequest{
		AccountId:  1,
		BucketName: "empty",
		Key:        "no",
	})
	t.Require().Error(err, "there must be an error when search for a non-existent key in a non-existent account in a non-existent bucket")
	t.Require().Empty(resp.GetKeyValuePair(), "the payload must be empty")

	// search for a non-existent key in a non-existent bucket
	resp, err = b.GetKey(&kvstorev1.GetKeyRequest{
		AccountId:  testAccountId,
		BucketName: "empty",
		Key:        "no",
	})
	t.Require().Error(err, "there must be an error when search for a non-existent key in a non-existent bucket")
	t.Require().Empty(resp.GetKeyValuePair(), "the payload must be empty")

	// search for a non-existent key
	resp, err = b.GetKey(&kvstorev1.GetKeyRequest{
		AccountId:  testAccountId,
		BucketName: testBucketName,
		Key:        "no",
	})
	t.Require().Error(err, "there must be an error when search for a non-existent key")
	t.Require().Empty(resp.GetKeyValuePair(), "the payload must be empty")

	// search for the target key
	resp, err = b.GetKey(&kvstorev1.GetKeyRequest{
		AccountId:  testAccountId,
		BucketName: testBucketName,
		Key:        expectedKvp.GetKey(),
	})
	t.Require().NoError(err, "there must not be an error when searching for an existing key")
	t.Require().NotEmpty(resp.GetKeyValuePair(), "the payload must not be empty")

	foundKvp := resp.GetKeyValuePair()
	t.Require().Equal(expectedKvp.GetKey(), foundKvp.GetKey(), "the keys must match")
	t.Require().Equal(expectedKvp.GetValue(), foundKvp.GetValue(), "the values must match")
	t.Require().Equal(expectedKvp.GetModRevision(), foundKvp.GetModRevision(), "the modify revisions must match")
	t.Require().Equal(expectedKvp.GetCreateRevision(), foundKvp.GetCreateRevision(), "the create revisions must match")
	t.Require().Equal(expectedKvp.GetVersion(), foundKvp.GetVersion(), "the versions must match")
	t.Require().Equal(expectedKvp.GetLease(), foundKvp.GetLease(), "the leases must match")
}

func (t *BBoltTestSuite) TestPutKey() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testBucketName := "test-bucket"
	testOwner := "test@test.com"

	// prep the test
	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	_, err = b.CreateBucket(&kvstorev1.CreateBucketRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
		Name:      testBucketName,
	})
	t.Require().NoError(err, "there must not be an error creating the test bucket")

	expectedRequest := &kvstorev1.PutKeyRequest{}

	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error putting an empty key")

	expectedRequest.AccountId = testAccountId
	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error putting a partial payload")

	expectedRequest.BucketName = testBucketName
	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error putting a partial payload")

	expectedKvp := &kvstorev1.KeyValue{
		Key:            "",
		CreateRevision: 0,
		ModRevision:    0,
		Version:        0,
		Value:          nil,
		Lease:          0,
	}

	expectedRequest.KeyValuePair = expectedKvp
	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error putting an empty kvp")

	expectedKvp.Key = "test-key"
	_, err = b.PutKey(expectedRequest)
	t.Require().NoError(err, "there must not be an error putting an empty value")

	expectedKvp.Version = 1
	_, err = b.PutKey(expectedRequest)
	t.Require().NoError(err, "there must be an error putting a partial kvp")

	expectedKvp.Version = 0
	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error putting a kvp with an older version")

	expectedKvp.Version = 1
	expectedKvp.Value = []byte(utils.RandomString(utils.RandomInt(0, 64)))
	_, err = b.PutKey(expectedRequest)
	t.Require().Error(err, "there must be an error trying to overwrite a key with the same value")

	expectedKvp.Version = 2
	expectedKvp.Value = []byte(utils.RandomString(utils.RandomInt(0, 64)))
	_, err = b.PutKey(expectedRequest)
	t.Require().NoError(err, "there must not be an error trying to overwrite a key with a new version")

	foundKvp := &kvstorev1.KeyValue{}
	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)

		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		bucket := acctBucket.Bucket([]byte(testBucketName))
		t.Require().NotNil(bucket, "the target bucket must not be nil")

		payload := bucket.Get([]byte(expectedKvp.Key))
		t.Require().NotNil(payload, "the expected kvp must not be empty")

		return foundKvp.UnmarshalVT(payload)
	})
	t.Require().NoError(err, "there must not be an error peeking into the database")

	t.Require().Equal(expectedKvp.GetKey(), foundKvp.GetKey(), "the keys must match")
	t.Require().Equal(expectedKvp.GetValue(), foundKvp.GetValue(), "the values must match")
	t.Require().Equal(expectedKvp.GetModRevision(), foundKvp.GetModRevision(), "the modify revisions must match")
	t.Require().Equal(expectedKvp.GetCreateRevision(), foundKvp.GetCreateRevision(), "the create revisions must match")
	t.Require().Equal(expectedKvp.GetVersion(), foundKvp.GetVersion(), "the versions must match")
	t.Require().Equal(expectedKvp.GetLease(), foundKvp.GetLease(), "the leases must match")
}

func (t *BBoltTestSuite) TestDeleteKey() {
	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)

	path := filepath.Join(t.T().TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, t.logger)
	t.Require().NoError(err, "there must be an error when passing a bad directory")
	t.Require().NotNil(b, "the bbolt store must be nil")

	testAccountId := rand.Uint64()
	testBucketName := "test-bucket"
	testOwner := "test@test.com"

	// prep the test
	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	_, err = b.CreateBucket(&kvstorev1.CreateBucketRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
		Name:      testBucketName,
	})
	t.Require().NoError(err, "there must not be an error creating the test bucket")

	expectedKvp := &kvstorev1.KeyValue{
		Key:            "test-key",
		Value:          []byte("test-value"),
	}

	expectedRequest := &kvstorev1.PutKeyRequest{
		AccountId: testAccountId,
		BucketName: testBucketName,
		KeyValuePair: expectedKvp,
	}

	_, err = b.PutKey(expectedRequest)
	t.Require().NoError(err, "there must not be an error putting an empty kvp")

	resp, err := b.DeleteKey(&kvstorev1.DeleteKeyRequest{
		AccountId:  testAccountId,
		BucketName: testBucketName,
		Key:        expectedKvp.Key,
	})
	t.Require().NoError(err, "there must not be an error when calling delete key")
	t.Require().True(resp.Ok, "the key must have been deleted")

	err = b.db.View(func(tx *bbolt.Tx) error {
		accountBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(accountBuf, testAccountId)

		acctBucket := tx.Bucket(accountBuf)
		t.Require().NotNil(acctBucket, "the account bucket must not be nil")

		bucket := acctBucket.Bucket([]byte(testBucketName))
		t.Require().NotNil(bucket, "the target bucket must not be nil")

		payload := bucket.Get([]byte(expectedKvp.Key))
		t.Require().Nil(payload, "the expected kvp must be empty")

		return nil
	})
	t.Require().NoError(err, "there must not be an error peeking into the database")
}

func FuzzBboltStore_CreateAccountBucket(f *testing.F) {

	if testing.Short() {
		f.Skipf("skipping account bucket fuzzing")
	}

	shardId, replicaId := uint64(10), uint64(20)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)
	logger := zerolog.New(zerolog.NewConsoleWriter())

	path := filepath.Join(f.TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, logger)
	require.NoError(f, err, "there must be an error when passing a bad directory")
	require.NotNil(f, b, "the bbolt store must be nil")

	for i := 0; i < fuzzRounds; i++ {
		f.Add(rand.Uint64())
	}

	f.Fuzz(func(t *testing.T, accountId uint64) {
		testOwner := "test@test.com"
		req := &kvstorev1.CreateAccountRequest{
			AccountId: accountId,
			Owner:     testOwner,
		}

		resp, err := b.CreateAccountBucket(req)
		if !errors.Is(err, bbolt.ErrBucketExists) {
			require.NoError(t, err, "there must not be an error when creating the account key")
		}
		require.NotEmpty(t, resp.GetAccountDescriptor(), "the account descriptor must not be empty")

		desc := resp.GetAccountDescriptor()
		require.Equal(t, accountId, desc.GetAccountId(), "the account ids must equal")
		require.Equal(t, testOwner, desc.GetOwner(), "the owner ids must equal")
		require.Equal(t, uint64(0), desc.GetBucketCount(), "the bucket count must be zero")
		require.NotNil(t, desc.GetCreated(), "the created timestamp must not be nil")
		require.NotNil(t, desc.GetLastUpdated(), "the last updated timestamp must not be nil")
	})
}

func FuzzBboltStore_CreateBucket(f *testing.F) {

	if testing.Short() {
		f.Skipf("skipping bucket fuzzing")
	}

	shardId, replicaId, accountId := uint64(10), uint64(20), uint64(30)
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)
	logger := zerolog.New(zerolog.NewConsoleWriter())

	path := filepath.Join(f.TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, logger)
	require.NoError(f, err, "there must be an error when passing a bad directory")
	require.NotNil(f, b, "the bbolt store must be nil")

	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: accountId,
		Owner:     "test@test.com",
	})
	require.NoError(f, err, "there must not be an error when creating the account key")

	for i := 0; i < fuzzRounds; i++ {
		f.Add(utils.RandomString(utils.RandomInt(0, 63)))
	}

	f.Fuzz(func(t *testing.T, bucketName string) {
		testOwner := "test@test.com"
		req := &kvstorev1.CreateBucketRequest{
			AccountId: accountId,
			Owner:     testOwner,
			Name:      bucketName,
		}

		resp, err := b.CreateBucket(req)

		// unrecoverable errors, very skippable
		if errors.Is(err, bbolt.ErrBucketExists) ||
			errors.Is(err, ErrEmptyBucketName) {
			return
		}

		require.NotEmpty(t, resp.GetBucketDescriptor(), "the account descriptor must not be empty")

		desc := resp.GetBucketDescriptor()
		require.Equal(t, testOwner, desc.GetOwner(), "the owner ids must equal")
		require.NotNil(t, desc.GetCreated(), "the created timestamp must not be nil")
		require.NotNil(t, desc.GetLastUpdated(), "the last updated timestamp must not be nil")
	})
}

func FuzzBboltStore_KeyStoreOperations(f *testing.F) {

	if testing.Short() {
		f.Skipf("skipping key operation fuzzing")
	}

	shardId, replicaId, accountId, testBucketName := uint64(10), uint64(20), uint64(30), "test-bucket"
	dbPath := fmt.Sprintf("shard-%d-replica-%d.db", shardId, replicaId)
	logger := zerolog.New(zerolog.NewConsoleWriter())

	path := filepath.Join(f.TempDir(), dbPath)
	b, err := newBboltStore(shardId, replicaId, path, logger)
	require.NoError(f, err, "there must be an error when passing a bad directory")
	require.NotNil(f, b, "the bbolt store must be nil")

	_, err = b.CreateAccountBucket(&kvstorev1.CreateAccountRequest{
		AccountId: accountId,
		Owner:     "test@test.com",
	})
	require.NoError(f, err, "there must not be an error when creating the account key")

	_, err = b.CreateBucket(&kvstorev1.CreateBucketRequest{
		AccountId: accountId,
		Owner:     "test@test.com",
		Name: testBucketName,
	})
	require.NoError(f, err, "there must not be an error when creating the account key")

	for i := 0; i < fuzzRounds; i++ {
		f.Add(utils.RandomString(utils.RandomInt(0, 63)))
	}

	f.Fuzz(func(t *testing.T, keyName string) {
		bytes, err := utils.RandomBytes(utils.RandomInt(0, 2<<4))
		require.NoError(t, err, "there must not be an error when generating random bytes")

		putReq := &kvstorev1.PutKeyRequest{
			AccountId: accountId,
			BucketName: testBucketName,
			KeyValuePair: &kvstorev1.KeyValue{
				Key:            keyName,
				Value:          bytes,
			},
		}

		_, err = b.PutKey(putReq)

		if !errors.Is(err, errors.New("empty key identifier")) {
			require.NoError(t, err, "there must not be an error putting a key")
		}

		getReq := &kvstorev1.GetKeyRequest{
			AccountId: accountId,
			BucketName: testBucketName,
			Key: keyName,
		}

		resp, err := b.GetKey(getReq)
		if errors.Is(err, errors.New("empty key identifier")) {
			return // can't compare empty keys
		}
		require.NoError(t, err, "there must not be an error getting a key")
		require.Equal(t, keyName, resp.GetKeyValuePair().GetKey(), "the key must be found and equal")

		// skip empty payload
		if len(bytes) == 0 {
			return
		}

		kvp := resp.GetKeyValuePair()
		require.Equal(t, bytes, kvp.GetValue(), "the value must be equal")

		delResp, err := b.DeleteKey(&kvstorev1.DeleteKeyRequest{
			AccountId: accountId,
			BucketName: testBucketName,
			Key: keyName,
		})
		if errors.Is(err, errors.New("empty key identifier")) {
			return // can't compare empty keys
		}
		require.NoError(t, err, "there must not be an error when deleting the key")
		require.True(t, delResp.GetOk(), "the key must have been deleted")
	})
}