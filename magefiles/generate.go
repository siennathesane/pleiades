//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
	"path/filepath"
)

const (
	vendoredStdPath = "-Ivendor/capnproto.org/go/capnp/v3/std"
)

type Gen mg.Namespace

// setup the generator tools and environment
func (Gen) Setup() {
	mg.Deps(func() error {
		fmt.Println("installing capn' proto compiler")
		return sh.RunWithV(nil, "brew", "install", "capnp")
	})

	mg.Deps(func() error {
		fmt.Println("installing capn' proto go compiler plugin")
		return sh.RunWithV(nil, "go", "install", "capnproto.org/go/capnp/v3/capnpc-go@latest")
	})

	mg.Deps(func() error {
		fmt.Println("installing cap'n proto golang compiler cli")
		return sh.RunWithV(map[string]string{}, "go", "get", "-v", "capnproto.org/go/capnp/v3")
	})

	mg.Deps(func() error {
		fmt.Println("installing cap'n proto golang compiler cli")
		return sh.RunWithV(map[string]string{}, "go", "mod", "vendor")
	})
}

// generate all schemas
func (Gen) All() error {
	if err := verifyVendor(); err != nil {
		return err
	}

	mg.SerialDeps(Gen.Host, Gen.Database)
	return nil
}

// compiles the database schemas and generates the go code
func (Gen) Database() error {
	if err := verifyVendor(); err != nil {
		return err
	}

	fmt.Println("generating database protocols")
	files, err := filepath.Glob("protocols/v1/database/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", vendoredStdPath, "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}

// compiles the host schemas and generates the go code
func (Gen) Host() error {
	if err := verifyVendor(); err != nil {
		return err
	}

	fmt.Println("generating host protocols")
	files, err := filepath.Glob("protocols/v1/host/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", vendoredStdPath, "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}
