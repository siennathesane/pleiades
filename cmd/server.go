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
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/mxplusb/cliflags/gen/gpflag"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var srvCmd = &cobra.Command{
	Use:   "server",
	Short: "run an instance of pleiades",
	Long: `Run an instance of the Pleiades Platform Operating System.

DeploymentID is used to determine whether two NodeHost instances belong to
the same deployment and thus allowed to communicate with each other. This
helps to prvent accidentially misconfigured NodeHost instances to cause
data corruption errors by sending out of context messages to unrelated
Raft nodes.

For a particular dragonboat based application, you can set DeploymentID
to the same uint64 value on all production NodeHost instances, then use
different DeploymentID values on your staging and dev environment. It is
also recommended to use different DeploymentID values for different
dragonboat based applications.

When not set, the default value 0 will be used as the deployment Id and
thus allowing all NodeHost instances with deployment Id 0 to communicate
with each other.
	   
WALDir is the directory used for storing the WAL of Raft entries. It is
recommended to use low latency storage such as NVME SSD with power loss
protection to store such WAL data. Leave WALDir to have zero value will
have everything stored in NodeHostDir.
	
NodeHostDir is where everything else is stored.
	
RTTMillisecond defines the average Rround Trip Time (RTT) in milliseconds
between two NodeHost instances. Such a RTT interval is internally used as
a logical clock tick, Raft heartbeat and election intervals are both
defined in term of how many such logical clock ticks (RTT intervals).
Note that RTTMillisecond is the combined delays between two NodeHost
instances including all delays caused by network transmission, delays
caused by NodeHost queuing and processing. As an example, when fully
loaded, the average Rround Trip Time between two of our NodeHost instances
used for benchmarking purposes is up to 500 microseconds when the ping time
between them is 100 microseconds. Set RTTMillisecond to 1 when it is less
than 1 million in your environment.
	
RaftAddress is a DNS name:port or IP:port address used by the transport
module for exchanging Raft messages, snapshots and metadata between
NodeHost instances. It should be set to the public address that can be
accessed from remote NodeHost instances.

When the NodeHostConfig.ListenAddress field is empty, NodeHost listens on
RaftAddress for incoming Raft messages. When hostname or domain name is
used, it will be resolved to IPv4 addresses first and Dragonboat listens
to all resolved IPv4 addresses.

By default, the RaftAddress value is not allowed to change between NodeHost
restarts. AddressByNodeHostID should be set to true when the RaftAddress
value might change after restart.
	
AddressByNodeHostID indicates that NodeHost instances should be addressed
by their NodeHostID values. This feature is usually used when only dynamic
addresses are available. When enabled, NodeHostID values should be used
as the target parameter when calling NodeHost's StartCluster,
AddNode, AddObserver and AddWitness methods.

Enabling AddressByNodeHostID also enables the internal gossip service,
NodeHostConfig.Gossip must be configured to control the behaviors of the
gossip service.

Note that once enabled, the AddressByNodeHostID setting can not be later
disabled after restarts.

Please see the godocs of the NodeHostConfig.Gossip field for a detailed
example on how AddressByNodeHostID and gossip works.
	
ListenAddress is an optional field in the hostname:port or IP:port address
form used by the transport module to listen on for Raft message and
snapshots. When the ListenAddress field is not set, The transport module
listens on RaftAddress. If 0.0.0.0 is specified as the IP of the
ListenAddress, Dragonboat listens to the specified port on all network
interfaces. When hostname or domain name is used, it will be resolved to
IPv4 addresses first and Dragonboat listens to all resolved IPv4 addresses.
	
MutualTLS defines whether to use mutual TLS for authenticating servers
and clients. Insecure communication is used when MutualTLS is set to
False.
See https://github.com/lni/dragonboat/wiki/TLS-in-Dragonboat for more
details on how to use Mutual TLS.
	
CAFile is the path of the CA certificate file. This field is ignored when
MutualTLS is false.
	
CertFile is the path of the node certificate file. This field is ignored
when MutualTLS is false.
	
KeyFile is the path of the node key file. This field is ignored when
MutualTLS is false.

EnableMetrics determines whether health metrics in Prometheus format should
be enabled.

MaxSendQueueSize is the maximum size in bytes of each send queue.
Once the maximum size is reached, further replication messages will be
dropped to restrict memory usage. When set to 0, it means the send queue
size is unlimited.
	
MaxReceiveQueueSize is the maximum size in bytes of each receive queue.
Once the maximum size is reached, further replication messages will be
dropped to restrict memory usage. When set to 0, it means the queue size
is unlimited.
	
MaxSnapshotSendBytesPerSecond defines how much snapshot data can be sent
every second for all Raft clusters managed by the NodeHost instance.
The default value 0 means there is no limit set for snapshot streaming.
	
MaxSnapshotRecvBytesPerSecond defines how much snapshot data can be
received each second for all Raft clusters managed by the NodeHost instance.
The default value 0 means there is no limit for receiving snapshot data.
	
NotifyCommit specifies whether clients should be notified when their
regular proposals and config change requests are committed. By default,
commits are not notified, clients are only notified when their proposals
are both committed and applied.
	
Gossip contains configurations for the gossip service. When the
AddressByNodeHostID field is set to true, each NodeHost instance will use
an internal gossip service to exchange knowledges of known NodeHost
instances including their RaftAddress and NodeHostID values. This Gossip
field contains configurations that controls how the gossip service works.

As an detailed example on how to use the gossip service in the situation
where all available machines have dynamically assigned IPs on reboot -

Consider that there are three NodeHost instances on three machines, each
of them has a dynamically assigned IP address which will change on reboot.
NodeHostConfig.RaftAddress should be set to the current address that can be
reached by remote NodeHost instance. In this example, we will assume they
are

10.0.0.100:24000
10.0.0.200:24000
10.0.0.300:24000

To use these machines, first enable the NodeHostConfig.AddressByNodeHostID
field and start the NodeHost instances. The NodeHostID value of each
NodeHost instance can be obtained by calling NodeHost.Id(). Let's say they
are

"nhid-xxxxx",
"nhid-yyyyy",
"nhid-zzzzz".

All these NodeHostID are fixed, they will never change after reboots.

When starting Raft nodes or requesting new nodes to be added, use the above
mentioned NodeHostID values as the target parameters (which are of the
Target type). Let's say we want to start a Raft Node as a part of a three
replicas Raft cluster, the initialMembers parameter of the StartCluster
method can be set to

initialMembers := map[uint64]Target {
	 1: "nhid-xxxxx",
  2: "nhid-yyyyy",
  3: "nhid-zzzzz",
}

This indicates that node 1 of the cluster will be running on the NodeHost
instance identified by the NodeHostID value "nhid-xxxxx", node 2 of the
same cluster will be running on the NodeHost instance identified by the
NodeHostID value of "nhid-yyyyy" and so on.

The internal gossip service exchanges NodeHost details, including their
NodeHostID and RaftAddress values, with all other known NodeHost instances.
Thanks to the nature of gossip, it will eventually allow each NodeHost
instance to be aware of the current details of all NodeHost instances.
As a result, let's say when Raft node 1 wants to send a Raft message to
node 2, it first figures out that node 2 is running on the NodeHost
identified by the NodeHostID value "nhid-yyyyy", RaftAddress information
from the gossip service further shows that "nhid-yyyyy" maps to a machine
currently reachable at 10.0.0.200:24000. Raft messages can thus be
delivered.

The Gossip field here is used to configure how the gossip service works.
In this example, let's say we choose to use the following configurations
for those three NodeHost instaces.

GossipConfig {
  BindAddress: "10.0.0.100:24001",
  Seed: []string{10.0.0.200:24001},
}

GossipConfig {
  BindAddress: "10.0.0.200:24001",
  Seed: []string{10.0.0.300:24001},
}

GossipConfig {
  BindAddress: "10.0.0.300:24001",
  Seed: []string{10.0.0.100:24001},
}

For those three machines, the gossip component listens on
"10.0.0.100:24001", "10.0.0.200:24001" and "10.0.0.300:24001" respectively
for incoming gossip messages. The Seed field is a list of known gossip end
points the local gossip service will try to talk to. The Seed field doesn't
need to include all gossip end points, a few well connected nodes in the
gossip network is enough.
`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var (
	cfg *configuration.NodeHostConfig = &configuration.NodeHostConfig{
		DevMode:                       true,
		DeploymentID:                  1,
		WALDir:                        "/var/pleiades/wal",
		NodeHostDir:                   "/var/pleiades/data",
		RTTMillisecond:                200,
		RaftAddress:                   "0.0.0.0:5001",
		AddressByNodeHostID:           false,
		ListenAddress:                 "0.0.0.0:5002",
		MutualTLS:                     false,
		CAFile:                        "/etc/pleiades/ca.pem",
		CertFile:                      "/etc/pleiades/cert.pem",
		KeyFile:                       "/etc/pleiades/cert.key",
		EnableMetrics:                 true,
		MaxSendQueueSize:              0,
		MaxReceiveQueueSize:           0,
		MaxSnapshotSendBytesPerSecond: 0,
		MaxSnapshotRecvBytesPerSecond: 0,
		NotifyCommit:                  true,
		Gossip: configuration.GossipConfig{
			AdvertiseAddress: "",
			BindAddress:      "",
			Seed:             []string{},
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// if we're on a mac, set different paths for the default config
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		cfg.WALDir = filepath.Join(dir, "Library", "pleiades", "wal")
		cfg.NodeHostDir = filepath.Join(dir, "Library", "pleiades", "data")
		cfg.CAFile = ""
		cfg.CertFile = ""
		cfg.KeyFile = ""
	}

	if err := gpflag.ParseTo(cfg, serverCmd.Flags()); err != nil {
		log.Logger.Err(err).Msg("cannot properly parse command strings")
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("error: %s", err)
		} else {
			fmt.Printf("error: %s", err)
		}
	}
}

func startServer() {
	//ctx, cancel := context.WithCancel(context.Background())

	logger := configuration.NewRootLogger()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	done := make(chan bool, 1)
	go func(sigs chan os.Signal, done chan bool) {
		<-sigs
		done <- true
	}(sigs, done)

	logger.Info().Msg("goodbye")
}
