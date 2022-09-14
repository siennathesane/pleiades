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
	"path/filepath"
	"runtime"

	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "the development commands",
	Long: `these commands are used to run pleiades in various development modes.

most of them are hacky, designed to test or implement various
features, functionalities, or other various things used by both
the pleiades development team, and developers at large. these 
commands can't be trusted for anything other than development 
purposes.

DO NOT USE THEM FOR PRODUCTION`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dev called")
	},
}

var (
	serverConfig    *configuration.ServerConfig
)

func init() {
	rootCmd.AddCommand(devCmd)

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		defaultBasePath = filepath.Join(dir, "Library", "pleiades")
	}

	serverConfig = &configuration.ServerConfig{
		Datastore: &configuration.Datastore{
			BasePath: defaultBasePath,
			LogDir:   filepath.Join(defaultBasePath, "logs"),
			DataDir:  filepath.Join(defaultBasePath, "data"),
		},
		Host: &configuration.Host{
			CaFile:        filepath.Join(defaultBasePath, "tls", "ca.pem"),
			CertFile:      filepath.Join(defaultBasePath, "tls", "cert.pem"),
			DeploymentId:  1,
			GrpcListenAddress: "0.0.0.0:5000",
			KeyFile:       filepath.Join(defaultBasePath, "tls", "key.pem"),
			ListenAddress: "0.0.0.0:5001",
			MutualTLS:     false,
			NotifyCommit:  false,
			Rtt:           1,
		},
	}
	config.Server = serverConfig

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		serverConfig.Host.GrpcListenAddress = "0.0.0.0:50000"
	}

	devCmd.PersistentFlags().StringVar(&serverConfig.Host.CaFile, "ca-cert", serverConfig.Host.CaFile, "tls ca")
	devCmd.PersistentFlags().StringVar(&serverConfig.Host.CertFile, "cert", serverConfig.Host.CertFile, "mtls cert")
	devCmd.PersistentFlags().StringVar(&serverConfig.Host.KeyFile, "cert-key", serverConfig.Host.KeyFile, "mtls key")
	devCmd.PersistentFlags().BoolVar(&serverConfig.Host.MutualTLS, "mtls", serverConfig.Host.MutualTLS, "enable mtls")

	devCmd.MarkFlagsRequiredTogether("mtls", "ca-cert", "cert", "cert-key")

	c := filepath.Join(configuration.DefaultBaseConfigPath, "pleiades.yaml")
	devCmd.PersistentFlags().StringP("config", "c", c, "config file location")
}
