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
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lni/dragonboat/v3"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestBBoltStoreManager(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping bbolt store manager tests")
	}
	suite.Run(t, new(bboltStoreManagerTestSuite))
}

type bboltStoreManagerTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	tm             *raftTransactionManager
	sm             *raftShardManager
	nh             *dragonboat.NodeHost
	defaultTimeout time.Duration
}

func (t *bboltStoreManagerTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
	t.defaultTimeout = 300 * time.Millisecond

	t.nh = buildTestNodeHost(t.T())
	t.tm = newTransactionManager(t.nh, t.logger)
	t.sm = newShardManager(t.nh, t.logger)

	// ensure that bbolt uses the temp directory
	viper.SetDefault("datastore.basePath", t.T().TempDir())

	// shardLimit+1
	var wg sync.WaitGroup
	for i := uint64(1); i < 257; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			err := t.sm.NewShard(i, rand.Uint64(), BBoltStateMachineType, utils.Timeout(t.defaultTimeout))
			t.Require().NoError(err, "there must not be an error when starting the bbolt state machine")
			utils.Wait(t.defaultTimeout)
		}()
		utils.Wait(100 * time.Millisecond)
	}
	wg.Wait()
}

func (t *bboltStoreManagerTestSuite) TestCreateAccount() {
	storeManager := newBboltStoreManager(t.tm, t.nh, t.logger)

	testBaseAccountId := rand.Uint64()
	testOwner := "test@test.com"

	// no transaction
	resp, err := storeManager.CreateAccount(&database.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().NotEmpty(resp.GetAccountDescriptor(), "the account descriptor must not be empty")

	// create 20 new accounts
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		resp, err := storeManager.CreateAccount(&database.CreateAccountRequest{
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

func (t *bboltStoreManagerTestSuite) TestDeleteAccount() {
	storeManager := newBboltStoreManager(t.tm, t.nh, t.logger)

	testBaseAccountId := rand.Uint64()
	testOwner := "test@test.com"

	// no transaction
	_, err := storeManager.CreateAccount(&database.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")

	// no transaction
	resp, err := storeManager.DeleteAccount(&database.DeleteAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().True(resp.Ok, "the request must be okay")

	// create and delete 20 new accounts
	for i := testBaseAccountId + 2; i < testBaseAccountId+2+20; i++ {
		createAccountReply, err := storeManager.CreateAccount(&database.CreateAccountRequest{
			AccountId:   i,
			Owner:       testOwner,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(createAccountReply, "the response must not be nil")
		t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")
		t.Require().NotEmpty(i, createAccountReply.GetAccountDescriptor().GetAccountId(), "the account descriptor must not be empty")

		deleteAccountReply, err := storeManager.DeleteAccount(&database.DeleteAccountRequest{
			AccountId:   i,
			Owner:       testOwner,
			Transaction: nil,
		})
		t.Require().NoError(err, "there must not be an error when creating an account")
		t.Require().NotNil(deleteAccountReply, "the response must not be nil")
		t.Require().True(deleteAccountReply.Ok, "the request must be okay")
	}
}

func (t *bboltStoreManagerTestSuite) TestCreateBucket() {
	storeManager := newBboltStoreManager(t.tm, t.nh, t.logger)

	testBaseAccountId := rand.Uint64()
	testBucketName := utils.RandomString(10)
	testOwner := "test@test.com"

	// no transaction
	createAccountReply, err := storeManager.CreateAccount(&database.CreateAccountRequest{
		AccountId:   testBaseAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(createAccountReply, "the response must not be nil")
	t.Require().NotEmpty(createAccountReply.GetAccountDescriptor(), "the account descriptor must not be empty")

	// no transaction
	createBucketReply, err := storeManager.CreateBucket(&database.CreateBucketRequest{
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
		resp, err := storeManager.CreateBucket(&database.CreateBucketRequest{
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