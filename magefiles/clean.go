//go:build mage
package main

import (
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Clean mg.Namespace

// clear the local build directory
func (Clean) Build() error {
	fmt.Println("cleaning build directory")
	if err := sh.Rm(buildDir); err != nil {
		return err
	}

	if err := sh.RunWithV(nil, "go", "clean"); err != nil {
		return err
	}

	return nil
}

func (Clean) Vendor() error {
	fmt.Println("cleaning vendor cache")
	return sh.Rm(vendorDir)
}

// clear the package cache
func (Clean) Cache() error {
	fmt.Println("cleaning mod cache...")
	return sh.RunWithV(nil, "go", "clean", "-modcache")
}

// clean the bin directory
func (Clean) Bindir() error {
	fmt.Println("removing bin directory")
	return sh.Rm(binDir)
}

// remove the homebrew tools
func (Clean) Homebrew() error {
	fmt.Println("removing homebrew tools")
	for idx := range homebrewTargets {
		if err := sh.RunWithV(nil, "brew", "remove", homebrewTargets[idx]); err != nil {
			return err
		}
	}
	return nil
}

// clear all tools and dependencies
func (Clean) All() {
	mg.Deps(Clean.Build, Clean.Cache, Clean.Bindir, Clean.Homebrew, Clean.Vendor)
}
