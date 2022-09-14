/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

const (
	DefaultBaseDataPath   string = "/var/pleiades"
	DefaultBaseConfigPath string = "/etc/pleiades"
)

type Configuration struct {
	ConfigFilePath string       `flag:"config" default:"/etc/pleiades/config.yaml" usage:"config file location"`
	Server         *ServerConfig `json:"server" yaml:"server" mapstructure:"server"`
	Client         *ClientConfig `json:"client" yaml:"client" mapstructure:"client"`
}

type ClientConfig struct{}

type ServerConfig struct {
	Datastore *Datastore `json:"datastore,omitempty" yaml:"datastore,omitempty" mapstructure:"datastore"`
	Host      *Host      `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host"`
}

type Datastore struct {
	BasePath string `flag:"base-path" default:"/var/pleiades" usage:"set the default base directory" json:"basePath,omitempty" yaml:"basePath,omitempty" mapstructure:""`
	LogDir   string `flag:"log-dir" default:"logs" usage:"folder path for the logs, relative to the base path" mapstructure:"" json:"logDir" yaml:"logDir"`
	DataDir  string `flag:"data-dir" default:"data" usage:"folder path for the data, relative to the base path" mapstructure:"" json:"dataDir" yaml:"dataDir"`
}

type Host struct {
	CaFile        string `flag:"ca-file" default:"/etc/pleiades/ca.pem" usage:"location of the certificate authority file" json:"caFile" yaml:"caFile" mapstructure:""`
	CertFile      string `flag:"cert-file" default:"/etc/pleiades/cert.pem" usage:"location of the tls cert file" json:"certFile" yaml:"certFile" mapstructure:""`
	DeploymentId  uint64 `flag:"deployment-id" default:"1" usage:"deployment id of this host" json:"deploymentId" yaml:"deploymentId" mapstructure:""`
	KeyFile       string `flag:"key-file" default:"/etc/pleiades/key.pem" usage:"location of the tls key file" json:"keyFile" yaml:"keyFile" mapstructure:""`
	ListenAddress string `flag:"listen-address" default:"0.0.0.0:5001" usage:"address to listen on" json:"listenAddress" yaml:"listenAddress" mapstructure:""`
	GrpcListenAddress string `flag:"grpc-listen-address" default:"0.0.0.0:5000" usage:"address to listen on" json:"grpcListenAddress" yaml:"grpcListenAddress" mapstructure:""`
	MutualTLS     bool   `flag:"mtls" default:"false" usage:"require mutual tls?" mapstructure:""`
	NotifyCommit  bool   `flag:"notify-commit" default:"false" usage:"notify consumers of when the data is committed?" json:"notifyCommit" yaml:"notifyCommit" mapstructure:""`
	Rtt           uint64 `flag:"rtt" default:"10" usage:"average round trip time, plus processing, in milliseconds to other hosts in the data centre" json:"rtt" yaml:"rtt" mapstructure:""`
}
