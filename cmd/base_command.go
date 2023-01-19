/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/kr/text"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"golang.org/x/net/http2"
)

type (
	FlagSetBit                uint
	ParseOptions              interface{}
	ParseOptionAllowRawFormat bool
)

const (
	FlagSetNone FlagSetBit = 1 << iota
	FlagSetLogging
	FlagSetHTTP
	FlagSetFormat
	FlagSetTls
	FlagSetTimeout

	globalFlagFormat = "format"
)

var (
	reRemoveWhitespace = regexp.MustCompile(`[\s]+`)
	maxLineLength      = 80

	globalFlags []string = []string{globalFlagFormat}
)

type BaseCommand struct {
	UI cli.Ui

	flags     *FlagSets
	flagsOnce sync.Once

	flagHost          string
	flagCACert        string
	flagTlsSkipVerify bool
	flagDebug         bool
	flagTrace         bool
	flagFormat        string
	flagTimeout int32

	client *http.Client
}

func (b *BaseCommand) Client() (*http.Client, error) {
	if b.client != nil {
		return b.client, nil
	}

	parsedHost, err := url.Parse(b.flagHost)
	if err != nil {
		return nil, err
	}

	if parsedHost.Scheme != "https" {
		return newInsecureClient(), nil
	}

	b.client = http.DefaultClient

	if b.flagCACert != "" {
		cfg, err := newTlsConfig(b, parsedHost.Host)
		if err != nil {
			return nil, err
		}

		b.client.Transport = &http2.Transport{
			AllowHTTP:       true,
			TLSClientConfig: cfg,
		}
	}

	return b.client, nil
}

func (b *BaseCommand) flagSet(bit FlagSetBit) *FlagSets {
	b.flagsOnce.Do(func() {
		set := NewFlagSets(b.UI)

		if bit&FlagSetTimeout != 0 {
			f := set.NewFlagSet("Timeout Options")

			f.Int32Var(&Int32Var{
				Name:              "timeout",
				Usage:             "Set the timeout for the command runtime, in seconds.",
				Default:           30,
				Target:            &b.flagTimeout,
				Completion:        complete.PredictNothing,
				ConfigurationPath: "client.timeout",
			})
		}

		if bit&FlagSetHTTP != 0 {
			f := set.NewFlagSet("HTTP Options")

			f.StringVar(&StringVar{
				Name:       flagNameHost,
				Usage:      "Address of pleiades server.",
				Default:    "http://localhost:8080",
				Hidden:     false,
				EnvVar:     EnvPleiadesUrl,
				Target:     &b.flagHost,
				Completion: complete.PredictAnything,
				ConfigurationPath: "client.http.address",
			})

			f.BoolVar(&BoolVar{
				Name:    "tls-skip-verify",
				Usage:   "Disable TLS SNI checking.",
				Default: false,
				Hidden:  false,
				EnvVar:  EnvPleiadesInsecureSkipVerify,
				ConfigurationPath: "tls.skip-verify",
				Target:  &b.flagTlsSkipVerify,
			})
		}

		if bit&FlagSetTls != 0 {
			f := set.NewFlagSet("TLS Options")

			f.StringVar(&StringVar{
				Name:       "ca-cert-file",
				Usage:      `Local on-disk path to a PEM-encoded CA certificate if using a custom TLS certificate.`,
				Default:    "",
				EnvVar:     EnvPleiadesCaCert,
				Target:     &b.flagCACert,
				ConfigurationPath: "tls.ca-cert-file",
				Completion: complete.PredictFiles("*"),
			})

			f.StringVar(&StringVar{
				Name:       "cert-file",
				Usage:      "Local on-disk path to a PEM-encoded certificate.",
				Default:    "",
				EnvVar:     EnvPleiadesCertFile,
				Target:     &b.flagCACert,
				ConfigurationPath: "tls.cert-file",
				Completion: complete.PredictFiles("*"),
			})

			f.StringVar(&StringVar{
				Name:       "key-file",
				Usage:      "Local on-disk path to a PEM-encoded certificate.",
				Default:    "",
				EnvVar:     EnvPleiadesKeyFile,
				Target:     &b.flagCACert,
				ConfigurationPath: "tls.key-file",
				Completion: complete.PredictFiles("*"),
			})
		}

		if bit&FlagSetLogging != 0 {
			logSet := set.NewFlagSet("Logging Options")

			logSet.BoolVar(&BoolVar{
				Name:    "debug",
				Usage:   "Enable debug logging. WARNING: this will output a lot of data.",
				Default: false,
				Hidden:  false,
				EnvVar:  EnvPleiadesDebug,
				Target:  &b.flagDebug,
				ConfigurationPath: "logging.debug",
			})

			logSet.BoolVar(&BoolVar{
				Name: "trace",
				Usage: "Enable trace logging. WARNING: this will substantially slow down pleiades and generate " +
					"extensive amounts of logging data. Never use this in a production environment.",
				Default: false,
				Hidden:  false,
				EnvVar:  EnvPleiadesTrace,
				Target:  &b.flagDebug,
				ConfigurationPath: "logging.trace",
			})
		}

		if bit&FlagSetFormat != 0 {
			outputSet := set.NewFlagSet("Output Options")

			outputSet.StringVar(&StringVar{
				Name:       "format",
				Usage:      "Print the output in the target format. The options are 'json', 'raw', 'yaml', or 'yml'.",
				Default:    "json",
				EnvVar:     EnvPleiadesDefaultOutput,
				Target:     &b.flagFormat,
				ConfigurationPath: "logging.clientFormat",
				Completion: complete.PredictSet("json", "yaml", "yml", "raw"),
			})
		}

		b.flags = set
	})

	return b.flags
}

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}

func newTlsConfig(b *BaseCommand, hostname string) (*tls.Config, error) {
	crt, err := os.ReadFile(b.flagCACert)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	return &tls.Config{
		RootCAs:            rootCAs,
		ServerName:         hostname,
		InsecureSkipVerify: b.flagTlsSkipVerify,
	}, nil
}

// FlagSets is a group of flag sets.
type FlagSets struct {
	flagSets    []*FlagSet
	mainSet     *flag.FlagSet
	hiddens     map[string]struct{}
	completions complete.Flags
	ui          cli.Ui
}

// NewFlagSets creates a new flag sets.
func NewFlagSets(ui cli.Ui) *FlagSets {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	// Errors and usage are controlled by the CLI.
	mainSet.Usage = func() {}
	mainSet.SetOutput(ioutil.Discard)

	return &FlagSets{
		flagSets:    make([]*FlagSet, 0, 6),
		mainSet:     mainSet,
		hiddens:     make(map[string]struct{}),
		completions: complete.Flags{},
		ui:          ui,
	}
}

// NewFlagSet creates a new flag set from the given flag sets.
func (f *FlagSets) NewFlagSet(name string) *FlagSet {
	flagSet := NewFlagSet(name)
	flagSet.mainSet = f.mainSet
	flagSet.completions = f.completions
	f.flagSets = append(f.flagSets, flagSet)
	return flagSet
}

// Completions returns the completions for this flag set.
func (f *FlagSets) Completions() complete.Flags {
	return f.completions
}

// Parse parses the given flags, returning any errors.
// Warnings, if any, regarding the arguments format are sent to stdout
func (f *FlagSets) Parse(args []string, opts ...ParseOptions) error {
	err := f.mainSet.Parse(args)

	warnings := generateFlagWarnings(f.Args())
	if warnings != "" && Format(f.ui) == "table" {
		f.ui.Warn(warnings)
	}

	if err != nil {
		return err
	}

	// Now surface any other errors.
	return generateFlagErrors(f, opts...)
}

// Parsed reports whether the command-line flags have been parsed.
func (f *FlagSets) Parsed() bool {
	return f.mainSet.Parsed()
}

// Args returns the remaining args after parsing.
func (f *FlagSets) Args() []string {
	return f.mainSet.Args()
}

// Visit visits the flags in lexicographical order, calling fn for each. It
// visits only those flags that have been set.
func (f *FlagSets) Visit(fn func(*flag.Flag)) {
	f.mainSet.Visit(fn)
}

// Help builds custom help for this command, grouping by flag set.
func (f *FlagSets) Help() string {
	var out bytes.Buffer

	for _, set := range f.flagSets {
		printFlagTitle(&out, set.name+":")
		set.VisitAll(func(f *flag.Flag) {
			// Skip any hidden flags
			if v, ok := f.Value.(FlagVisibility); ok && v.Hidden() {
				return
			}
			printFlagDetail(&out, f)
		})
	}

	return strings.TrimRight(out.String(), "\n")
}

type FlagSet struct {
	name        string
	flagSet     *flag.FlagSet
	mainSet     *flag.FlagSet
	completions complete.Flags
}

// NewFlagSet creates a new flag set.
func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		name:    name,
		flagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

// Name returns the name of this flag set.
func (f *FlagSet) Name() string {
	return f.name
}

func (f *FlagSet) Visit(fn func(*flag.Flag)) {
	f.flagSet.Visit(fn)
}

func (f *FlagSet) VisitAll(fn func(*flag.Flag)) {
	f.flagSet.VisitAll(fn)
}

// printFlagTitle prints a consistently-formatted title to the given writer.
func printFlagTitle(w io.Writer, s string) {
	fmt.Fprintf(w, "%s\n\n", s)
}

// printFlagDetail prints a single flag to the given writer.
func printFlagDetail(w io.Writer, f *flag.Flag) {
	// Check if the flag is hidden - do not print any flag detail or help output
	// if it is hidden.
	if h, ok := f.Value.(FlagVisibility); ok && h.Hidden() {
		return
	}

	// Check for a detailed example
	example := ""
	if t, ok := f.Value.(FlagExample); ok {
		example = t.Example()
	}

	if example != "" {
		fmt.Fprintf(w, "  -%s=<%s>\n", f.Name, example)
	} else {
		fmt.Fprintf(w, "  -%s\n", f.Name)
	}

	usage := reRemoveWhitespace.ReplaceAllString(f.Usage, " ")
	indented := wrapAtLengthWithPadding(usage, 6)
	fmt.Fprintf(w, "%s\n\n", indented)
}

// wrapAtLengthWithPadding wraps the given text at the maxLineLength, taking
// into account any provided left padding.
func wrapAtLengthWithPadding(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}

func generateFlagWarnings(args []string) string {
	var trailingFlags []string
	for _, arg := range args {
		// "-" can be used where a file is expected to denote stdin.
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			continue
		}

		isGlobalFlag := false
		trimmedArg, _, _ := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		for _, flag := range globalFlags {
			if trimmedArg == flag {
				isGlobalFlag = true
			}
		}
		if isGlobalFlag {
			continue
		}

		trailingFlags = append(trailingFlags, arg)
	}

	if len(trailingFlags) > 0 {
		return fmt.Sprintf("Command flags must be provided before positional arguments. "+
			"The following arguments will not be parsed as flags: [%s]", strings.Join(trailingFlags, ","))
	} else {
		return ""
	}
}

func generateFlagErrors(f *FlagSets, opts ...ParseOptions) error {
	if Format(f.ui) == "raw" {
		canUseRaw := false
		for _, opt := range opts {
			if value, ok := opt.(ParseOptionAllowRawFormat); ok {
				canUseRaw = bool(value)
			}
		}

		if !canUseRaw {
			return fmt.Errorf("This command does not support the -format=raw option.")
		}
	}

	return nil
}

type PleiadesUI struct {
	cli.Ui
	format string
}
