/*
 * Copyright (c) 2022-2023 Sienna Lloyd
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

var (
	_ cli.Command = (*BucketCommand)(nil)
)

type BucketCommand struct {
	*BaseCommand
}

func (a *BucketCommand) Help() string {
	helpText := `bucket help text`

	return helpText
}

func (a *BucketCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (a *BucketCommand) Synopsis() string {
	return "Operations on bucket!"
}

