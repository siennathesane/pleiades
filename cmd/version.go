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
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/anthropos-labs/pleiades/pkg"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Pleiades",
	Run: func(cmd *cobra.Command, args []string) {
		printVersionInfo()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&jPrint, "json", false, "output in json format")
}

var (
	jPrint bool = false
)

type versionInfo struct {
	GoVersion string `json:"go_version"`
	Sha  string `json:"commit"`
	BuildTime string `json:"build_time"`
	Dirty string `json:"dirty"`
}

func printVersionInfo() {
	vi := versionInfo{
		GoVersion: pkg.GoVersion,
		Sha: pkg.Sha,
		BuildTime: pkg.BuildTime,
		Dirty: pkg.Dirty,
	}

	if jPrint {
		target, err := json.MarshalIndent(vi, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(target))
	} else {
		fmt.Printf("pleiades version: %s\n", pkg.Sha)
		fmt.Printf("go version: %s\n", pkg.GoVersion)
		fmt.Printf("build time: %s\n", pkg.BuildTime)
		fmt.Printf("dirty: %v\n", pkg.Dirty)
	}
}