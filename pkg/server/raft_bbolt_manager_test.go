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
	for i := uint64(1); i < 257; i++ {
		err := t.sm.NewShard(i, rand.Uint64(), BBoltStateMachineType, utils.Timeout(t.defaultTimeout))
		t.Require().NoError(err, "there must not be an error when starting the bbolt state machine")
		utils.Wait(300 * time.Millisecond)
	}
}

func (t *bboltStoreManagerTestSuite) TestCreateAccount() {
	storeManager := newBboltStoreManager(t.tm, t.nh, t.logger)

	testAccountId := rand.Uint64()
	//testBucketId := utils.RandomString(10)
	testOwner := "test@test.com"

	// no transaction
	resp, err := storeManager.CreateAccount(&database.CreateAccountRequest{
		AccountId:   testAccountId,
		Owner:       testOwner,
		Transaction: nil,
	})
	t.Require().NoError(err, "there must not be an error when creating an account")
	t.Require().NotNil(resp, "the response must not be nil")
}
