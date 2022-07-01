/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package services

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/utils"
)

type TestStruct struct {
	Key string
	Val TestNestedStruct
}

type TestNestedStruct struct {
	Bool bool
}

type TManagerTestSuite struct {
	suite.Suite
	logger zerolog.Logger
}

func TestStoreManager(t *testing.T) {
	suite.Run(t, new(TManagerTestSuite))
}

func (s *TManagerTestSuite) SetupSuite() {
	s.logger = utils.NewTestLogger(s.T())
}

func (s *TManagerTestSuite) TestNewGenericManager() {
	manager := NewStoreManager(s.T().TempDir(), s.logger)

	var err error
	require.NotPanics(s.T(), func() {
		err = manager.Start(false)
	})
	require.Nil(s.T(), err, "there should be no error starting the store manager")

	require.NotPanics(s.T(), func() {
		err = manager.Stop(false)
	}, "the manager instantiation shouldn't panic")
	require.Nil(s.T(), err, "there shouldn't be an error stopping the manager")
}

func (s *TManagerTestSuite) BeforeTest(suiteName, testName string) {
	fullPath, err := dbPath(s.T().TempDir())
	if err != nil {
		s.T().Error(err)
	}

	if err := os.RemoveAll(fullPath); err != nil {
		s.T().Error(err)
	}
}

func (s *TManagerTestSuite) TestManagerInitOnPut() {
	t := utils.NewTestLogger(s.T())

	manager := NewStoreManager(s.T().TempDir(), t)
	require.Nil(s.T(), manager.Start(false), "there should be no errors starting the store manager")

	testStruct := &TestStruct{
		Key: "test-key",
		Val: TestNestedStruct{Bool: true},
	}

	val, err := json.Marshal(testStruct)
	require.Nil(s.T(), err, "there should be no serialization error")

	require.NotPanics(s.T(), func() {
		err = manager.Put("test-key", val, reflect.TypeOf(testStruct))
	}, "putting a type shouldn't panic.")
	require.Nil(s.T(), err, "putting a type shouldn't have an error")
}

func (s *TManagerTestSuite) TestManagerGet() {
	t := utils.NewTestLogger(s.T())

	manager := NewStoreManager(s.T().TempDir(), t)
	require.Nil(s.T(), manager.Start(false), "there should be no errors starting the store manager")

	testStruct := &TestStruct{
		Key: "test-key",
		Val: TestNestedStruct{Bool: true},
	}

	val, _ := json.Marshal(testStruct)

	var err error
	require.NotPanics(s.T(), func() {
		err = manager.Put("test-key", val, reflect.TypeOf(testStruct))
	}, "putting a type shouldn't panic.")
	require.Nil(s.T(), err, "putting a type shouldn't have an error")

	val, err = manager.Get("test-key", reflect.TypeOf(testStruct))
	require.Nil(s.T(), err, "the test struct shouldn't throw an error on fetch")

	var res *TestStruct
	err = json.Unmarshal(val, &res)
	require.Nil(s.T(), err, "there shouldn't be a deserialize error")
	require.NotNil(s.T(), res, "the test struct result shouldn't be nil")
	require.Equal(s.T(), testStruct, res, "the test struct and the returned value should be the same")
}
