//go:build mage

/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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
			if err := compileWithPath(fmt.Sprintf("build/pleiades-%s-%s", platform, variant), map[string]string{
				"GOOS":   platform,
				"GOARCH": variant,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// compile pleiades with the local build information
func compileWithPath(path string, env map[string]string) error {
	return sh.RunWithV(env, "go", "build", "-v", fmt.Sprintf("-ldflags=%s", ldflags()), "-o", path, "./main.go")
}

func ldflags() string {
	fmt.Println("generating ldflags...")

	writeComma := func(sb *strings.Builder) {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
	}

	fmt.Println("getting git head...")
	localRepo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	head, err := localRepo.Head()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	fmt.Printf("got git head: %s\n", head.Hash().String())
	headHash := head.Hash().String()

	worktreeStatus, err := localRepo.Worktree()
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	status, err := worktreeStatus.Status()
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	var dirtyHead bool
	if status.IsClean() {
		dirtyHead = false
	} else {
		dirtyHead = true
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

	shortHead := headHash[len(headHash)-7:]
	fmt.Printf("using git hash: %s\n", shortHead)
	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.Sha=")
	sb.WriteString(shortHead)
	sb.WriteString("'")

	fmt.Printf("is head dirty: %v\n", dirtyHead)
	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.Dirty=")
	sb.WriteString(strconv.FormatBool(dirtyHead))
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
