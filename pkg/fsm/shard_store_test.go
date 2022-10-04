/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package fsm

import (
	"testing"

	raftv1 "github.com/mxplusb/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestShardStore(t *testing.T) {
	suite.Run(t, new(shardStoreTestSuite))
}

type shardStoreTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func (t *shardStoreTestSuite) SetupSuite() {
	t.logger = utils.NewTestLogger(t.T())
}

func (t *shardStoreTestSuite) TestLifecycle() {
	configuration.Get().SetDefault("server.datastore.basePath", t.T().TempDir())

	store, err := NewShardStore(t.logger)
	t.Require().NoError(err, "there must not be an error when opening the shard store")
	t.Require().NotNil(store, "the shard store must not be empty")

	testPayload := &raftv1.NewShardRequest{
		ShardId:   1,
		ReplicaId: 1,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
		Hostname:  "test.local",
		Timeout:   100,
	}

	err = store.Put(testPayload)
	t.Require().NoError(err, "there must not be an error when putting a configuration")

	resp, err := store.Get(testPayload.GetShardId())
	t.Require().NoError(err, "there must not be an error when fetching a configuration")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().Equal(testPayload.GetTimeout(), resp.GetTimeout(), "the test values must match up")

	err = store.Delete(testPayload.GetShardId())
	t.Require().NoError(err, "there must not be an error deleting a configuration")
}