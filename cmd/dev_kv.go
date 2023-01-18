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

// kvCmd represents the kv command
var kvCmd = &cobra.Command{
	Use:   "kv",
	Short: "operations on keys!",
}

var (
	payload []byte
	key string
)

func init() {
	rootCmd	.AddCommand(kvCmd)
	kvCmd.PersistentFlags().Uint64Var(&accountId, "account-id", 0, "the account to operate on")
	kvCmd.PersistentFlags().StringVarP(&bucketName, "bucket", "b", "","bucket to place the key in")
}
