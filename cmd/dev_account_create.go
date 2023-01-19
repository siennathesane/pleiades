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
	"context"
	"net/http"
	"net/url"

	kvstorev1 "github.com/mxplusb/api/kvstore/v1"
	"github.com/mxplusb/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// accountCreateCmd represents the accountCreate command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create an account",
	Run:   createAccount,
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
	//accountCreateCmd.PersistentFlags().StringVar(&accountOwner, "owner", "", "the email owning the account")
}

func createAccount(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	if accountId == 0 {
		logger.Fatal().Msg("account id cannot be zero")
	}

	logger.Debug().Str("host", config.GetString("client.grpcAddr")).Msg("creating client")

	targetHost, err := url.Parse(config.GetString("client.grpcAddr"))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't parse remote host")
	}

	var host kvstorev1connect.KvStoreServiceClient
	if targetHost.Scheme != "https" {
		host = kvstorev1connect.NewKvStoreServiceClient(newInsecureClient(),targetHost.String())
	} else {
		host = kvstorev1connect.NewKvStoreServiceClient(http.DefaultClient, targetHost.String())
	}

	descriptor, err := host.CreateAccount(context.Background(), connect.NewRequest(&kvstorev1.CreateAccountRequest{AccountId: accountId, Owner: accountOwner}))
	if err != nil {
		logger.Fatal().Err(err).Uint64("account-id", accountId).Msg("can't create account")
	}

	print(protojson.Format(descriptor.Msg))
}
