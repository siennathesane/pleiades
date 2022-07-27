/*
 * Copyright (c) 2022 Sienna Lloyd
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
	"path/filepath"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	vendoredStdPath = "-Ivendor/capnproto.org/go/capnp/v3/std"
	nodeJsBinPath = "node_modules/.bin"
)

var (
	goProtoFlags = []string{
		"-I",
		".",
		"--plugin",
		fmt.Sprintf("protoc-gen-ts=%s/protoc-gen-ts", nodeJsBinPath),
		"--plugin",
		fmt.Sprintf("protoc-gen-go=%s/protoc-gen-go", binDir),
		"--plugin",
		fmt.Sprintf("protoc-gen-go-vtproto=%s/protoc-gen-go-vtproto", binDir),
		"--plugin",
		fmt.Sprintf("protoc-gen-go-starpc=%s/protoc-gen-go-starpc", binDir),
		"--js_out=import_style=commonjs,binary:.",
		"--ts_out=.",
		"--ts_opt=esModuleInterop=true",
		"--ts_opt=forceLong=long",
		"--ts_opt=oneof=unions",
		"--ts_opt=outputServices=default,outputServices=generic-definitions",
		"--ts_opt=useDate=true",
		"--ts_opt=useAsyncIterable=true",
		"--ts_opt=fileSuffix=.pb",
		"--ts_opt=iportSuffix=.js",
		"--go_opt=paths=source_relative",
		"--go_out=.",
		"--go-starpc_out=.",
		"--go-starpc_opt=paths=source_relative",
		"--go-vtproto_out=.",
		"--go-vtproto_opt=features=marshal+unmarshal+size+equal+pool",
		"--go-vtproto_opt=paths=source_relative",
	}
)

type Gen mg.Namespace

// setup the generator tools and environment
func (Gen) Setup() {
	defer func() {
		verifyVendor()
	}()

	mg.Deps(Clean.Vendor)

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
		fmt.Println("installing vitness protobufs")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto"); err != nil {
			return err
		}

		return sh.RunWithV(nil, "go",
			"build",
			"-v",
			"-o",
			fmt.Sprintf("%s/protoc-gen-go-vtproto", binDir),
			"github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto")
	})

	mg.Deps(func() error {
		fmt.Println("installing protobuf golang generator")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"google.golang.org/protobuf/cmd/protoc-gen-go"); err != nil {
			return err
		}

		return sh.RunWithV(nil, "go",
			"build",
			"-v",
			"-o",
			fmt.Sprintf("%s/protoc-gen-go", binDir),
			"google.golang.org/protobuf/cmd/protoc-gen-go")
	})

	mg.Deps(func() error {
		fmt.Println("installing the starpc protobuf generator")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"github.com/aperturerobotics/starpc/cmd/protoc-gen-go-starpc"); err != nil {
			return err
		}

		return sh.RunWithV(nil, "go",
			"build",
			"-v",
			"-o",
			fmt.Sprintf("%s/protoc-gen-go-starpc", binDir),
			"github.com/aperturerobotics/starpc/cmd/protoc-gen-go-starpc")
	})

	mg.Deps(func() error {
		fmt.Println("installing node protobuf compiler")
		return sh.RunWithV(nil, "npm", "install")
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

	pbFiles, err := filepath.Glob("api/v1/database/*.proto")
	if err != nil {
		return err
	}
	localProtoFlags := append(goProtoFlags, pbFiles...)
	if err := sh.RunWithV(nil, "protoc", localProtoFlags...); err != nil {
		return err
	}

	capnpFiles, err := filepath.Glob("protocols/v1/database/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", vendoredStdPath, "-ogo:pkg"}
	args = append(args, capnpFiles...)
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
