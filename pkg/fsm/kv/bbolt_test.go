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

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
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

	req := &database.CreateAccountRequest{}
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
	foundAcctDescriptor := &database.AccountDescriptor{}
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
	req := &database.CreateBucketRequest{
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

	_, err = b.CreateAccountBucket(&database.CreateAccountRequest{
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

	acctDescriptor := &database.AccountDescriptor{}
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

	t.Require().Equal(uint64(1),acctDescriptor.GetBucketCount(), "the number of buckets must match")
	t.Require().Equal(testBucketName,acctDescriptor.GetBuckets()[0], "the bucket names must match")
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
	_, err = b.CreateAccountBucket(&database.CreateAccountRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
	})
	t.Require().NoError(err, "there must not be an error creating the test account bucket")

	_, err = b.CreateBucket(&database.CreateBucketRequest{
		AccountId: testAccountId,
		Owner:     testOwner,
		Name: testBucketName,
	})
	t.Require().NoError(err, "there must not be an error creating the test bucket")

	// bad request
	req := &database.DeleteBucketRequest{
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

	acctDescriptor := &database.AccountDescriptor{}
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

	t.Require().Equal(uint64(0),acctDescriptor.GetBucketCount(), "the number of buckets must match")
	t.Require().Empty(acctDescriptor.GetBuckets(), "the bucket names must empty")
}