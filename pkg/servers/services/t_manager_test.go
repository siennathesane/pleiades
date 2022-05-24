package services

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx/fxtest"
	"r3t.io/pleiades/pkg/conf"
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
	lifecycle *fxtest.Lifecycle
	client    *api.Client
	env       *conf.EnvironmentConfig
}

func TestStoreManager(t *testing.T) {
	suite.Run(t, new(TManagerTestSuite))
}

func (s *TManagerTestSuite) SetupSuite() {
	var err error
	s.lifecycle = fxtest.NewLifecycle(s.T())
	s.client, err = conf.NewConsulClient(s.lifecycle)
	require.Nil(s.T(), err, "failed to connect to consul")
	require.NotNil(s.T(), s.client, "the consul client can't be empty")

	s.env, err = conf.NewEnvironmentConfig(s.client)
	require.Nil(s.T(), err, "the environment config is needed")
	require.NotNil(s.T(), s.env, "the environment config must be rendered")
}

func (s *TManagerTestSuite) TestNewGenericManager() {
	t := &conf.MockLogger{}

	manager := NewStoreManager(s.env, t, s.client)

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
	fullPath, err := dbPath(s.env.BaseDir)
	if err != nil {
		s.T().Error(err)
	}

	if err := os.RemoveAll(fullPath); err != nil {
		s.T().Error(err)
	}
}

func (s *TManagerTestSuite) TestManagerInitOnPut() {
	t := &conf.MockLogger{}

	manager := NewStoreManager(s.env, t, s.client)
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
	t := &conf.MockLogger{}

	manager := NewStoreManager(s.env, t, s.client)
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
