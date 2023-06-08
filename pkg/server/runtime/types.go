/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

type StateMachineType uint64
type ResultCode int

const (
	TimeoutResultCode ResultCode = iota
	CompletedResultCode
	TerminatedResultCode
	RejectedResultCode
	DroppedResultCode
	AbortedResultCode
	CommittedResultCode
)

type HostInfo struct {
	HostId          string
	RaftAddress     string
	Gossip          GossipInfo
	ClusterInfoList []ClusterInfo
	LogInfo         []NodeInfo
}

type GossipInfo struct {
	Enabled             bool
	AdvertiseAddress    string
	NumOfLiveKnownHosts int
}

type ClusterInfo struct {
	ShardId           uint64
	ReplicaId         uint64
	Nodes             map[uint64]string
	ConfigChangeIndex uint64
	IsLeader          bool
	IsObserver        bool
	IsWitness         bool
	Pending           bool
}

type NodeInfo struct {
	ShardId   uint64
	ReplicaId uint64
}

type HostConfig struct {
	DeploymentID        uint64
	WALDir              string
	NodeHostDir         string
	RTTMillisecond      uint64
	RaftAddress         string
	AddressByNodeHostID bool
	ListenAddress       string
	MutualTLS           bool
	CAFile              string
	CertFile            string
	KeyFile             string
	EnableMetrics       bool
	NotifyCommit        bool
}

type HostInfoOption struct {
	SkipLogInfo bool
}

type SnapshotOption struct {
	CompactionOverhead         uint64
	ExportPath                 string
	Exported                   bool
	OverrideCompactionOverhead bool
}

type OperationResult struct {
	Status ResultCode
	Index  uint64
	Data   []byte
}

// MembershipEntry is the struct used to describe Raft cluster membership query
// results.
type MembershipEntry struct {
	// ConfigChangeId is the Raft entry index of the last applied membership
	// change entry.
	ConfigChangeId uint64
	// Nodes is a map of ReplicaId values to NodeHost Raft addresses for all regular
	// Raft nodes.
	Replicas map[uint64]string
	// Observers is a map of ReplicaId values to NodeHost Raft addresses for all
	// observers in the Raft cluster.
	Observers map[uint64]string
	// Witnesses is a map of ReplicaId values to NodeHost Raft addresses for all
	// witnesses in the Raft cluster.
	Witnesses map[uint64]string
	// Removed is a set of ReplicaId values that have been removed from the Raft
	// cluster. They are not allowed to be added back to the cluster.
	Removed map[uint64]struct{}
}
