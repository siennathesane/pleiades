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
//)
//
//type RaftConfig struct {
//	NodeID                  uint64                 `json:"node_id,omitempty"`
//	ClusterID               uint64                 `json:"cluster_id,omitempty"`
//	CheckQuorum             bool                   `json:"check_quorum,omitempty"`
//	ElectionRTT             uint64                 `json:"election_rtt,omitempty"`
//	HeartbeatRTT            uint64                 `json:"heartbeat_rtt,omitempty"`
//	SnapshotEntries         uint64                 `json:"snapshot_entries,omitempty"`
//	CompactionOverhead      uint64                 `json:"compaction_overhead,omitempty"`
//	OrderedConfigChange     bool                   `json:"ordered_config_change,omitempty"`
//	MaxInMemLogSize         uint64                 `json:"max_in_mem_log_size,omitempty"`
//	SnapshotCompressionType config.CompressionType `json:"snapshot_compression_type,omitempty"`
//	EntryCompressionType    config.CompressionType `json:"entry_compression_type,omitempty"`
//	DisableAutoCompactions  bool                   `json:"disable_auto_compactions,omitempty"`
//	IsObserver              bool                   `json:"is_observer,omitempty"`
//	IsWitness               bool                   `json:"is_witness,omitempty"`
//	Quiesce                 bool                   `json:"quiesce,omitempty"`
//}
//
//// NewRaftConfig generates a new Raft node configuration
//// todo (sienna): this should read from 169.254.169.254
//func NewRaftConfig(client *etcd.Client, sugaredLogger sugaredLogger.ILogger) (config.Config, error) {
//	configLogger = sugaredLogger
//
//	hostname, err := os.Hostname()
//	if err != nil {
//		return config.Config{}, err
//	}
//
//	pair, _, err := client.KV().Get(fmt.Sprintf("hosts/%s/conf/raft", hostname), &etcd.QueryOptions{})
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
//	return config.Config(c), nil
//}
//
//// Validate validates the Config instance and return an error when any member
//// field is considered as invalid.
//func (c RaftConfig) validate() error {
//
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
