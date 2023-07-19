/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cli

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
	_ cli.Command             = (*AccountCreateCommand)(nil)
	_ cli.CommandAutocomplete = (*AccountCreateCommand)(nil)
)

type AccountCreateCommand struct {
	*BaseCommand

	flagAcountOwner string
	flagAccountId   uint64
}

func (a *AccountCreateCommand) Flags() *FlagSets {
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

func (a *AccountCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (a *AccountCreateCommand) AutocompleteFlags() complete.Flags {
	return a.Flags().Completions()
}

func (a *AccountCreateCommand) Help() string {
	helpText := `Create an account.

` + a.Flags().Help()

	return helpText
}

func (a *AccountCreateCommand) Run(args []string) int {
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

	descriptor, err := client.CreateAccount(context.Background(), connect.NewRequest(&kvpb.CreateAccountRequest{
		AccountId:   a.flagAccountId,
		Owner:       a.flagAcountOwner,
		Transaction: nil,
	}))
	if err != nil {
		a.UI.Error(fmt.Sprintf("error creating account: %s", err))
		return exitCodeRemote
	}

	if descriptor.Msg != nil {
		OutputData(a.UI, descriptor.Msg)
	}

	return exitCodeGood
}

func (a *AccountCreateCommand) Synopsis() string {
	return "Create an account in the key value store"
}
