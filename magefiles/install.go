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
	"os"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

var (
	homebrewTargets = []string{
		"vault",
		"yq",
		"nvm",
		"krew",
		"kubernetes-cli",
		"graphviz",
	}
)

type Install mg.Namespace

// install pleiades to your local directory
func (Install) Local() error {
	mg.Deps(Install.Tools, Build.Compile)
	fmt.Println("installing...")
	return os.Rename("build/pleiades", "/usr/bin/pleiades")
}

// install the binary to a homebrew location - only for use with homebrew
func (Install) Homebrew(path string) error {
	fmt.Println("installing to homebrew...")
	return compileWithPath(path, nil)
}

// fetch the go dependencies
func (Install) Godeps() error {
	fmt.Println("installing go dependencies")

	mg.Deps(func() error {
		if err := sh.RunWithV(nil, "go", "install", "github.com/boumenot/gocover-cobertura@latest"); err != nil {
			return err
		}
		return nil
	})

	mg.Deps(func() error {
		if err := sh.RunWithV(nil, "go", "install", "gotest.tools/gotestsum@latest"); err != nil {
			return err
		}
		return nil
	})

	err := sh.RunWithV(nil, "go", "get", "-v", "./...")
	if err != nil {
		return err
	}

	return nil
}

// fetch the nodejs dependencies
func (Install) Node() error {
	fmt.Println("installing nodejs dependencies")
	return sh.RunWithV(nil, "npm", "install")
}

// install necessary tools and dependencies to develop pleiades
func (Install) Tools() error {
	fmt.Println("installing tools...")

	// each of these should be their own dep :shrug:
	mg.Deps(func() error {
		for idx := range homebrewTargets {
			if err := sh.RunWithV(nil, "brew", "install", homebrewTargets[idx]); err != nil {
				return err
			}
		}
		return nil
	})

	mg.Deps(func() error {
		if err := sh.RunWithV(nil, "go", "install", "github.com/spf13/cobra-cli@latest"); err != nil {
			return err
		}
		return nil
	})

	return nil
}
