/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package blaze

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// StreamReceiverTest tests the StreamReceiver class.
func TestStreamReceiver(t *testing.T) {
	suite.Run(t, new(StreamReceiverTest))
}

type StreamReceiverTest struct {
	suite.Suite
	logger zerolog.Logger
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_Read() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_Write() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_Close() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose_WithError() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose_WithEOF() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose_WithEOF_WithError() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose_WithError_WithEOF() {
	//TODO implement me
	panic("implement me")
}

// StreamReceiverTest tests the StreamReceiver class.
func (s *StreamReceiverTest) TestStreamReceiver_ReadWriteClose_WithError_WithEOF_WithError() {
	//TODO implement me
	panic("implement me")
}