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
	"time"

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mxplusb/pleiades/pkg/server"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// bucketCreateCmd represents the bucketGet command
var bucketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new bucket",
	Run: bucketCreate,
}

func init() {
	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCreateCmd.PersistentFlags().StringVar(&accountOwner, "owner", "", "the email owning the bucket")
	bucketCreateCmd.PersistentFlags().StringVar(&bucketName, "name", "", "name of the bucket")
}

func bucketCreate(cmd *cobra.Command, args []string) {
	err := cmd.Flags().Parse(args)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("can't parse flags")
	}

	var logger zerolog.Logger
	if debug {
		logger = configuration.NewRootLogger().Level(zerolog.DebugLevel)
	} else {
		logger = configuration.NewRootLogger()
	}
	logger = logger.With().Uint64("account-id", accountId).Str("bucket", bucketName).Logger()

	if accountId == 0 {
		logger.Fatal().Msg("account id cannot be zero")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Fatal().Err(err).Msg("error dialing server")
	}

	client := server.NewKVStoreServiceClient(conn)

	logger.Info().Msg("creating bucket")
	descriptor, err := client.CreateBucket(context.Background(), &database.CreateBucketRequest{AccountId: accountId, Owner: accountOwner, Name: bucketName})
	if err != nil {
		logger.Fatal().Err(err).Msg("can't create bucket")
	}

	print(proto.MarshalTextString(descriptor))
}
