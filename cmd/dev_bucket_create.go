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

// bucketCreateCmd represents the bucketGet command
var bucketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new bucket",
	Run:   bucketCreate,
}

func init() {
	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCreateCmd.PersistentFlags().StringVar(&accountOwner, "owner", "", "the email owning the bucket")
	bucketCreateCmd.PersistentFlags().StringVar(&bucketName, "name", "", "name of the bucket")
}

func bucketCreate(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Fatal().Err(err).Msg("can't parse flags")
	}

	logger := setupLogger(cmd, args)

	logger = logger.With().Uint64("account-id", accountId).Str("bucket", bucketName).Logger()

	if accountId == 0 {
		logger.Fatal().Msg("account id cannot be zero")
	}

	client := kvstorev1connect.NewKvStoreServiceClient(http.DefaultClient, "http://localhost:8080")

	descriptor, err := client.CreateBucket(context.Background(), connect.NewRequest(&kvstorev1.CreateBucketRequest{AccountId: accountId, Owner: accountOwner, Name: bucketName}))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't create bucket")
	}

	print(protojson.Format(descriptor.Msg))
}
