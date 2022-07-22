//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

// prep your workspace
func Setup() {
	fmt.Println("setting up workspace")
	mg.SerialDeps(Install.Tools, CI.Setup, Install.Homebrew, Install.Godeps)
}
