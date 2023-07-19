/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

//go:build mage

package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	rootPath = "."
	nodeJsBinPath = "./api/node_modules/.bin"
)

var (
	goProtoFlags = []string{
		"-I",
		rootPath,
		//"--plugin",
		//fmt.Sprintf("protoc-gen-ts=%s/protoc-gen-ts", nodeJsBinPath),
		"--plugin",
		fmt.Sprintf("protoc-gen-go=%s/protoc-gen-go", binDir),
		"--plugin",
		fmt.Sprintf("protoc-gen-go-vtproto=%s/protoc-gen-go-vtproto", binDir),
		"--plugin",
		fmt.Sprintf("protoc-gen-go-connect=%s/protoc-gen-go-connect", binDir),
		//"--js_out=import_style=commonjs,binary:.",
		//"--ts_out=.",
		//"--ts_opt=esModuleInterop=true",
		//"--ts_opt=forceLong=long",
		//"--ts_opt=oneof=unions",
		//"--ts_opt=outputServices=default,outputServices=generic-definitions",
		//"--ts_opt=useDate=true",
		//"--ts_opt=useAsyncIterable=true",
		//"--ts_opt=fileSuffix=.pb",
		//"--ts_opt=importSuffix=.js",
		//"--ts_opt=useDate=true",
		"--go_opt=paths=source_relative",
		"--go_out=.",
		"--go-vtproto_out=.",
		"--go-vtproto_opt=features=marshal+unmarshal+size+equal+pool",
		"--go-vtproto_opt=paths=source_relative",
		"--go-connect_out=.",
		"--go-connect_opt=paths=source_relative",
	}

	grpcFlags = []string{
		"-I",
		".",
		"--plugin",
		fmt.Sprintf("protoc-gen-go=%s/protoc-gen-go", binDir),
		//"--plugin",
		//fmt.Sprintf("protoc-gen-grpc-web=%s/protoc-gen-grpc-web", binDir),
		"--go_opt=paths=source_relative",
		"--go_out=.",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		//"--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:.",
	}

	fuckingNodeJsFlags = []string{
		"-I",
		".",
		"--js_out=import_style=commonjs,binary:.",
		"--grpc_out=grpc_js:.",
	}
)

type Gen mg.Namespace

func (Gen) Print() {
	fmt.Printf("run this from the `pkg` directory after running gen:setup")
	fmt.Printf("protoc %s\n", strings.Join(goProtoFlags, " "))
}

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

	// go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
	mg.Deps(func() error {
		fmt.Println("installing connect golang generator")

		if err := sh.RunWithV(nil, "go",
			"get",
			"-v",
			"github.com/bufbuild/connect-go"); err != nil {
			return err
		}

		return sh.RunWithV(nil, "go",
			"build",
			"-v",
			"-o",
			fmt.Sprintf("%s/protoc-gen-connect-go", binDir),
			"github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go")
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

	//mg.Deps(func() error {
	//	fmt.Println("installing node generator")
	//
	//	if err := sh.RunWithV(nil, "wget", "-O",fmt.Sprintf("%s/protoc-gen-grpc-web", binDir), "https://github.com/grpc/grpc-web/releases/download/1.3.1/protoc-gen-grpc-web-1.3.1-darwin-x86_64"); err != nil {
	//		return err
	//	}
	//
	//	if err := sh.RunWithV(nil, "chmod", "a+x", fmt.Sprintf("%s/protoc-gen-grpc-web", binDir)); err != nil {
	//		return err
	//	}
	//
	//	return nil
	//})

	//mg.Deps(func() error {
	//	fmt.Println("installing node protobuf compiler")
	//
	//	err := os.Chdir("pkg/api")
	//	if err != nil {
	//		return err
	//	}
	//	defer func() {
	//		err := os.Chdir("../..")
	//		if err != nil {
	//			return
	//		}
	//	}()
	//
	//	return sh.RunWithV(nil, "npm", "install")
	//})
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

	errorPbFiles, err := filepath.Glob("pkg/errorspb/*.proto")
	if err != nil {
		return err
	}
	localProtoFlags = append(goProtoFlags, errorPbFiles...)
	if err := sh.RunWithV(nil, "protoc", localProtoFlags...); err != nil {
		return err
	}

	return nil
}

// compiles the raftpb schemas and generates the go code
func (Gen) Raft() error {

	fmt.Println("generating raftpb protocols")

	raftPbFiles, err := filepath.Glob("api/raftpb/*.proto")
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
