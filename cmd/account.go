/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"github.com/mitchellh/cli"
)

var (
	_ cli.Command = (*AccountCommand)(nil)
)

type AccountCommand struct {
	*BaseCommand
}

func (a *AccountCommand) Help() string {
	helpText := `Commands to manage accounts.`

	return helpText
}

func (a *AccountCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (a *AccountCommand) Synopsis() string {
	return "Operations on accounts!"
}

