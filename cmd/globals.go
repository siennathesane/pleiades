/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"github.com/mitchellh/cli"
)

const (
	/* A group! */
	EnvPleiadesUrl                = "PLEIADES_ADDR"
	EnvPleiadesInsecureSkipVerify = "PLEIADES_TLS_SKIP_VERIFY"
	EnvPleiadesCaCert             = "PLEIADES_CA_CERT"
	EnvPleiadesDebug              = "PLEIADES_DEBUG"
	EnvPleiadesTrace              = "PLEIADES_TRACE"
	EnvPleiadesDefaultOutput      = "PLEIADES_OUTPUT"

	flagNameHost = "address"

	exitCodeGood               = 0
	exitCodeGenericBad         = 1
	exitCodeFailureToParseArgs = 2
	exitCodeRemote = 3
)

var Commands map[string]cli.CommandFactory

func initCommands(ui cli.Ui) {
	getBaseCmd := func() *BaseCommand {
		return &BaseCommand{
			UI: ui,
		}
	}

	Commands = map[string]cli.CommandFactory{
		"kv": func() (cli.Command, error) {
			return &KvCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"kv get": func() (cli.Command, error) {
			return &KvGetCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
	}
}
