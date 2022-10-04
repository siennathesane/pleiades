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
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	kvstorev1 "github.com/mxplusb/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestRaftBBoltStoreManagerGrpcAdapter(t *testing.T) {
	suite.Run(t, new(raftBboltStoreManagerGrpcAdapterTestSuite))
}

type raftBboltStoreManagerGrpcAdapterTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	tm             *raftTransactionManager
	sm             *raftShardManager
	storem         *bboltStoreManager
	conn           *grpc.ClientConn
	srv            *grpc.Server
	nh             *dragonboat.NodeHost
	rh             *raftHost
	defaultTimeout time.Duration
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 300 * time.Millisecond

	// ensure that bbolt uses the temp directory
	configuration.Get().SetDefault("server.datastore.dataDir", t.T().TempDir())

	t.nh = buildTestNodeHost(t.T())
	t.tm = newTransactionManager(t.nh, t.logger)
	t.sm = newShardManager(t.nh, t.logger)
	t.storem = newBboltStoreManager(t.tm, t.nh, t.logger)

	// shardLimit+1
	var wg sync.WaitGroup
	for i := uint64(1); i < 257; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			err := t.sm.NewShard(i, i*4, BBoltStateMachineType, utils.Timeout(t.defaultTimeout))
			t.Require().NoError(err, "there must not be an error when starting the bbolt state machine")
			utils.Wait(t.defaultTimeout)
		}()
		utils.Wait(100 * time.Millisecond)
	}
	defer wg.Wait()

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	ctx := context.Background()
	t.srv = grpc.NewServer()

	kvstorev1.RegisterKvStoreServiceServer(t.srv, &raftBBoltStoreManagerGrpcAdapter{
		logger:       t.logger,
		storeManager: t.storem,
	})

	go func() {
		if err := t.srv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	t.conn, _ = grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) TestCreateAccount() {
	client := kvstorev1.NewKvStoreServiceClient(t.conn)

	testBaseAccountId := rand.Uint64()
	testOwner := "test@test.com"

	// no transaction
	resp, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().NotEmpty(resp.GetAccountDescriptor(), "the account descriptor must not be empty")

	// create 20 new accounts
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		resp, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
			AccountId:   i,
			Owner:       testOwner,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(resp, "the response must not be nil")
		t.Require().NotEmpty(resp.GetAccountDescriptor(), "the account descriptor must not be empty")
		t.Require().NotEmpty(i, resp.GetAccountDescriptor().GetAccountId(), "the account descriptor must not be empty")
	}
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) TestDeleteAccount() {
	client := kvstorev1.NewKvStoreServiceClient(t.conn)

	testBaseAccountId := rand.Uint64()
	testOwner := "test@test.com"

	// no transaction
	_, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")

	// no transaction
	resp, err := client.DeleteAccount(context.TODO(), &kvstorev1.DeleteAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().True(resp.Ok, "the request must be okay")

	// create and delete 20 new accounts
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		createAccountReply, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
			AccountId:   i,
			Owner:       testOwner,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(createAccountReply, "the response must not be nil")
		t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")
		t.Require().NotEmpty(i, createAccountReply.GetAccountDescriptor().GetAccountId(), "the account descriptor must not be empty")

		deleteAccountReply, err := client.DeleteAccount(context.TODO(), &kvstorev1.DeleteAccountRequest{
			AccountId:   i,
			Owner:       testOwner,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(deleteAccountReply, "the response must not be nil")
		t.Require().True(deleteAccountReply.Ok, "the request must be okay")
	}
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) TestCreateBucket() {
	client := kvstorev1.NewKvStoreServiceClient(t.conn)

	testBaseAccountId := rand.Uint64()
	testBucketName := utils.RandomString(10)
	testOwner := "test@test.com"

	// no transaction
	createAccountReply, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createAccountReply, "the response must not be nil")
	t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")

	// no transaction
	createBucketReply, err := client.CreateBucket(context.TODO(), &kvstorev1.CreateBucketRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Name:        testBucketName,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createBucketReply, "the response must not be nil")
	t.Require().NotEmpty(createBucketReply.GetBucketDescriptor(), "the account descriptor must not be empty")

	//create 20 new buckets
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		testBucketName = utils.RandomString(10)
		resp, err := client.CreateBucket(context.TODO(), &kvstorev1.CreateBucketRequest{
			AccountId:   testBaseAccountId,
			Owner:       testOwner,
			Name:        testBucketName,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(resp, "the response must not be nil")
		t.Require().NotEmpty(resp.GetBucketDescriptor(), "the bucket descriptor must not be empty")
		t.Require().Empty(resp.GetBucketDescriptor().GetKeyCount(), "the key count must be zero")
	}
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) TestDeleteBucket() {
	client := kvstorev1.NewKvStoreServiceClient(t.conn)

	testBaseAccountId := rand.Uint64()
	testBucketName := utils.RandomString(10)
	testOwner := "test@test.com"

	// no transaction
	createAccountReply, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createAccountReply, "the response must not be nil")
	t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")

	// no transaction
	createBucketReply, err := client.CreateBucket(context.TODO(), &kvstorev1.CreateBucketRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Name:        testBucketName,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createBucketReply, "the response must not be nil")
	t.Require().NotEmpty(createBucketReply.GetBucketDescriptor(), "the account descriptor must not be empty")

	// no transaction
	deleteBucketReply, err := client.DeleteBucket(context.TODO(), &kvstorev1.DeleteBucketRequest{
		AccountId:   testBaseAccountId,
		Name:        testBucketName,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(deleteBucketReply, "the response must not be nil")
	t.Require().True(deleteBucketReply.GetOk(), "the account descriptor must not be empty")

	//create 20 new buckets
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		testBucketName = utils.RandomString(10)
		bucketReply, err := client.CreateBucket(context.TODO(), &kvstorev1.CreateBucketRequest{
			AccountId:   testBaseAccountId,
			Owner:       testOwner,
			Name:        testBucketName,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(bucketReply, "the response must not be nil")
		t.Require().NotEmpty(bucketReply.GetBucketDescriptor(), "the bucket descriptor must not be empty")
		t.Require().Empty(bucketReply.GetBucketDescriptor().GetKeyCount(), "the key count must be zero")

		deleteBucket, err := client.DeleteBucket(context.TODO(), &kvstorev1.DeleteBucketRequest{
			AccountId:   testBaseAccountId,
			Name:        testBucketName,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(deleteBucket, "the response must not be nil")
		t.Require().True(deleteBucket.GetOk(), "the account descriptor must not be empty")
	}
}

func (t *raftBboltStoreManagerGrpcAdapterTestSuite) TestKeyLifecycle() {
	client := kvstorev1.NewKvStoreServiceClient(t.conn)

	testBaseAccountId := rand.Uint64()
	testBucketName := utils.RandomString(10)
	testOwner := "test@test.com"

	// no transaction
	createAccountReply, err := client.CreateAccount(context.TODO(), &kvstorev1.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createAccountReply, "the response must not be nil")
	t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")

	// no transaction
	createBucketReply, err := client.CreateBucket(context.TODO(), &kvstorev1.CreateBucketRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Name:        testBucketName,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating a bucket")
	t.Require().NotNil(createBucketReply, "the response must not be nil")
	t.Require().NotEmpty(createBucketReply.GetBucketDescriptor(), "the account descriptor must not be empty")

	testPutValue, _ := utils.RandomBytes(128)
	testKvp := &kvstorev1.KeyValue{
		Key:            "test-key",
		CreateRevision: 0,
		ModRevision:    0,
		Version:        0,
		Value:          testPutValue,
		Lease:          0,
	}

	putKeyReply, err := client.PutKey(context.TODO(), &kvstorev1.PutKeyRequest{
		AccountId:    testBaseAccountId,
		BucketName:   testBucketName,
		KeyValuePair: testKvp,
		Transaction:  nil,
	})
	t.Require().NoError(err, "there must not be an error when putting a key")
	t.Require().NotNil(putKeyReply, "the key response must not be empty")

	getKeyReply, err := client.GetKey(context.TODO(), &kvstorev1.GetKeyRequest{
		AccountId:  testBaseAccountId,
		BucketName: testBucketName,
		Key:        testKvp.Key,
	})
	t.Require().NoError(err, "there must not be an error when getting a key")
	t.Require().NotNil(getKeyReply, "the reply must not be nil")
	t.Require().NotEmpty(getKeyReply.GetKeyValuePair(), "the key value pair must not be empty")
	t.Require().Equal(testKvp, getKeyReply.GetKeyValuePair(), "the kvps must match")

	deleteKeyReply, err := client.DeleteKey(context.TODO(), &kvstorev1.DeleteKeyRequest{
		AccountId:   testBaseAccountId,
		BucketName:  testBucketName,
		Key:         testKvp.Key,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when deleting a key")
	t.Require().NotNil(deleteKeyReply, "the reply must not be nil")
	t.Require().True(deleteKeyReply.GetOk(), "the key must have been deleted")

	// handle the lifecycle for 20 random keys
	for i := 0; i < 20; i++ {
		testPutValue, _ = utils.RandomBytes(128)
		testKvp = &kvstorev1.KeyValue{
			Key:            fmt.Sprintf("test-key-%d", i),
			CreateRevision: 0,
			ModRevision:    0,
			Version:        0,
			Value:          testPutValue,
			Lease:          0,
		}

		putKeyReply, err = client.PutKey(context.TODO(), &kvstorev1.PutKeyRequest{
			AccountId:    testBaseAccountId,
			BucketName:   testBucketName,
			KeyValuePair: testKvp,
			Transaction:  nil,
		})
		t.Require().NoError(err, "there must not be an error when putting a key")
		t.Require().NotNil(putKeyReply, "the key response must not be empty")

		getKeyReply, err = client.GetKey(context.TODO(), &kvstorev1.GetKeyRequest{
			AccountId:  testBaseAccountId,
			BucketName: testBucketName,
			Key:        testKvp.Key,
		})
		t.Require().NoError(err, "there must not be an error when getting a key")
		t.Require().NotNil(getKeyReply, "the reply must not be nil")
		t.Require().NotEmpty(getKeyReply.GetKeyValuePair(), "the key value pair must not be empty")
		t.Require().Equal(testKvp, getKeyReply.GetKeyValuePair(), "the kvps must match")

		deleteKeyReply, err = client.DeleteKey(context.TODO(), &kvstorev1.DeleteKeyRequest{
			AccountId:   testBaseAccountId,
			BucketName:  testBucketName,
			Key:         testKvp.Key,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when deleting a key")
		t.Require().NotNil(deleteKeyReply, "the reply must not be nil")
		t.Require().True(deleteKeyReply.GetOk(), "the key must have been deleted")
	}
}
