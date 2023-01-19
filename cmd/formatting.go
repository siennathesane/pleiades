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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v3"
)

const (
	PleiadesDefaultOutputFormat = "json"
)

func OutputData(ui cli.Ui, data interface{}) int {
	return outputWithFormat(ui, data)
}

func outputWithFormat(ui cli.Ui, data interface{}) int {
	format := Format(ui)
	formatter, ok := Formatters[format]
	if !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	if err := formatter.Output(ui, data); err != nil {
		ui.Error(fmt.Sprintf("Could not parse output: %s", err.Error()))
		return 1
	}
	return 0
}

type Formatter interface {
	Output(ui cli.Ui, data interface{}) error
	Format(data interface{}) ([]byte, error)
}

var Formatters = map[string]Formatter{
	"json": JsonFormatter{},
	"yaml": YamlFormatter{},
	"yml":  YamlFormatter{},
	"raw":  RawFormatter{},
}

func Format(ui cli.Ui) string {
	switch ui := ui.(type) {
	case *PleiadesUI:
		return ui.format
	}

	return PleiadesDefaultOutputFormat
}

func Detailed(ui cli.Ui) bool {
	return false
}

// An output formatter for raw output of the original request object
type RawFormatter struct{}

func (r RawFormatter) Format(data interface{}) ([]byte, error) {
	byte_data, ok := data.([]byte)
	if !ok {
		return nil, fmt.Errorf("this command does not support the -format=raw option")
	}

	return byte_data, nil
}

func (r RawFormatter) Output(ui cli.Ui, data interface{}) error {
	b, err := r.Format(data)
	if err != nil {
		return err
	}
	ui.Output(string(b))
	return nil
}

// An output formatter for json output of an object
type JsonFormatter struct{}

func (j JsonFormatter) Format(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (j JsonFormatter) Output(ui cli.Ui, data interface{}) error {
	b, err := j.Format(data)
	if err != nil {
		return err
	}

	ui.Output(string(b))
	return nil
}

// An output formatter for yaml output format of an object
type YamlFormatter struct{}

func (y YamlFormatter) Format(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}

func (y YamlFormatter) Output(ui cli.Ui, data interface{}) error {
	b, err := y.Format(data)
	if err == nil {
		ui.Output(strings.TrimSpace(string(b)))
	}
	return err
}
