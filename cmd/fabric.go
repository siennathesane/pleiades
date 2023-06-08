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
	_ cli.Command = (*FabricCommand)(nil)
)

type FabricCommand struct {
	*BaseCommand
}

func (f *FabricCommand) Help() string {
	helpText := `Commands to manage the Pleiades data fabric.`

	return helpText
}

func (f *FabricCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (f *FabricCommand) Synopsis() string {
	return "Operations on the data fabric."
}
