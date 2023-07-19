/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cli

import (
	"strings"

	"github.com/mitchellh/cli"
)

var (
	payload []byte
	key     string

	_ cli.Command = (*KvCommand)(nil)
)

type KvCommand struct {
	*BaseCommand
}

func (k *KvCommand) Help() string {
	helpText := `
Commands to operate on key value pairs.
`
	return strings.TrimSpace(helpText)
}

func (k *KvCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (k *KvCommand) Synopsis() string {
	return "Operations on key value pairs!"
}
