/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package fsm
//
//import (
//	"bytes"
//	"fmt"
//	"reflect"
//	"testing"
//
//	hostv1 "github.com/mxplusb/pleiades/pkg/protocols/v1/host"
//	"github.com/mxplusb/pleiades/pkg/services"
//	"github.com/mxplusb/pleiades/pkg/utils"
//	"capnproto.org/go/capnp/v3"
//	"github.com/rs/zerolog"
//	"github.com/stretchr/testify/suite"
//)
//
//func TestConfigServiceStoreManager(t *testing.T) {
//	suite.Run(t, new(ConfigServiceStoreManagerTests))
//}
//
//type ConfigServiceStoreManagerTests struct {
//	suite.Suite
//	logger zerolog.Logger
//	store  *services.StoreManager
//}
//
//func (test *ConfigServiceStoreManagerTests) BeforeTest(suiteName, testName string) {
//	test.logger = utils.NewTestLogger(test.T())
//	test.store = services.NewStoreManager(test.T().TempDir(), test.logger)
//}
//
//func (test *ConfigServiceStoreManagerTests) Test_Get_Returns_Error_If_Key_Does_Not_Exist() {
//	logger := test.logger.With().Str("test", "Test_Get_Returns_Error_If_Key_Does_Not_Exist").Logger()
//	cssm, err := NewConfigServiceStoreManager(logger, test.store)
//	_, err = cssm.Get("test")
//	test.Assert().Error(err, "there must be an error getting the RaftConfiguration")
//}
//
//// write a test to get a configuration from the store and compare it to the expected configuration
//func (test *ConfigServiceStoreManagerTests) Test_Get_Returns_Correct_Configuration() {
//	logger := test.logger.With().Str("test", "Test_Get_Returns_Correct_Configuration").Logger()
//	cssm, err := NewConfigServiceStoreManager(logger, test.store)
//	test.Require().NoError(err, "there must not be an error creating the ConfigServiceStoreManager")
//
//	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
//	test.Require().NoError(err, "there must be no error creating a new message")
//
//	config, err := hostv1.NewRootRaftConfiguration(seg)
//	test.Require().NoError(err, "there must not be an error creating a new RaftConfiguration")
//
//	err = config.SetId("test")
//	test.Require().NoError(err, "there must not be an error setting the id")
//
//	var buf bytes.Buffer
//	err = capnp.NewEncoder(&buf).Encode(msg)
//	test.Require().NoError(err, "there must not be an error encoding the RaftConfiguration")
//
//	err = test.store.Put("test", buf.Bytes(), reflect.TypeOf(hostv1.RaftConfiguration{}))
//	test.Require().NoError(err, "there must not be an error putting the RaftConfiguration")
//
//	res, err := cssm.Get("test")
//	test.Assert().NoError(err, "there must not be an error getting the RaftConfiguration")
//	test.Require().NotNil(res, "there must be a RaftConfiguration")
//	id, err := res.Id()
//	test.Require().NoError(err, "there must not be an error getting the id")
//	test.Require().Equal(id, "test", "the RaftConfiguration must have the correct id")
//}
//
//// write a test to get all configurations from the store and compare them to the expected configurations
//func (test *ConfigServiceStoreManagerTests) Test_GetAll_Returns_Correct_Configurations() {
//	logger := test.logger.With().Str("test", "Test_GetAll_Returns_Correct_Configurations").Logger()
//	cssm, err := NewConfigServiceStoreManager(logger, test.store)
//	test.Require().NoError(err, "there must not be an error creating the ConfigServiceStoreManager")
//
//	for i := 0; i < 10; i++ {
//		testName := fmt.Sprintf("test_%d", i)
//		msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
//		test.Require().NoError(err, "there must be no error creating a new message")
//
//		config, err := hostv1.NewRootRaftConfiguration(seg)
//		test.Require().NoError(err, "there must not be an error creating a new RaftConfiguration")
//
//		err = config.SetId(testName)
//		test.Require().NoError(err, "there must not be an error setting the id")
//
//		var buf bytes.Buffer
//		err = capnp.NewEncoder(&buf).Encode(msg)
//		test.Require().NoError(err, "there must not be an error encoding the RaftConfiguration")
//
//		err = test.store.Put(testName, buf.Bytes(), reflect.TypeOf(hostv1.RaftConfiguration{}))
//		test.Require().NoError(err, "there must not be an error putting the RaftConfiguration")
//	}
//
//	configs, err := cssm.GetAll()
//	test.Require().NoError(err, "there must not be an error getting the RaftConfigurations")
//	test.Assert().Len(configs, 10, "there must be 10 RaftConfigurations")
//
//	for i := 0; i < 10; i++ {
//		testName := fmt.Sprintf("test_%d", i)
//		test.Assert().Contains(configs, testName, "there must be a RaftConfiguration with the correct id")
//	}
//}
//
//// write a test to put a configuration into the store and compare it to the expected configuration
//func (test *ConfigServiceStoreManagerTests) Test_Put_Returns_Correct_Configuration() {
//	logger := test.logger.With().Str("test", "Test_Put_Returns_Correct_Configuration").Logger()
//	cssm, err := NewConfigServiceStoreManager(logger, test.store)
//	test.Require().NoError(err, "there must not be an error creating the ConfigServiceStoreManager")
//
//	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
//	test.Require().NoError(err, "there must be no error creating a new message")
//
//	config, err := hostv1.NewRootRaftConfiguration(seg)
//	test.Require().NoError(err, "there must not be an error creating a new RaftConfiguration")
//
//	err = config.SetId("test")
//	test.Require().NoError(err, "there must not be an error setting the id")
//	config.SetCheckQuorum(true)
//
//	err = cssm.Put("test", &config)
//	test.Require().NoError(err, "there must not be an error putting the RaftConfiguration")
//}
//
//func (test *ConfigServiceStoreManagerTests) Test_Put_Returns_Error_If_Payload_Is_Nil() {
//	logger := test.logger.With().Str("test", "Test_Put_Returns_Error_If_Payload_Is_Nil").Logger()
//	cssm, err := NewConfigServiceStoreManager(logger, test.store)
//	test.Require().NoError(err, "there must not be an error creating the ConfigServiceStoreManager")
//
//	err = cssm.Put("test", nil)
//	test.Assert().Error(err, "there must be an error putting the RaftConfiguration")
//}
