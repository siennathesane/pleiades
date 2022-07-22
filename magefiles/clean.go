//go:build mage
package main

import (
	"fmt"
	"os"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Clean mg.Namespace

// clear the local build directory
func (Clean) Build() error {
	fmt.Println("cleaning build cache")
	if err := sh.Rm("build"); err != nil {
		return err
	}

	if err := sh.RunWithV(nil, "go", "clean"); err != nil {
		return err
	}

	return nil
}

// clear the package cache
func (Clean) Cache() error {
	fmt.Println("cleaning mod cache...")
	return sh.RunWithV(nil, "go", "clean", "-modcache")
}

// clear all tools and dependencies
func (Clean) All() error {
	fmt.Println("removing build directory...")
	err := os.RemoveAll("build")
	if err != nil {
		return err
	}

	fmt.Println("cleaning mod cache...")
	if err := sh.RunWithV(nil, "go", "clean", "-modcache"); err != nil {
		return err
	}

	fmt.Println("removing homebrew tools")
	for idx := range homebrewTargets {
		if err := sh.RunWithV(nil, "brew", "remove", homebrewTargets[idx]); err != nil {
			return err
		}
	}
	return nil
}
