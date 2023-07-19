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
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/mxplusb/pleiades/pkg/kvpb"
	"github.com/mxplusb/pleiades/pkg/kvpb/kvpbconnect"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*BucketDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*BucketDeleteCommand)(nil)
)

type BucketDeleteCommand struct {
	*BaseCommand

	flagAccountId   uint64
	flagBucketName  string
}

func (a *BucketDeleteCommand) Flags() *FlagSets {
	set := a.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging)
	f := set.NewFlagSet("Account Options")

	f.Uint64Var(&Uint64Var{
		Name:              "account-id",
		Usage:             "Account ID to associate the bucket with.",
		Target:            &a.flagAccountId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.bucket.create.account-id",
	})

	f.StringVar(&StringVar{
		Name:              "bucket",
		Usage:             "Name of the bucket.",
		Target:            &a.flagBucketName,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.bucket.create.name",
	})

	return set
}

func (a *BucketDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (a *BucketDeleteCommand) AutocompleteFlags() complete.Flags {
	return a.Flags().Completions()
}

func (a *BucketDeleteCommand) Help() string {
	helpText := `Create a bucket.

` + a.Flags().Help()

	return helpText
}

func (a *BucketDeleteCommand) Run(args []string) int {
	f := a.Flags()

	if err := f.Parse(args); err != nil {
		a.UI.Error(err.Error())
		return exitCodeFailureToParseArgs
	}

	trace := config.GetBool("logging.trace")
	if trace {
		OutputData(a.UI, config.AllSettings())
	}

	if a.flagAccountId == 0 {
		a.UI.Error("the account-id must not be 0")
		return exitCodeGenericBad
	}

	httpClient, err := a.Client()
	if err != nil {
		a.UI.Error(err.Error())
		return exitCodeGenericBad
	}

	client := kvpbconnect.NewKvStoreServiceClient(httpClient, a.BaseCommand.flagHost)

	descriptor, err := client.DeleteBucket(context.Background(), connect.NewRequest(&kvpb.DeleteBucketRequest{
		AccountId:   a.flagAccountId,
		Name:        a.flagBucketName,
		Transaction: nil,
	}))
	if err != nil {
		a.UI.Error(fmt.Sprintf("error deleting bucket: %s", err))
		return exitCodeRemote
	}

	if descriptor.Msg != nil {
		OutputData(a.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (a *BucketDeleteCommand) Synopsis() string {
	return "create an account in the key value store"
}
