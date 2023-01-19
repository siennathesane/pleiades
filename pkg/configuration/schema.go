/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

import (
	"net"
)

const (
	DefaultBaseDataPath   string = "/var/pleiades"
	DefaultBaseConfigPath string = "/etc/pleiades"

	DefaultLogDir = "logs"
	DefaultDataDir = "data"
)

type Configuration struct {
	ConfigFilePath string        `flag:"config" default:"/etc/pleiades/config.yaml" usage:"config file location" mapstructure:"configPath"`
	Debug          bool          `flag:"debug" default:"false" yaml:"debug" json:"debug" mapstructure:"debug"`
	Trace          bool          `flag:"trace" default:"false" yaml:"trace" json:"trace" mapstructure:"trace"`
	Server         *ServerConfig `json:"server" yaml:"server" mapstructure:"server"`
	Client         *ClientConfig `json:"client" yaml:"client" mapstructure:"client"`
}

type ClientConfig struct {
	GrpcAddr string `flag:"host" default:"http://localhost:8080" usage:"address to call" json:"grpcAddr" yaml:"grpcAddr" mapstructure:"grpcAddr"`
}

type ServerConfig struct {
	Datastore *Datastore `json:"datastore,omitempty" yaml:"datastore,omitempty" mapstructure:"datastore"`
	Host      *Host      `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host"`
	Reset     bool       // internal flag to reset the dev server
}

type Datastore struct {
	BasePath string `flag:"base-path" default:"/var/pleiades" usage:"set the default base directory" json:"basePath,omitempty" yaml:"basePath,omitempty" mapstructure:"basePath"`
	LogDir   string `flag:"log-dir" default:"logs" usage:"folder path for the logs, relative to the base path" json:"logDir" yaml:"logDir" mapstructure:"logDir"`
	DataDir  string `flag:"data-dir" default:"data" usage:"folder path for the data, relative to the base path" json:"dataDir" yaml:"dataDir" mapstructure:"dataDir"`
}

type Host struct {
	CaFile                  string `flag:"ca-file" default:"/etc/pleiades/ca.pem" usage:"location of the certificate authority file" json:"caFile" yaml:"caFile" mapstructure:"caFile"`
	CertFile                string `flag:"cert-file" default:"/etc/pleiades/cert.pem" usage:"location of the tls cert file" json:"certFile" yaml:"certFile" mapstructure:"certFile"`
	DeploymentId            uint64 `flag:"deployment-id" default:"1" usage:"deployment id of this host" json:"deploymentId" yaml:"deploymentId" mapstructure:"deploymentId"`
	KeyFile                 string `flag:"key-file" default:"/etc/pleiades/key.pem" usage:"location of the tls key file" json:"keyFile" yaml:"keyFile" mapstructure:"keyFile"`
	ListenAddress           *net.IP `flag:"listen-address" default:"0.0.0.0" usage:"address to listen on" json:"listenAddress" yaml:"listenAddress" mapstructure:"listenAddress"`
	GrpcListenPort          uint   `flag:"grpc-port" default:"8080" usage:"address to listen on" json:"grpcListenPort" yaml:"grpcListenPort" mapstructure:"grpcListenPort"`
	RaftListenPort          uint   `flag:"raft-port" default:"8081" usage:"address to listen on" json:"raftListenPort" yaml:"raftListenPort" mapstructure:"raftListenPort"`
	MutualTLS               bool   `flag:"mtls" default:"false" usage:"require mutual tls?" mapstructure:"mtls"`
	NotifyCommit            bool   `flag:"notify-commit" default:"false" usage:"notify consumers of when the data is committed?" json:"notifyCommit" yaml:"notifyCommit" mapstructure:"notifyCommit"`
	Rtt                     uint64 `flag:"rtt" default:"10" usage:"average round trip time, plus processing, in milliseconds to other hosts in the data centre" json:"rtt" yaml:"rtt" mapstructure:"rtt"`
	ServiceDiscoveryAddress string `flag:"sd-address" default:"" json:"serviceDiscoveryAddress" yaml:"serviceDiscoveryAddress" mapstructure:"serviceDiscoveryAddress"`
}
