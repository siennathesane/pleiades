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
	"github.com/spf13/viper"
)

var (
	config   *viper.Viper
	Commands map[string]cli.CommandFactory
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
}

type PleiadesUi struct {
	cli.Ui
	format string
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

		fmt.Fprintf(tw, "Usage: pleiades <command> [args]\n")

		otherCommands := make([]string, 0, len(commands))
		for k := range commands {
			otherCommands = append(otherCommands, k)
		}
		sort.Strings(otherCommands)

		fmt.Fprintf(tw, "\n")
		fmt.Fprintf(tw, "Commands:\n")
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

func initCommands(ui cli.Ui) {
	getBaseCmd := func() *BaseCommand {
		return &BaseCommand{
			UI: ui,
		}
	}

	Commands = map[string]cli.CommandFactory{
		"kv": func() (cli.Command, error) {
			return &KvCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"kv get": func() (cli.Command, error) {
			return &KvGetCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"kv put": func() (cli.Command, error) {
			return &KvPutCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"kv delete": func() (cli.Command, error) {
			return &KvDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account": func() (cli.Command, error) {
			return &AccountCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account create": func() (cli.Command, error) {
			return &AccountCreateCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"account delete": func() (cli.Command, error) {
			return &AccountDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket": func() (cli.Command, error) {
			return &BucketCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket create": func() (cli.Command, error) {
			return &BucketCreateCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"bucket delete": func() (cli.Command, error) {
			return &BucketDeleteCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric": func() (cli.Command, error) {
			return &FabricCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric add-shard": func() (cli.Command, error) {
			return &FabricAddShardCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric add-replica": func() (cli.Command, error) {
			return &FabricAddReplicaCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric start-replica": func() (cli.Command, error) {
			return &FabricStartReplicaCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric stop-replica": func() (cli.Command, error) {
			return &FabricStopReplicaCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric start-replica-observer": func() (cli.Command, error) {
			return &FabricStartReplicaObserverCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric remove-replica": func() (cli.Command, error) {
			return &FabricRemoveReplicaCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric remove-data": func() (cli.Command, error) {
			return &FabricRemoveDataCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric add-replica-observer": func() (cli.Command, error) {
			return &FabricAddReplicaObserverCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric add-replica-witness": func() (cli.Command, error) {
			return &FabricAddReplicaWitnessCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric get-leader-id": func() (cli.Command, error) {
			return &FabricGetLeaderIdCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
		"fabric get-shard-members": func() (cli.Command, error) {
			return &FabricGetShardMembersCommand{
				BaseCommand: getBaseCmd(),
			}, nil
		},
	}
}
