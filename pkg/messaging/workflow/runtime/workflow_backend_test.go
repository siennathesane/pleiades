/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

import (
	"testing"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func TestWorkflowBackend(t *testing.T) {
	suite.Run(t, new(WorkflowBackendTestSuite))
}

type WorkflowBackendTestSuite struct {
	suite.Suite
	logger zerolog.Logger
	store  *systemstore.SystemStore
}

func (t *WorkflowBackendTestSuite) SetupSuite() {
	configuration.Get().SetDefault("server.datastore.basePath", t.T().TempDir())

	t.logger = utils.NewTestLogger(t.T())
	store, err := systemstore.NewSystemStore(t.logger)
	t.Require().NoError(err, "there must not be an error opening the system store")
	t.store = store
}

func (t *WorkflowBackendTestSuite) TestWorkflowState() {
	ws, err := NewWorkflowStateStore(t.store, t.logger)
	t.Require().NoError(err)
	t.Require().NotNil(ws)

	t.Equal("", ws.bucketName)
	ws.Configure("test", "test")
	t.Require().Equal("workflow-test-test", ws.bucketName)

	err = ws.Set("test", "another-test")
	t.Require().NoError(err)

	val, err := ws.Get("test")
	t.Require().NoError(err)
	t.Require().Equal("another-test", val)

	err = ws.Update("test", "", "yet-another-test")
	t.Require().NoError(err)

	err = ws.Cleanup()
	t.Require().NoError(err)
}

func (t *WorkflowBackendTestSuite) TestWorkflowStore() {
	ws, err := NewWorkflowStore(t.store, t.logger)
	t.Require().NoError(err)
	t.Require().NotNil(ws)

	t.Equal("", ws.bucketName)
	ws.Configure("test", "test")
	t.Require().Equal("workflow-test-test", ws.bucketName)

	err = ws.Set("test", []byte("another-test"))
	t.Require().NoError(err)

	val, err := ws.Get("test")
	t.Require().NoError(err)
	t.Require().Equal([]byte("another-test"), val)

	err = ws.Del("test")
	t.Require().NoError(err)

	err = ws.Cleanup()
	t.Require().NoError(err)
}
