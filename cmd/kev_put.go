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
	"context"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/kvpb/kvpbconnect"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KvPutCommand)(nil)
	_ cli.CommandAutocomplete = (*KvPutCommand)(nil)
)

type KvPutCommand struct {
	*BaseCommand

	flagAccountId  uint64
	flagBucketName string
	flagKey        string
	flagKeyVersion uint32
	flagKeyFile    string
}

func (k *KvPutCommand) Flags() *FlagSets {
	set := k.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging)
	f := set.NewFlagSet("Key Value Options")

	f.StringVar(&StringVar{
		Name:              "key",
		Usage:             "Name of the key.",
		Target:            &k.flagKey,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kvpb.put.key",
	})

	f.StringVar(&StringVar{
		Name:              "value-from-file",
		Usage:             "Local filepath of the key data.",
		Target:            &k.flagKeyFile,
		Completion:        complete.PredictFiles("*"),
		ConfigurationPath: "client.kvpb.put.key-file-path",
	})

	f.Uint32Var(&Uint32Var{
		Name: "key-version",
		Usage: `Version of the key. Pleiades implements monotonic versions, so the key version must be 
higher than the exising key version value.`,
		Target:            &k.flagKeyVersion,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kvpb.put.key-version",
	})

	f.Uint64Var(&Uint64Var{
		Name:              "account-id",
		Usage:             "Account ID for the key.",
		Target:            &k.flagAccountId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kvpb.put.account-id",
	})

	f.StringVar(&StringVar{
		Name:              "bucket",
		Usage:             "Name of the bucket.",
		Target:            &k.flagBucketName,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.kvpb.put.bucket-name",
	})

	return set
}

func (k *KvPutCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (k *KvPutCommand) AutocompleteFlags() complete.Flags {
	return k.Flags().Completions()
}

func (k *KvPutCommand) Help() string {
	helpText := `Put a key.

` + k.Flags().Help()

	return helpText
}

func (k *KvPutCommand) Run(args []string) int {
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

	if k.flagKeyFile == "" {
		k.UI.Error("value-from-file must not be empty")
		return exitCodeGenericBad
	}

	payload, err := os.ReadFile(k.flagKeyFile)
	if err != nil {
		k.UI.Error(err.Error())
		return exitCodeGenericBad
	}

	httpClient, err := k.Client()
	if err != nil {
		k.UI.Error(err.Error())
		return exitCodeGenericBad
	}

	client := kvpbconnect.NewKvStoreServiceClient(httpClient, k.BaseCommand.flagHost)

	now := time.Now().UnixMilli()

	descriptor, err := client.PutKey(context.Background(), connect.NewRequest(&kvpb.PutKeyRequest{
		AccountId:  k.flagAccountId,
		BucketName: k.flagBucketName,
		KeyValuePair: &kvpb.KeyValue{
			Key:            []byte(k.flagKey),
			CreateRevision: now,
			ModRevision:    now,
			Version:        k.flagKeyVersion,
			Value:          payload,
			Lease:          0,
		}}))
	if err != nil {
		k.UI.Error(err.Error())
	}

	if descriptor != nil {
		OutputData(k.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (k *KvPutCommand) Synopsis() string {
	return "Put a key."
}
