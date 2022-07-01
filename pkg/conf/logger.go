
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package conf

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
	dlog "github.com/lni/dragonboat/v3/logger"
	zlog "github.com/rs/zerolog"
	"go.uber.org/fx"
)

func ProvideLogger() fx.Option {
	return fx.Provide(NewConsulClient, NewEnvironmentConfig, NewLogger)
}

var _ dlog.ILogger = Logger{}
var _ dlog.Factory = Logger{}.LoggerFactory

type Logger struct {
	lifecycle fx.Lifecycle
	client    *api.Client
	env       *EnvironmentConfig
	logger    zlog.Logger
	gcpLogger gcpLogger
}

func NewLogger(lifecycle fx.Lifecycle, client *api.Client, env *EnvironmentConfig) (Logger, error) {

	l := Logger{env: env, lifecycle: lifecycle, client: client}

	// set up gcp and console logging
	var err error
	l.gcpLogger, err = newGcpLogger(l)
	if err != nil {
		return Logger{}, err
	}
	consoleWriter := zlog.ConsoleWriter{Out: os.Stdout}
	multiWriter := zlog.MultiLevelWriter(consoleWriter, l.gcpLogger)

	// adds `Lshortfile` equivalence
	zlog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	switch env.Environment {
	case Development:
		l.logger = zlog.New(multiWriter).With().
			Str("owner", "root").
			Str("env", "development").
			Timestamp().
			Caller().
			Logger().
			Level(zlog.DebugLevel)

	case Production:
		l.logger = zlog.New(multiWriter).With().
			Str("owner", "root").
			Str("env", "production").
			Timestamp().
			Caller().
			Logger().
			Level(zlog.InfoLevel)
	}

	lifecycle.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		if err := l.gcpLogger.Close(); err != nil {
			return err
		}
		return nil
	}})

	return l, nil
}

func (l Logger) LoggerFactory(pkgName string) dlog.ILogger {
	logz := l.logger.With().Str("owner", pkgName).Logger()
	return Logger{
		lifecycle: l.lifecycle,
		client:    l.client,
		env:       l.env,
		gcpLogger: l.gcpLogger,
		logger:    logz,
	}
}

var internalSeverityMap = map[dlog.LogLevel]zlog.Level{
	dlog.CRITICAL: zlog.FatalLevel,
	dlog.ERROR:    zlog.ErrorLevel,
	dlog.WARNING:  zlog.WarnLevel,
	dlog.INFO:     zlog.InfoLevel,
	dlog.DEBUG:    zlog.DebugLevel,
}

var reverseInternalSeverityMap = map[zlog.Level]dlog.LogLevel{
	zlog.FatalLevel: dlog.CRITICAL,
	zlog.ErrorLevel: dlog.ERROR,
	zlog.WarnLevel:  dlog.WARNING,
	zlog.InfoLevel:  dlog.INFO,
	zlog.DebugLevel: dlog.DEBUG,
}

func (l Logger) GetLevel() dlog.LogLevel {
	return reverseInternalSeverityMap[l.logger.GetLevel()]
}

func (l Logger) SetLevel(logLevel dlog.LogLevel) {
	l.logger = l.logger.Level(internalSeverityMap[logLevel])
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l Logger) Warningf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Stack().Msgf(format, args...)
}

func (l Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}
