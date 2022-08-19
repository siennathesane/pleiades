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
	"runtime"
	"strconv"

	"github.com/mxplusb/pleiades/pkg"
	dlog "github.com/lni/dragonboat/v3/logger"
	zlog "github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
)

var (
	rootLogger  zlog.Logger
	writers     []io.Writer
	multiWriter = zlog.MultiLevelWriter(writers...)
)

func init() {
	if runtime.GOOS != "darwin" {
		writers = append(writers, journald.NewJournalDWriter())
	} else {
		writers = append(writers, zlog.ConsoleWriter{Out: os.Stdout})
	}

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

	rootLogger = zlog.New(multiWriter).With().
		Str("sha", pkg.Sha).
		Timestamp().
		Caller().
		Logger().
		Level(zlog.InfoLevel)
}

func NewRootLogger() zlog.Logger {
	return rootLogger
}

func DragonboatLoggerFactory(pkgName string) dlog.ILogger {
	logz := rootLogger.With().Str("pkg", pkgName).Logger()
	return DragonboatLoggerAdapter{
		logger: logz,
	}
}

var _ dlog.ILogger = DragonboatLoggerAdapter{}

type DragonboatLoggerAdapter struct {
	logger zlog.Logger
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

func (l DragonboatLoggerAdapter) SetLevel(logLevel dlog.LogLevel) {
	l.logger = l.logger.Level(internalSeverityMap[logLevel])
}

func (l DragonboatLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l DragonboatLoggerAdapter) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l DragonboatLoggerAdapter) Warningf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l DragonboatLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Error().Stack().Msgf(format, args...)
}

func (l DragonboatLoggerAdapter) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}
