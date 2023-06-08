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

	kvstorev1 "github.com/mxplusb/pleiades/pkg/api/kvstore/v1"
	"github.com/mxplusb/pleiades/pkg/api/kvstore/v1/kvstorev1connect"
	"github.com/bufbuild/connect-go"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AccountDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*AccountDeleteCommand)(nil)
)

type AccountDeleteCommand struct {
	*BaseCommand

	flagAcountOwner string
	flagAccountId   uint64
}

func (a *AccountDeleteCommand) Flags() *FlagSets {
	set := a.flagSet(FlagSetHTTP | FlagSetFormat | FlagSetLogging)
	f := set.NewFlagSet("Account Options")

	f.Uint64Var(&Uint64Var{
		Name:              "account-id",
		Usage:             "Account ID to create",
		Target:            &a.flagAccountId,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.account.create.account-id",
	})

	f.StringVar(&StringVar{
		Name:              "account-owner-email",
		Usage:             "Email address of the account owner",
		Target:            &a.flagAcountOwner,
		Completion:        complete.PredictNothing,
		ConfigurationPath: "client.account.create.account-owner-email",
	})

	return set
}

func (a *AccountDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (a *AccountDeleteCommand) AutocompleteFlags() complete.Flags {
	return a.Flags().Completions()
}

func (a *AccountDeleteCommand) Help() string {
	helpText := `Delete an account.

` + a.Flags().Help()

	return helpText
}

func (a *AccountDeleteCommand) Run(args []string) int {
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

	client := kvstorev1connect.NewKvStoreServiceClient(httpClient, a.BaseCommand.flagHost)

	descriptor, err := client.DeleteAccount(context.Background(), connect.NewRequest(&kvstorev1.DeleteAccountRequest{
		AccountId:   a.flagAccountId,
		Owner:       a.flagAcountOwner,
		Transaction: nil,
	}))
	if err != nil {
		a.UI.Error(fmt.Sprintf("error deleting account: %s", err))
		return exitCodeRemote
	}

	if descriptor.Msg != nil {
		OutputData(a.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (a *AccountDeleteCommand) Synopsis() string {
	return "Delete an account from the key value store"
}
