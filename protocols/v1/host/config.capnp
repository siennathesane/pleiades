using Go = import "/go.capnp";
@0x9f0a9459eb121bf0;
$Go.package("host");
$Go.import("r3t.io/pkg/protocols/v1/host");

################################################################################
#                  __ _
#  ___ ___  _ __  / _(_) __ _   ___  ___ _ ____   _____ _ __
# / __/ _ \| '_ \| |_| |/ _` | / __|/ _ \ '__\ \ / / _ \ '__|
#| (_| (_) | | | |  _| | (_| | \__ \  __/ |   \ V /  __/ |
# \___\___/|_| |_|_| |_|\__, | |___/\___|_|    \_/ \___|_|
#                       |___/
#
################################################################################

interface ConfigService @0xcd55e3c0a182ac77 {
    getConfig @0 (request :GetConfigurationRequest) -> (response :GetConfigurationResponse);
    putConfig @1 (request :PutConfigurationRequest) -> (response :PutConfigurationResponse);
}

struct GetConfigurationRequest @0xc0e43eb9670b8d20 {
    enum Type {
        all @0;
        raft @1;
    }
    what @0 :Type;

    enum Specificity {
        one @0;
        everything @1;
    }
    amount @1 :Specificity;
    id @2 :Text;
}

struct GetConfigurationResponse @0xad93807af77fe1b9 {
    union {
        all @0 :AllConfigurations;
        raft @1 :List(RaftConfiguration);
    }
}

struct AllConfigurations @0xa3cf4f7f955be932 {
    raft @0 :List(RaftConfiguration);
}

struct PutConfigurationRequest @0x93c59921a137c8db {
    enum Type {
        raft @0;
        nodeHost @1;
    }
    union {
        raft @0 :RaftConfiguration;
        nodeHost @1 :NodeHostConfiguration;
    }
}

struct PutConfigurationResponse @0x9f8f8f8f8f8f8f8f {
    enum Type {
            raft @0;
            nodeHost @1;
    }
    union {
        raft @0 :RaftConfiguration;
        nodeHost @1 :NodeHostConfiguration;
    }
    success @2 :Bool;
    status @3 :Text;
    type @4 :Type;
}

# ServiceType is the initial message payload sent by the client to the server so the stream can be

# mapped to a specific server implementation.

struct ServiceType @0x94e84f47e297127c {
    enum Type {
        test @0;
        configService @1;
    }
    type @0 :Type;
}

################################################################################
#             __ _
#  _ __ __ _ / _| |_
# | '__/ _` | |_| __|
# | | | (_| |  _| |_
# |_|  \__,_|_|  \__|
#
################################################################################

# RaftConfig is the configuration for a Raft node.
struct RaftConfiguration @0xdb9a661a7821150d {
    id @0 :Text;
    nodeId @1 :UInt64;
    clusterId @2 :UInt64;
    checkQuorum @3 :Bool;
    electionTimeout @4 :UInt64;
    heartbeatTimeout @5 :UInt64;
    snapshotEntries @6 :UInt64;
    compactionOverhead @7 :UInt64;
    orderedConfigurationChange @8 :Bool;
    maxInMemoryLogSize @9 :UInt64;
    snapshotCompressionType @10 :UInt64;
    entryCompressionType @11 :UInt64;
    disableAutoCompaction @12 :Bool;
    isObserver @13 :Bool;
    isWitness @14 :Bool;
    quiesce @15 :Bool;
    configType @16 :ConfigType;
}

enum ConfigType @0xadb2d0b69445303c {
    system @0;
    exchange @1;
    customerFsm @2;
}

# ListRaftConfigsRequest is the request for listing Raft nodes.
struct ListRaftConfigurationRequest @0xc75a0f30e41b37f5 {}

# ListRaftConfigsResponse is the response for listing Raft nodes.
struct ListRaftConfigurationResponse @0xc5195060b33b8218 {
    configs @0 :List(RaftConfiguration);
}

# GetRaftConfigRequest is the request for getting a Raft configuration.
struct GetRaftConfigurationRequest @0xa233c1204c18c976 {
    id @0 :Text;
}

# GetRaftConfigResponse is the response for getting a Raft configuration.
struct GetRaftConfigurationResponse @0xb4fe5e6f0ef85636 {
    config @0 :RaftConfiguration;
}

# UpdateRaftConfigRequest is the request for updating or creating a Raft configuration.
struct PutRaftConfigurationRequest @0xefe67d057faf5d90 {
    enable @0 :Bool;
    name @1 :Text;
    config @2 :RaftConfiguration;

}

# UpdateRaftConfigResponse is the response for updating or creating a Raft configuration.
struct PutRaftConfigurationResponse @0x8f8f8f8f8f8f8f8f {
    valid @0 :Bool;
    name @1 :Text;
    error @2 :Text;
}

struct NodeHostConfiguration @0x859698645e9c4a44 {
    deploymentId @0 :UInt64;
    writeAheadLogDir @1 :Text;
    nodeHostDir @2 :Text;
    roundTripTimeMilliseconds @3 :UInt64;
    raftAddress @4 :Text;
    apiAddress @5 :Text;
    mutualTls @6 :Bool;
    caFile @7 :Text;
    certFile @8 :Text;
    keyFile @9 :Text;
}