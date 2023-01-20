/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

//go:build mage

package main

import (
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

const (
	containerImageName string = "gcr.io/pleiades-353402/pleiades"
)

// build the docker image
func (Docker) Build() error {
	mg.Deps(Build.Compile)

	//goland:noinspection GoBoolExpressions
	if runtime.GOARCH == "arm64" {
		if err := sh.RunWith(nil, "docker", "buildx", "build", "--platform", "linux/amd64", "-t", containerImageName, "."); err != nil {
			return err
		}
	} else {
		if err := sh.RunWith(nil, "docker", "build", "-t", containerImageName, "."); err != nil {
			return err
		}
	}

	return nil
}

// build and push the docker image
func (Docker) Push() error {
	mg.Deps(Docker.Build)

	if err := sh.RunWith(nil, "docker", "push", containerImageName); err != nil {
		return err
	}

	return nil
}

func (Docker) Run() error {
	mg.Deps(Docker.Build)

	args := append(make([]string, 0), []string{
		"run",
		"-p",
		"8080:8080",
		"-p",
		"8081:8081",
		containerImageName,
		"server",
		"--debug",
		"--round-trip",
		"10"}...)
	return sh.RunWith(nil, "docker", args...)
}
