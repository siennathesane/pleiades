//go:build mage
package main

import (
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

// run all tests
func (Test) Cleanroom() error {
	fmt.Println("running all tests in cleanroom environment")
	mg.SerialDeps(Clean.Build, Clean.Cache, Install.Godeps)
	return sh.RunWithV(nil, "go", "test", "-v", "./...")
}

// run all tests
func (Test) All() error {
	fmt.Println("running all tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./...")
}

// run blaze server tests
func (Test) Blaze() error {
	fmt.Println("running blaze tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./pkg/blaze/...")
}

// run config tests
func (Test) Config() error {
	fmt.Println("running config tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./pkg/conf/...")
}

// run fsm tests
func (Test) FSM() error {
	fmt.Println("running fsm tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./pkg/fsm/...")
}

// run routing tests
func (Test) Routing() error {
	fmt.Println("running routing tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./pkg/routing/...")
}

// run service tests
func (Test) Services() error {
	fmt.Println("running service tests")
	return sh.RunWithV(nil, "go", "test", "-v", "./pkg/services/...")
}