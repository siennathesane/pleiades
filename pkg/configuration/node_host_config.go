/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

type NodeHostConfig struct {
	DevMode                       bool         `flag:"dev d" desc:"enable dev mode?"`
	DeploymentID                  uint64       `desc:"deployment id of the node host"`
	WALDir                        string       `desc:"wal log directory, defaults to node-host-dir location if unspecified"`
	NodeHostDir                   string       `desc:"data directory"`
	RTTMillisecond                uint64       `flag:"rtt" desc:"average round trip time, plus processing, in milliseconds to other hosts in the cluster"`
	RaftAddress                   string       `desc:"address of the node host as seen by other hosts"`
	AddressByNodeHostID           bool         `desc:"addressByNodeHostID indicates that NodeHost instances should be addressed by their NodeHostID values"`
	ListenAddress                 string       `desc:"address to listen on"`
	MutualTLS                     bool         `desc:"require mutual tls?"`
	CAFile                        string       `desc:"location of the certificate authority file"`
	CertFile                      string       `desc:"location of the tls cert file"`
	KeyFile                       string       `desc:"location of the tls key file"`
	EnableMetrics                 bool         `desc:"enable metrics?"`
	MaxSendQueueSize              uint64       `flag:"max-send-queue" desc:"maximum size in bytes of each send queue"`
	MaxReceiveQueueSize           uint64       `flag:"max-receive-queue" desc:"maximum size in bytes of each receive queue"`
	MaxSnapshotSendBytesPerSecond uint64       `flag:"max-send-snapshot" desc:"how much snapshot data can be sent every second for all Raft clusters managed by the node host instance. 0 means there is no limit"`
	MaxSnapshotRecvBytesPerSecond uint64       `flag:"max-receive-snapshot" desc:"how much snapshot data can be received every second for all Raft clusters managed by the node host instance. 0 means there is no limit"`
	NotifyCommit                  bool         `desc:"notify consumers of when the data is committed?"`
	Gossip                        GossipConfig
}

type NodeHostGossipConfig struct {
	BindAddress      string   `desc:"address for the gossip service to bind to and listen on"`
	AdvertiseAddress string   `desc:"address to advertise to other NodeHost instances used for NAT traversal"`
	Seed             []string `desc:"list of advertise addresses of remote node host instances"`
}
