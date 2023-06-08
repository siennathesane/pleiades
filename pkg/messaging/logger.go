/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package messaging

import (
	"github.com/nats-io/nats-server/v2/server"
	zlog "github.com/rs/zerolog"
)

var _ server.Logger = (*messagingLogger)(nil)

type messagingLogger struct {
	logger zlog.Logger
}

func (n *messagingLogger) Noticef(format string, v ...interface{}) {
	n.logger.Info().Msgf(format, v...)
}

func (n *messagingLogger) Warnf(format string, v ...interface{}) {
	n.logger.Warn().Msgf(format, v...)
}

func (n *messagingLogger) Fatalf(format string, v ...interface{}) {
	n.logger.Fatal().Msgf(format, v...)
}

func (n *messagingLogger) Errorf(format string, v ...interface{}) {
	n.logger.Error().Msgf(format, v...)
}

func (n *messagingLogger) Debugf(format string, v ...interface{}) {
	n.logger.Debug().Msgf(format, v...)
}

func (n *messagingLogger) Tracef(format string, v ...interface{}) {
	n.logger.Trace().Msgf(format, v...)
}
