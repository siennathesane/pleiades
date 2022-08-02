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
	"io"
	"os"
	"strconv"

	"gitlab.com/anthropos-labs/pleiades/pkg"
	dlog "github.com/lni/dragonboat/v3/logger"
	zlog "github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
)

var _ dlog.ILogger = Logger{}

type Logger struct {
	logger zlog.Logger
}

func NewLogger(writers ...io.Writer) (Logger, error) {

	l := Logger{}

	writers = append(writers, zlog.ConsoleWriter{Out: os.Stdout}, journald.NewJournalDWriter())

	// write to both console and journald for linux
	multiWriter := zlog.MultiLevelWriter(writers...)

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

	l.logger = zlog.New(multiWriter).With().
		Str("sha", pkg.Sha).
		Timestamp().
		Caller().
		Logger().
		Level(zlog.InfoLevel)

	return l, nil
}

func (l Logger) LoggerFactory(pkgName string) dlog.ILogger {
	logz := l.logger.With().Str("pkg", pkgName).Logger()
	return Logger{
		logger: logz,
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

func (l Logger) GetLogger() zlog.Logger {
	return l.logger
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
