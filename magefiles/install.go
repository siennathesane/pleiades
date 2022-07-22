//go:build mage
package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

var (
	homebrewTargets = []string{
		"capnp",
		"vault",
		"yq",
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
	return compileWithPath(path)
}

// fetch the go dependencies
func (Install) Godeps() error {
	fmt.Println("installing go dependencies")
	err := sh.RunWithV(nil, "go", "get", "-v", "./...")
	if err != nil {
		return err
	}

	return verifyVendor()
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

	mg.Deps(func() error {
		if err := sh.RunWithV(nil, "go", "install", "github.com/nomad-software/vend@latest"); err != nil {
			return err
		}
		return nil
	})

	return nil
}
