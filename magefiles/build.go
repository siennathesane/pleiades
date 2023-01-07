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
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

var (
	platforms  = []string{"linux", "darwin"}
	invariants = []string{"amd64", "arm64"}
)

func (Build) Setup() {
	mg.Deps(Install.Godeps)
}

// compile pleiades with the local build information
func (Build) Compile() error {
	for _, platform := range platforms {
		for _, variant := range invariants {
			fmt.Printf("##teamcity[progressMessage 'Started %s-%s compilation']\n", platform, variant)
			if err := compileWithPath(fmt.Sprintf("build/pleiades-%s-%s", platform, variant), map[string]string{
				"GOOS":   platform,
				"GOARCH": variant,
				"CGO_ENABLED": "0",
			}); err != nil {
				fmt.Printf("##teamcity[progressMessage 'Error with %s-%s compilation']\n", platform, variant)
				return err
			}
			fmt.Printf("##teamcity[progressMessage 'Finished %s-%s compilation']\n", platform, variant)
		}
	}
	return nil
}

// compile pleiades with the local build information
func compileWithPath(path string, env map[string]string) error {
	return sh.RunWith(env, "go", "build", "-v", fmt.Sprintf("-ldflags=%s", ldflags()), "-o", path, "./main.go")
}

func ldflags() string {
	fmt.Println("generating ldflags...")

	writeComma := func(sb *strings.Builder) {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
	}

	sb := strings.Builder{}

	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.GoVersion=")
	sb.WriteString(runtime.Version())
	sb.WriteString("'")
	writeComma(&sb)

	now := time.Now().Format(time.RFC3339)
	fmt.Printf("using build time: %s\n", now)

	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.BuildTime=")
	sb.WriteString(now)
	sb.WriteString("'")
	writeComma(&sb)

	version := os.Getenv("BUILD_TAG")
	if len(version) == 0 {
		x, _ := newVersion()
		version = x.String()
	}
	fmt.Printf("using version: %s\n", version)
	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.Version=")
	sb.WriteString(version)
	sb.WriteString("'")

	fmt.Printf("using ldflags: %s\n", sb.String())

	return sb.String()
}

// clean rebuild of pleiades
func (Build) Rebuild() error {
	fmt.Println("cleaning...")
	err := sh.Rm("build")
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "clean")
	err = cmd.Run()
	if err != nil {
		return err
	}
	mg.Deps(Build.Compile)
	return nil
}

// lint the repo
func (Build) Vet() error {
	fmt.Println("running linter")
	return sh.RunWithV(nil, "go", "vet", "./...")
}
