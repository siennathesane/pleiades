package conf

import (
	"github.com/lni/dragonboat/v3/logger"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (t MockLogger) SetLevel(level logger.LogLevel) {
	t.Called(level)
}

func (t MockLogger) Debugf(format string, args ...interface{}) {
	t.Called(format, args)
}

func (t MockLogger) Infof(format string, args ...interface{}) {
	t.Called(format, args)
}

func (t MockLogger) Warningf(format string, args ...interface{}) {
	t.Called(format, args)
}

func (t MockLogger) Errorf(format string, args ...interface{}) {
	t.On("Errorf").Return(nil)
}

func (t MockLogger) Panicf(format string, args ...interface{}) {
	t.Called(format, args)
}
