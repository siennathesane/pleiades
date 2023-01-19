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

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KvDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*KvDeleteCommand)(nil)
)

type KvDeleteCommand struct {
	*BaseCommand

	flagAccountId  uint64
	flagBucketName string
	flagKey        string
}

func (k *KvDeleteCommand) Flags() *FlagSets {
	set := k.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging)
	f := set.NewFlagSet("Key Value Options")

	f.StringVar(&StringVar{
		Name:              "key",
		Usage:             "Name of the key.",
		Target:            &k.flagKey,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.put.key",
	})

	f.Uint64Var(&Uint64Var{
		Name:              "account-id",
		Usage:             "Account ID for the key.",
		Target:            &k.flagAccountId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.put.account-id",
	})

	f.StringVar(&StringVar{
		Name:              "bucket",
		Usage:             "Name of the bucket.",
		Target:            &k.flagBucketName,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kv.put.bucket-name",
	})

	return set
}

func (k *KvDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (k *KvDeleteCommand) AutocompleteFlags() complete.Flags {
	return k.Flags().Completions()
}

func (k *KvDeleteCommand) Help() string {
	helpText := `Delete a key.

` + k.Flags().Help()

	return helpText
}

func (k *KvDeleteCommand) Run(args []string) int {
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
		k.UI.Error("account-id cannot be 0")
		return exitCodeGenericBad
	}

	if k.flagBucketName == "" {
		k.UI.Error("bucket cannot be empty")
		return exitCodeGenericBad
	}

	if k.flagKey == "" {
		k.UI.Error("key cannot be empty")
		return exitCodeGenericBad
	}

	httpClient, err := k.Client()
	if err != nil {
		k.UI.Error(err.Error())
		return exitCodeGenericBad
	}

	client := kvstorev1connect.NewKvStoreServiceClient(httpClient, k.BaseCommand.flagHost)

	descriptor, err := client.DeleteKey(context.Background(), connect.NewRequest(&kvstorev1.DeleteKeyRequest{
		AccountId:  k.flagAccountId,
		BucketName: k.flagBucketName,
		Key: []byte(k.flagKey),
	}))
	if err != nil {
		k.UI.Error(err.Error())
	}

	if descriptor != nil {
		OutputData(k.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (k *KvDeleteCommand) Synopsis() string {
	return "Delete a key."
}
