/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

// raftCmd represents the raft command
var raftCmd = &cobra.Command{
	Use:   "raft",
	Short: "operations on the raft subsystem",
	Long: `various operations used to control various aspects of the raft subsystems`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("raft called")
	},
}

func init() {
	rootCmd.AddCommand(raftCmd)

	raftCmd.PersistentFlags().String("host", "http://localhost:8080", "target host for a pleiades cluster")
	config.BindPFlag("server.client.grpcAddr", raftCmd.PersistentFlags().Lookup("host"))
}

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}