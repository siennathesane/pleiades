package servers

import (
	"context"
	"errors"
	"fmt"
	"math"

	dconf "github.com/lni/dragonboat/v3/config"
	dlog "github.com/lni/dragonboat/v3/logger"
	"github.com/pbnjay/memory"
	"go.uber.org/atomic"
	"r3t.io/pleiades/pkg"
	"r3t.io/pleiades/pkg/fsm"
	configv1 "r3t.io/pleiades/pkg/pb/config/v1"
)

var (
	maxInMemoryLogSize uint64
)

func init() {
	func(target *uint64) {
		total := memory.TotalMemory()
		// subtract 25% from the total available memory, then divide it evenly amongst the max amount of
		// clusters per hosts. better to be safe upfront than allow for noisy neighbours
		totalMinus25Percent := float64(total - uint64(math.Round(float64(total)*0.25)))
		*target = uint64(math.Floor(totalMinus25Percent / float64(pkg.MaxRaftClustersPerHost)))
	}(&maxInMemoryLogSize)
}

type RaftConfigServer struct {
	configv1.DRPCConfigServiceUnimplementedServer
	manager *fsm.RaftManager[configv1.RaftConfig]
	logger  dlog.ILogger

	allCache map[string]*configv1.RaftConfig
	count    *atomic.Uint64
}

func NewRaftConfigServer(manager *fsm.RaftManager[configv1.RaftConfig], logger dlog.ILogger) *RaftConfigServer {
	count := atomic.NewUint64(0)
	all, _ := manager.GetAll()
	if all != nil {
		count.Add(uint64(len(all)))
	}
	return &RaftConfigServer{manager: manager, logger: logger, count: count}
}

func (rcs *RaftConfigServer) PutConfiguration(ctx context.Context, payload *configv1.PutRaftConfigurationRequest) (*configv1.PutRaftConfigurationResponse, error) {
	if err := rcs.validatePutRaftConfig(payload); err != nil {
		errMsg := err.Error()
		return &configv1.PutRaftConfigurationResponse{Name: &payload.Name, Error: &errMsg, Valid: false}, err
	}

	// prep the in-memory cache
	rcs.allCache[payload.Name] = payload.Config

	return &configv1.PutRaftConfigurationResponse{Name: &payload.Name, Error: nil, Valid: true}, nil
}

func (rcs *RaftConfigServer) validatePutRaftConfig(payload *configv1.PutRaftConfigurationRequest) error {
	if payload == nil {
		return errors.New("cannot validate an empty configuration")
	}

	config := payload.Config

	switch config.Type {
	case configv1.IsConfigType_IS_CONFIG_TYPE_SYSTEM:
		return rcs.validateSystemCluster(payload)
	}

	return nil
}

func (rcs *RaftConfigServer) validateSystemCluster(payload *configv1.PutRaftConfigurationRequest) error {
	config := payload.Config

	if config.ClusterId < pkg.SystemClusterIdRangeStartingValue || config.ClusterId > pkg.SystemClusterIdRangeEndingValue {
		return fmt.Errorf("a system cluster must have a cluster range between %d and %d", pkg.SystemClusterIdRangeStartingValue, pkg.SystemClusterIdRangeEndingValue)
	}

	if config.NodeId == 0 {
		return errors.New("invalid NodeID, it must be >= 1")
	}

	if config.HeartbeatRoundTripTime == 0 {
		return errors.New("HeartbeatRoundTripTime must be > 0")
	}

	if config.ElectionRoundTripTime == 0 {
		return errors.New("ElectionRoundTripTime must be > 0")
	}

	if config.ElectionRoundTripTime <= 2*config.HeartbeatRoundTripTime {
		return errors.New("invalid ElectionRoundTripTime, must be at least two times the value of HeartbeatRoundTripTime")
	}

	if config.ElectionRoundTripTime > 10*config.HeartbeatRoundTripTime {
		rcs.logger.Warningf("ElectionRoundTripTime needs to be less than 10 times the value of HeartbeatRoundTripTime")
	}

	if config.GetMaxInMemoryLogSize() > 0 && config.GetMaxInMemoryLogSize() > maxInMemoryLogSize {
		return fmt.Errorf("MaxInMemLogSize must be more than 0 and less than %d for this host", maxInMemoryLogSize)
	}

	if dconf.CompressionType(config.SnapshotCompressionType) != dconf.Snappy &&
		dconf.CompressionType(config.SnapshotCompressionType) != dconf.NoCompression {
		return errors.New("unknown compression type")
	}

	if config.IsWitness && config.SnapshotEntries > 0 {
		return errors.New("witness node can not take snapshot")
	}

	if config.IsWitness && config.IsObserver {
		return errors.New("witness node can not be an observer")
	}

	if err := rcs.manager.Put(payload.Name, payload.Config); err != nil {
		return err
	}

	return nil
}

func (rcs *RaftConfigServer) GetConfiguration(ctx context.Context, payload *configv1.GetRaftConfigurationRequest) (*configv1.GetRaftConfigurationResponse, error) {
	if payload == nil {
		return nil, errors.New("cannot fetch an empty configuration")
	}

	if payload.Name == "" {
		return nil, errors.New("cannot fetch an empty configuration")
	}

	// try the fast path
	val, ok := rcs.allCache[payload.Name]
	if ok {
		return &configv1.GetRaftConfigurationResponse{Configuration: val}, nil
	}

	// now we do the slow path
	val, err := rcs.manager.Get(payload.Name)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, fmt.Errorf("raft configuration %s cannot be found", payload.Name)
	}

	return nil, nil
}

func (rcs *RaftConfigServer) ListConfigurations(ctx context.Context, payload *configv1.ListRaftConfigurationRequest) (*configv1.ListRaftConfigurationResponse, error) {
	var err error
	all := make(map[string]*configv1.RaftConfig)
	if len(rcs.allCache) == 0 {
		all, err = rcs.manager.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		all = rcs.allCache
	}

	return &configv1.ListRaftConfigurationResponse{AvailableConfigs: all}, nil
}

func (rcs *RaftConfigServer) mustEmbedUnimplementedRaftConfigServiceServer() {}
