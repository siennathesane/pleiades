/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

// build the docker image
func (Docker) Build() error {
	if err := compileWithPath("build/pleiades", map[string]string{
		"GOARCH": "arm64",
		"GOOS":   "linux",
	}); err != nil {
		return err
	}

	if err := sh.RunWithV(nil, "docker", "buildx", "build", "-t", "registry.gitlab.com/anthropos-labs/pleiades/pleiades", "."); err != nil {
		return err
	}

	return compileWithPath("build/pleiades", map[string]string{
		"GOARCH": "arm64",
		"GOOS":   "darwin",
	})
}

func (Docker) Run() error {
	mg.Deps(Docker.Build)

	args := append(make([]string, 0), []string{
		"run",
		"-p",
		"8080:8080",
		"-p",
		"8081:8081",
		"registry.gitlab.com/anthropos-labs/pleiades/pleiades",
		"server",
		"--debug",
		"--round-trip",
		"10"}...)
	return sh.RunWithV(nil, "docker", args...)
}
