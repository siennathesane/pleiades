/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

import (
	"os"

	"github.com/mxplusb/pleiades/pkg"
	zlog "github.com/rs/zerolog"
)

var (
	rootLogger zlog.Logger
)

func init() {
	rootLogger = zlog.New(zlog.ConsoleWriter{Out: os.Stdout}).
		With().
		Str("sha", pkg.Sha).
		Timestamp().
		Logger().
		Level(zlog.InfoLevel)
}

func NewRootLogger() zlog.Logger {
	return rootLogger
}
