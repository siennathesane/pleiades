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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mxplusb/pleiades/pkg"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pleiades",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	config *viper.Viper
)

func init() {
	if config == nil {
		config = configuration.Get()
	}

	viper.SetConfigName("pleiades") // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, _ := homedir.Dir()
		viper.AddConfigPath(filepath.Join(dir, "Library", "pleiades"))
	} else {
		viper.AddConfigPath(configuration.DefaultBaseConfigPath) // path to look for the config file in
	}

	viper.AddConfigPath("$HOME/.pleiades") // call multiple times to add many search paths
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

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
	config.Set("server.datastore.basePath", defaultDataBasePath)

	//rootCmd.PersistentFlags().Bool("trace", false, "enable trace logging")
	//config.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))
	//
	//rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	//config.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	//rootCmd.MarkFlagsMutuallyExclusive("debug", "trace")

	// mtls settings
	//region
	rootCmd.PersistentFlags().String("ca-cert", filepath.Join(defaultDataBasePath, "tls", "ca.pem"), "mtls ca")
	config.BindPFlag("server.host.caFile", rootCmd.PersistentFlags().Lookup("ca-cert"))

	rootCmd.PersistentFlags().String("cert-file", filepath.Join(defaultDataBasePath, "tls", "cert.pem"), "mtls cert")
	config.BindPFlag("server.host.certFile", rootCmd.PersistentFlags().Lookup("cert-file"))

	rootCmd.PersistentFlags().String("cert-key", filepath.Join(defaultDataBasePath, "tls", "key.pem"), "mtls key")
	config.BindPFlag("server.host.keyFile", rootCmd.PersistentFlags().Lookup("cert-key"))

	rootCmd.MarkFlagsRequiredTogether("ca-cert", "cert-file", "cert-key")
	//endregion
}

func setupLogger() zerolog.Logger {
	var logger zerolog.Logger
	if config.GetBool("logging.trace") {
		logger = configuration.NewRootLogger().Level(zerolog.TraceLevel)
	} else if config.GetBool("logging.debug") {
		logger = configuration.NewRootLogger().Level(zerolog.DebugLevel)
	} else {
		logger = configuration.NewRootLogger().Level(zerolog.InfoLevel)
	}
	return logger
}

type PleiadesUi struct {
	cli.Ui
	format string
}

// setupEnv parses args and may replace them and sets some env vars to known
// values based on format options
func setupEnv(args []string) (retArgs []string, format string) {
	var nextArgFormat bool

	for _, arg := range args {
		if nextArgFormat {
			nextArgFormat = false
			format = arg
			continue
		}

		if arg == "--" {
			break
		}

		if len(args) == 1 && (arg == "-v" || arg == "-version" || arg == "--version") {
			args = []string{"version"}
			break
		}

		// Parse a given flag here, which overrides the env var
		if isGlobalFlagWithValue(arg, globalFlagFormat) {
			format = getGlobalFlagValue(arg)
		}
		// For backwards compat, it could be specified without an equal sign
		if isGlobalFlag(arg, globalFlagFormat) {
			nextArgFormat = true
		}
	}

	envPleiadesFormat := os.Getenv(EnvPleiadesDefaultOutput)
	// If we did not parse a value, fetch the env var
	if format == "" && envPleiadesFormat != "" {
		format = envPleiadesFormat
	}
	// Lowercase for consistency
	format = strings.ToLower(format)
	if format == "" {
		format = "json"
	}

	return args, format
}

func isGlobalFlag(arg string, flag string) bool {
	return arg == "-"+flag || arg == "--"+flag
}

func isGlobalFlagWithValue(arg string, flag string) bool {
	return strings.HasPrefix(arg, "--"+flag+"=") || strings.HasPrefix(arg, "-"+flag+"=")
}

func getGlobalFlagValue(arg string) string {
	_, value, _ := strings.Cut(arg, "=")

	return value
}

type RunOptions struct {
	Stdout io.Writer
	Stderr io.Writer
}

func Run(args []string) int {
	runOpts := &RunOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	_, format := setupEnv(args)

	ui := PleiadesUi{
		Ui: &cli.BasicUi{
			Reader:      bufio.NewReader(os.Stdin),
			Writer:      runOpts.Stdout,
			ErrorWriter: runOpts.Stderr,
		},
		format: format,
	}

	if _, ok := Formatters[format]; !ok {
		ui.Error(fmt.Sprintf("invalid output format: %s", format))
		return 1
	}

	initCommands(ui)

	root := &cli.CLI{
		Args:                       args,
		Commands:                   Commands,
		Name:                       "pleiades",
		Version:                    pkg.Version,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		AutocompleteGlobalFlags:    nil,
		HelpFunc:                   groupedHelpFunc(cli.BasicHelpFunc("pleiades")),
	}

	code, err := root.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return 1
	}

	return code
}

func groupedHelpFunc(f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var b bytes.Buffer
		tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)

		fmt.Fprintf(tw, "usage: pleiades <command> [args]\n")

		otherCommands := make([]string, 0, len(commands))
		for k := range commands {
			otherCommands = append(otherCommands, k)
		}
		sort.Strings(otherCommands)

		fmt.Fprintf(tw, "\n")
		fmt.Fprintf(tw, "commands:\n")
		for _, v := range otherCommands {
			printCommand(tw, v, commands[v])
		}

		tw.Flush()

		return strings.TrimSpace(b.String())
	}
}

func printCommand(w io.Writer, name string, cmdFn cli.CommandFactory) {
	cmd, err := cmdFn()
	if err != nil {
		panic(fmt.Sprintf("failed to load %q command: %s", name, err))
	}
	fmt.Fprintf(w, "    %s\t%s\n", name, cmd.Synopsis())
}
