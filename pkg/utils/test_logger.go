package utils

import (
	"testing"

	"github.com/rs/zerolog"
)

func NewTestLogger(t *testing.T) zerolog.Logger {
	testWriter := zerolog.NewTestWriter(t)
	return zerolog.New(testWriter)
}
