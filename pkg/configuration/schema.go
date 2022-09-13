/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

type Configuration struct {
	Datastore *Storage   `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	Host      HostConfig `json:"host,omitempty" yaml:"host,omitempty"`
}

type Storage struct {
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`
}

type HostConfig struct {
	DeploymentId  uint64 `desc:"deployment id of the node host"`
	LogDir  string `desc:"wal log directory"`
	DataDir       string `desc:"data directory"`
	Rtt           uint64 `flag:"rtt" desc:"average round trip time, plus processing, in milliseconds to other hosts in the data centre"`
	ListenAddress string `desc:"address to listen on"`
	MutualTLS     bool   `desc:"require mutual tls?"`
	CaFile        string `desc:"location of the certificate authority file"`
	CertFile      string `desc:"location of the tls cert file"`
	KeyFile       string `desc:"location of the tls key file"`
	NotifyCommit  bool   `desc:"notify consumers of when the data is committed?"`
}

func DefaultConfiguration() *Configuration {
	return &Configuration{
		Datastore: &Storage{
			BasePath: "/var/pleiades/shards",
		},
		Host: HostConfig{
			DeploymentId:  1,
			LogDir:  "/var/pleiades/logs",
			DataDir:       "/var/pleiades/shards",
			Rtt:           10,
			ListenAddress: "0.0.0.0:5000",
			MutualTLS:     true,
			CaFile:        "/etc/pleiades/ca.pem",
			CertFile:      "/etc/pleiades/cert.pem",
			KeyFile:       "/etc/pleiades/key.pem",
			NotifyCommit:  true,
		},
	}
}
