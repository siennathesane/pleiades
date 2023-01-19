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
	EnvPleiadesInsecureSkipVerify = "PLEIADES_INSECURE_SKIP_VERIFY"
	EnvPleiadesDebug              = "PLEIADES_DEBUG"
	EnvPleiadesTrace              = "PLEIADES_TRACE"
	EnvPleiadesDefaultOutput      = "PLEIADES_OUTPUT"

	/* TLS Configs */
	EnvPleiadesCaCert   = "PLEIADES_CA_CERT_FILE"
	EnvPleiadesCertFile = "PLEIADES_CERT_FILE"
	EnvPleiadesKeyFile  = "PLEIADES_KEY_FILE"

	/* Server Variables */
	EnvPleiadesDeploymentId      = "PLEIADES_DEPLOYMENT_ID"
	EnvPleiadesDataDir           = "PLEIADES_DATA_DIR"
	EnvPleiadesFabricAddr        = "PLEIADES_FABRIC_ADDR"
	EnvPleidesListenAddr         = "PLEIADES_LISTEN_ADDR"
	EnvPleiadesHttpPort          = "PLEIADES_HTTP_PORT"
	EnvPleiadesFabricPort        = "PLEIADES_FABRIC_PORT"
	EnvPleiadesConstellationPort = "PLEIADES_CONSTELLATION_PORT"
	EnvPleiadesNotifyCommit      = "PLEIADES_NOTIFY_COMMIT"
	EnvPleiadesRoundTrip         = "PLEIADES_ROUND_TRIP_MS"

	flagNameHost = "address"

	exitCodeGood               = 0
	exitCodeGenericBad         = 1
	exitCodeFailureToParseArgs = 2
	exitCodeRemote             = 3
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
		"kv put": func() (cli.Command, error) {
			return &KvPutCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"kv delete": func() (cli.Command, error) {
			return &KvDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account": func() (cli.Command, error) {
			return &AccountCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account create": func() (cli.Command, error) {
			return &AccountCreateCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account delete": func() (cli.Command, error) {
			return &AccountDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket": func() (cli.Command, error) {
			return &BucketCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket create": func() (cli.Command, error) {
			return &BucketCreateCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket delete": func() (cli.Command, error) {
			return &BucketDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
	}
}
