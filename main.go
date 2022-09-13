/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package main

import (
	"github.com/mxplusb/pleiades/cmd"
	"github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
)

func init() {
	encoding.RegisterCodec(grpc.Codec{})
}

func main() {
	viper.SetConfigName("pleiades") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/pleiades/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.pleiades")  // call multiple times to add many search paths
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	cmd.Execute()
}
