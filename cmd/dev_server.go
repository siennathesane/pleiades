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
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mxplusb/cliflags/gen/gpflag"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run a development server",
	Long: `runs a development server.

it will boot with 256 predefined shards, configured in 
insecure mode, and will generally be buggy. it will run
the latest and greatest, which means it may or may not 
be usable for consuming applications. there may be unversioned
changes in this command which are not available as part of
the cloud offering. this command is unsupported beyond 
filing bugs against it the team may or may not get to

DO NOT USE THIS IN PRODUCTION`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	devCmd.AddCommand(serverCmd)

	cfg := configuration.DefaultConfiguration()
cfg.Host.MutualTLS = false

	// if we're on a mac, set different paths for the default config
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		rootDir := filepath.Join(dir, "Library", "pleiades")
		cfg.Host.DataDir = filepath.Join(rootDir, "logs")
		cfg.Host.LogDir = filepath.Join(rootDir, "shards")
		cfg.Host.CaFile = ""
		cfg.Host.CertFile = ""
		cfg.Host.KeyFile = ""

		cfg.Datastore.BasePath = filepath.Join(rootDir, "pleiades")
	}

	if err := gpflag.ParseTo(cfg, serverCmd.Flags()); err != nil {
		log.Logger.Err(err).Msg("cannot properly parse command strings")
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {} else {
			log.Logger.Error().Err(err).Msg("configuration file not found")
		}
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
