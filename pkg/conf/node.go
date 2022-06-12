package conf

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"os"
//
//	"github.com/hashicorp/consul/etcd"
//	"github.com/lni/dragonboat/v3/config"
//	"github.com/lni/dragonboat/v3/sugaredLogger"
//	"github.com/lni/dragonboat/v3/raftio"
//)
//
//var (
//	configLogger sugaredLogger.ILogger
//)
//
//type NodeConfig struct {
//	DeploymentID                  uint64                      `json:"deployment_id"`
//	WALDir                        string                      `json:"wal_dir"`
//	NodeHostDir                   string                      `json:"node_host_dir"`
//	RTTMillisecond                uint64                      `json:"rtt_millisecond"`
//	RaftAddress                   string                      `json:"raft_address"`
//	AddressByNodeHostID           bool                        `json:"address_by_node_host_id"`
//	ListenAddress                 string                      `json:"listen_address"`
//	MutualTLS                     bool                        `json:"mutual_tls"`
//	CAFile                        string                      `json:"ca_file"`
//	CertFile                      string                      `json:"cert_file"`
//	KeyFile                       string                      `json:"key_file"`
//	EnableMetrics                 bool                        `json:"enable_metrics"`
//	RaftEventListener             raftio.IRaftEventListener   `json:"-"`
//	SystemEventListener           raftio.ISystemEventListener `json:"-"`
//	MaxSendQueueSize              uint64                      `json:"max_send_queue_size,omitempty"`
//	MaxReceiveQueueSize           uint64                      `json:"max_receive_queue_size,omitempty"`
//	MaxSnapshotSendBytesPerSecond uint64                      `json:"max_snapshot_send_bytes_per_second,omitempty"`
//	MaxSnapshotRecvBytesPerSecond uint64                      `json:"max_snapshot_recv_bytes_per_second,omitempty"`
//	NotifyCommit                  bool                        `json:"notify_commit,omitempty"`
//	Gossip                        config.GossipConfig         `json:"-"`
//	Expert                        config.ExpertConfig         `json:"-"`
//}
//
//// NewNodeConfig generates a new node configuration for multiraft
//// todo (sienna): this should read from 169.254.169.254
//func NewNodeConfig(client *etcd.Client, sugaredLogger sugaredLogger.ILogger) (config.Config, error) {
//	configLogger = sugaredLogger
//
//	hostname, err := os.Hostname()
//	if err != nil {
//		return config.Config{}, err
//	}
//
//	pair, _, err := client.KV().Get(fmt.Sprintf("hosts/%s/conf/raft-conf", hostname), &etcd.QueryOptions{})
//	if err != nil {
//		configLogger.Errorf("can't get configuration from consul: %s", err)
//		return config.Config{}, err
//	}
//
//	var c RaftConfig
//	if err := json.Unmarshal(pair.Value, &c); err != nil {
//		configLogger.Errorf("can't unmarshal configuration from consul", err)
//		return config.Config{}, err
//	}
//
//	// and validate it
//	if err := validate(c); err != nil {
//		return config.Config{}, err
//	}
//
//	// yay conversions!
//	return config.NodeHostConfig(c), nil
//}
//
//// Validate validates the Config instance and return an error when any member
//// field is considered as invalid.
//func validate(c RaftConfig) error {
//	if c.NodeID == 0 {
//		return errors.New("invalid NodeID, it must be >= 1")
//	}
//	if c.HeartbeatRTT == 0 {
//		return errors.New("HeartbeatRTT must be > 0")
//	}
//	if c.ElectionRTT == 0 {
//		return errors.New("ElectionRTT must be > 0")
//	}
//	if c.ElectionRTT <= 2*c.HeartbeatRTT {
//		return errors.New("invalid election rtt")
//	}
//	if c.ElectionRTT < 10*c.HeartbeatRTT {
//		configLogger.Warningf("election_rtt is not a magnitude larger than heartbeat_rtt")
//	}
//	// todo (sienna): figure out how to validate `MaxInMemLogSize`
//	//if c.MaxInMemLogSize > 0 &&
//	//	c.MaxInMemLogSize < settings.EntryNonCmdFieldsSize+1 {
//	//	return errors.New("MaxInMemLogSize is too small")
//	//}
//	if c.SnapshotCompressionType != config.Snappy &&
//		c.SnapshotCompressionType != config.NoCompression {
//		return errors.New("unknown compression type")
//	}
//	if c.EntryCompressionType != config.Snappy &&
//		c.EntryCompressionType != config.NoCompression {
//		return errors.New("unknown compression type")
//	}
//	if c.IsWitness && c.SnapshotEntries > 0 {
//		return errors.New("witness node can not take snapshot")
//	}
//	if c.IsWitness && c.IsObserver {
//		return errors.New("witness node can not be an observer")
//	}
//	return nil
//}
