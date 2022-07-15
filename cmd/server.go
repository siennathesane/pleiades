

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
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/mxplusb/pleiades/pkg/blaze"
	"github.com/mxplusb/pleiades/pkg/conf"
	"github.com/mxplusb/pleiades/pkg/services/v1/config"
	"github.com/lucas-clemente/quic-go"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run an instance of pleiades",
	Long: `Run an instance of the Pleiades database server.

This command will start the server's listening socket and accept connections,
however it will not start any of the services. This command is intended to primarily be used
for development purpose because it doesn't actually _do_ anything other than open a socket
and listen for connections.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var (
	devMode bool = false
	listenerPort int = 0
	tlsCa string = ""
	tlsCert string = ""
	tlsKey string = ""
	hostname string = ""
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().BoolVar(&devMode, "dev", false, "enable dev mode")
	serverCmd.Flags().IntVar(&listenerPort, "port", 8080, "which port to listen on")
	serverCmd.Flags().StringVar(&tlsCa, "tls-ca-path", "", "location of the certificate authority file")
	serverCmd.Flags().StringVar(&tlsCert, "tls-cert-path", "", "location of the certificate file")
	serverCmd.Flags().StringVar(&tlsCert, "tls-key-path", "", "location of the key file")
	serverCmd.Flags().StringVar(&hostname, "hostname", "", "hostname to use for the server. it must match the hostname in the certificate")
}

func startServer() {
	logger, err := conf.NewLogger()
	if err != nil {
		err = fmt.Errorf("could not instantiate logger: %w", err)
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	l := logger.GetLogger()

	registry, err := config.NewRegistry(logger.GetLogger())
	if err != nil {
		l.Error().Err(err).Msg("cannot instantiate new registry")
		os.Exit(1)
	}

	var tlsConfig *tls.Config
	certPool := x509.NewCertPool()

	if devMode {
		keyPair, err := tls.X509KeyPair([]byte(blaze.DevTlsCert), []byte(blaze.DevTlsKey))
		if err != nil {
			l.Error().Err(err).Msg("cannot instantiate dev tls key pair")
			os.Exit(1)
		}
		certPool.AppendCertsFromPEM([]byte(blaze.DevTlsCa))

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			RootCAs:      certPool,
			ServerName:   "localhost",
			NextProtos: []string{"pleiades"},
		}
	} else {
		ca, err := ioutil.ReadFile(tlsCa)
		if err != nil {
			l.Error().Err(err).Msg("cannot read certificate authority file")
		}

		certPool.AppendCertsFromPEM(ca)

		keyPair, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			l.Error().Err(err).Msg("cannot load key pair")
			os.Exit(1)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			RootCAs:      certPool,
			ServerName:   "localhost",
			NextProtos: []string{"pleiades"},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
	}

	listener, err := quic.ListenAddr("0.0.0.0:"+fmt.Sprintf("%d", listenerPort), tlsConfig, nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	done := make(chan bool, 1)
	go func(done chan bool) {
		<-sigs
		done <- true
	}(done)

	ctx := context.Background()

	server := blaze.NewServer(listener, logger.GetLogger(), registry)
	err = server.Start(ctx)
	if err != nil {
		l.Error().Err(err).Msg("cannot start server")
		os.Exit(1)
	}

	// wait until we get a signal
	<-done

	err = server.Stop(ctx)
	if err != nil {
		l.Error().Err(err).Msg("cannot stop server safely")
		os.Exit(1)
	}
}