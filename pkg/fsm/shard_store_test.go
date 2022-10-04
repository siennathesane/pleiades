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
	"fmt"
	"testing"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
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

	count := uint64(10)
	for i := uint64(0); i < count; i++ {
		testPayload := &raftv1.ShardState{
			ShardId:   i,
			Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
			Replicas: map[uint64]string{
				i+4: fmt.Sprintf("test.local.%d", i),
			},
		}
		err = store.Put(testPayload)
		t.Require().NoError(err, "there must not be an error when putting a configuration")
	}

	testPayload := &raftv1.ShardState{
		ShardId:   1,
		Type:      raftv1.StateMachineType_STATE_MACHINE_TYPE_TEST,
	}

	resp, err := store.Get(testPayload.GetShardId())
	t.Require().NoError(err, "there must not be an error when fetching a configuration")
	t.Require().NotNil(resp, "the response must not be nil")
	t.Require().Equal(testPayload.GetShardId(), resp.GetShardId(), "the test values must match up")

	resps, err := store.GetAll()
	t.Require().NoError(err, "there must not be an error when fetching a configuration")
	t.Require().NotEmpty(resps, "the response must not be nil")
	t.Require().Len(resps, int(count), "there must be ten items")
	t.Require().Equal(testPayload.GetShardId(), resps[1].GetShardId(), "the test values must match up")

	err = store.Delete(testPayload.GetShardId())
	t.Require().NoError(err, "there must not be an error deleting a configuration")
}