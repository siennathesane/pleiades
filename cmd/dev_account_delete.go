/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */
package cmd

import (
	"context"
	"net/http"

	kvstorev1 "github.com/mxplusb/api/kvstore/v1"
	"github.com/mxplusb/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// accountDeleteCmd represents the accountPut command
var accountDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an account",
	Run:   deleteAccountCmd,
}

func init() {
	accountCmd.AddCommand(accountDeleteCmd)
	accountDeleteCmd.PersistentFlags().StringVar(&accountOwner, "owner", "", "the email owning the account")
}

func deleteAccountCmd(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	if accountId == 0 {
		logger.Fatal().Msg("account id cannot be zero")
	}

	client := kvstorev1connect.NewKvStoreServiceClient(http.DefaultClient, "http://localhost:8080")

	descriptor, err := client.DeleteAccount(context.Background(), connect.NewRequest(&kvstorev1.DeleteAccountRequest{AccountId: accountId, Owner: accountOwner}))
	if err != nil {
		logger.Fatal().Err(err).Uint64("account-id", accountId).Msg("can't create account")
	}

	print(protojson.Format(descriptor.Msg))
}
