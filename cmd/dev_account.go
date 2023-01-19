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
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "operations on accounts!",
}

var (
	accountId    uint64
	accountOwner string
)

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.PersistentFlags().String("host", "http://localhost:8080", "target host for a pleiades cluster")
	//config.BindPFlag("client.grpcAddr", accountCmd.PersistentFlags().Lookup("host"))

	accountCmd.PersistentFlags().Uint64Var(&accountId, "account-id", 0, "the account to operate on")
}
