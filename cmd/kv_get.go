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

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// kvGetCmd represents the kvGet command
var (
	kvGetCmd = &cobra.Command{
		Use:   "get",
		Short: "get a key",
		Run:   getKey,
	}
)

func init() {
	kvCmd.AddCommand(kvGetCmd)

	kvGetCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to look for")
}

func getKey(cmd *cobra.Command, args []string) {
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

	descriptor, err := client.GetKey(context.Background(), connect.NewRequest(&kvstorev1.GetKeyRequest{
		AccountId:  accountId,
		BucketName: bucketName,
		Key:        []byte(key),
	}))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't delete bucket")
	}

	print(protojson.Format(descriptor.Msg))
}

var (
	_ cli.Command             = (*KvGetCommand)(nil)
	_ cli.CommandAutocomplete = (*KvGetCommand)(nil)
)

type KvGetCommand struct {
	*BaseCommand

	flagAccountId  uint64
	flagBucket     string
	flagKey        string
	flagKeyVersion uint32
}

func (k *KvGetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (k *KvGetCommand) AutocompleteFlags() complete.Flags {
	return k.Flags().Completions()
}

func (k *KvGetCommand) Help() string {
	helpText := `Get a key from as specific bucket.

` + k.Flags().Help()
	return helpText
}

func (k *KvGetCommand) Flags() *FlagSets {
	set := k.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging)
	f := set.NewFlagSet("Key Value Options")

	f.Uint64Var(&Uint64Var{
		Name:              "account-id",
		Usage:             "Account ID where the bucket resides.",
		Target:            &k.flagAccountId,
		Default:           0,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.get.account-id",
	})

	f.StringVar(&StringVar{
		Name:              "bucket",
		Usage:             "The bucket where key information can be found.",
		Default:           "",
		Target:            &k.flagBucket,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.get.bucket",
	})

	f.StringVar(&StringVar{
		Name:              "key",
		Usage:             "The specific key to fetch.",
		Default:           "",
		Target:            &k.flagKey,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.get.key",
	})

	f.Uint32Var(&Uint32Var{
		Name: "key-version",
		Usage: `The specific key revision to retrieve. This is an optional field and will fetch the latest value
if omitted.`,
		Default:           0,
		Target:            &k.flagKeyVersion,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.get.key-version",
	})

	return set
}

func (k *KvGetCommand) Run(args []string) int {
	f := k.Flags()

	if err := f.Parse(args); err != nil {
		k.UI.Error(err.Error())
		return exitCodeFailureToParseArgs
	}

	trace := config.GetBool("logging.trace")
	if trace {
		OutputData(k.UI, config.AllSettings())
	}

	if k.flagAccountId == 0 {
		k.UI.Error("account id cannot be zero")
		return exitCodeFailureToParseArgs
	}

	if k.flagBucket == "" {
		k.UI.Error("bucket name cannot be empty")
		return exitCodeFailureToParseArgs
	}

	if k.flagKey == "" {
		k.UI.Error("key cannot be empty")
		return exitCodeFailureToParseArgs
	}

	httpClient, err := k.Client()
	if err != nil {
		k.UI.Error(err.Error())
		return exitCodeFailureToParseArgs
	}

	client := kvstorev1connect.NewKvStoreServiceClient(httpClient, k.BaseCommand.flagHost)

	descriptor, err := client.GetKey(context.Background(), connect.NewRequest(&kvstorev1.GetKeyRequest{
		AccountId:  k.flagAccountId,
		BucketName: k.flagBucket,
		Key:        []byte(k.flagKey),
		Version:    &k.flagKeyVersion,
	}))

	if err != nil {
		k.UI.Error(err.Error())
		return exitCodeRemote
	}

	if descriptor.Msg != nil {
		OutputData(k.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (k *KvGetCommand) Synopsis() string {
	return "Fetch a key from a bucket"
}
