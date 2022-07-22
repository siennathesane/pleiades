//go:build mage
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Gen mg.Namespace

// generate all schemas
func (Gen) All() {
	mg.SerialDeps(Gen.Host, Gen.Database)
}

// compiles the database schemas and generates the go code
func (Gen) Database() error {
	gopath := os.Getenv("GOPATH")

	fmt.Println("generating database protocols")
	files, err := filepath.Glob("protocols/v1/database/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", fmt.Sprintf("-I%s/src/capnproto.org/go/capnp/std", gopath), "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}

// compiles the host schemas and generates the go code
func (Gen) Host() error {
	gopath := os.Getenv("GOPATH")

	fmt.Println("generating host protocols")
	files, err := filepath.Glob("protocols/v1/host/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", fmt.Sprintf("-I%s/src/capnproto.org/go/capnp/std", gopath), "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}
