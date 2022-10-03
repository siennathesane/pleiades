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
}

func init() {
	rootCmd.AddCommand(devCmd)

	if config == nil {
		config = configuration.Get()
	}

	defaultDataBasePath := ""
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		defaultDataBasePath = filepath.Join(dir, "Library", "pleiades")
	} else {
		defaultDataBasePath = configuration.DefaultBaseDataPath
	}

	//serverConfig := &configuration.ServerConfig{
	//	Datastore: &configuration.Datastore{
	//		BasePath: defaultDataBasePath,
	//		LogDir:   filepath.Join(defaultDataBasePath, "logs"),
	//		DataDir:  filepath.Join(defaultDataBasePath, "data"),
	//	},
	//	Host: &configuration.Host{
	//		CaFile:            filepath.Join(defaultDataBasePath, "tls", "ca.pem"),
	//		CertFile:          filepath.Join(defaultDataBasePath, "tls", "cert.pem"),
	//		DeploymentId:      1,
	//		GrpcListenAddress: "0.0.0.0:5000",
	//		KeyFile:           filepath.Join(defaultDataBasePath, "tls", "key.pem"),
	//		ListenAddress:     "0.0.0.0:5001",
	//		MutualTLS:         false,
	//		NotifyCommit:      false,
	//		Rtt:               1,
	//	},
	//}

	// mtls settings
	//region
	devCmd.PersistentFlags().String("ca-cert", filepath.Join(defaultDataBasePath, "tls", "ca.pem"), "mtls ca")
	config.BindPFlag("server.host.caCert", devCmd.PersistentFlags().Lookup("ca-cert"))

	devCmd.PersistentFlags().String("cert-file", filepath.Join(defaultDataBasePath, "tls", "cert.pem"), "mtls cert")
	config.BindPFlag("server.host.certFile", devCmd.PersistentFlags().Lookup("cert-file"))

	devCmd.PersistentFlags().String("cert-key", filepath.Join(defaultDataBasePath, "tls", "key.pem"), "mtls key")
	config.BindPFlag("server.host.certKey", devCmd.PersistentFlags().Lookup("cert-key"))

	devCmd.PersistentFlags().Bool("mtls", false, "enable mtls")
	config.BindPFlag("server.host.mtls", devCmd.PersistentFlags().Lookup("mtls"))

	devCmd.MarkFlagsRequiredTogether("mtls", "ca-cert", "cert-file", "cert-key")
	//endregion
}
