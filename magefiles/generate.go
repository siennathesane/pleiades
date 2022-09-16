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
		"--js_out=import_style=commonjs,binary:.",
		"--ts_out=.",
		"--ts_opt=esModuleInterop=true",
		"--ts_opt=forceLong=long",
		"--ts_opt=oneof=unions",
		"--ts_opt=outputServices=default,outputServices=generic-definitions",
		"--ts_opt=useDate=true",
		"--ts_opt=useAsyncIterable=true",
		"--ts_opt=fileSuffix=.pb",
		"--ts_opt=importSuffix=.js",
		"--ts_opt=useDate=true",
		"--go_opt=paths=source_relative",
		"--go_out=.",
		"--go-vtproto_out=.",
		"--go-vtproto_opt=features=marshal+unmarshal+size+equal+pool",
		"--go-vtproto_opt=paths=source_relative",
	}

	grpcFlags = []string{
		"-I",
		".",
		"--plugin",
		fmt.Sprintf("protoc-gen-go=%s/protoc-gen-go", binDir),
		"--plugin",
		fmt.Sprintf("protoc-gen-go=%s/protoc-gen-ts", nodeJsBinPath),
		"--go_opt=paths=source_relative",
		"--go_out=.",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
	}

	fuckingNodeJsFlags = []string{
		"-I",
		".",
		"--js_out=import_style=commonjs,binary:.",
		"--grpc_out=grpc_js:.",
	}
)

type Gen mg.Namespace

// setup the generator tools and environment
func (Gen) Setup() {
	mg.Deps(Clean.Vendor)

	mg.Deps(func() error {
		fmt.Println("installing vitness protobufs")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@0ae748f"); err != nil {
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
		fmt.Println("installing golang grpc generator")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc"); err != nil {
			return err
		}

		return sh.RunWithV(nil, "go",
			"build",
			"-v",
			"-o",
			fmt.Sprintf("%s/protoc-gen-go-grpc", binDir),
			"google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	})

	mg.Deps(func() error {
		fmt.Println("installing node protobuf compiler")
		return sh.RunWithV(nil, "npm", "install")
	})
}

// generate all schemas
func (Gen) All() error {

	mg.SerialDeps(Gen.Raft, Gen.DB, Gen.Server)
	return nil
}

// compiles the database schemas and generates the go code
func (Gen) DB() error {

	fmt.Println("generating database protocols")

	pbFiles, err := filepath.Glob("api/v1/database/*.proto")
	if err != nil {
		return err
	}
	localProtoFlags := append(goProtoFlags, pbFiles...)
	if err := sh.RunWithV(nil, "protoc", localProtoFlags...); err != nil {
		return err
	}

	errorPbFiles, err := filepath.Glob("api/v1/errors/*.proto")
	if err != nil {
		return err
	}
	localProtoFlags = append(goProtoFlags, errorPbFiles...)
	if err := sh.RunWithV(nil, "protoc", localProtoFlags...); err != nil {
		return err
	}

	return nil
}

// compiles the raft schemas and generates the go code
func (Gen) Raft() error {

	fmt.Println("generating raft protocols")

	raftPbFiles, err := filepath.Glob("api/v1/raft/*.proto")
	if err != nil {
		return err
	}
	localProtoFlags := append(goProtoFlags, raftPbFiles...)
	if err := sh.RunWithV(nil, "protoc", localProtoFlags...); err != nil {
		return err
	}

	return nil
}

// generates the server grpc code
func (Gen) Server() error {

	fmt.Println("generating server instances")

	serverFiles, err := filepath.Glob("pkg/server/*.proto")
	if err != nil {
		return err
	}
	localGrpcFlags := append(grpcFlags, serverFiles...)
	if err := sh.RunWithV(nil, "protoc", localGrpcFlags...); err != nil {
		return err
	}

	fuckNode := append(fuckingNodeJsFlags, serverFiles...)
	if err := sh.RunWithV(nil, fmt.Sprintf("%s/grpc_tools_node_protoc", nodeJsBinPath), fuckNode...); err != nil {
		return err
	}

	return nil
}
