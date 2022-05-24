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
	"r3t.io/pleiades/pkg/pb"
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
	pb.UnimplementedRaftConfigServiceServer
	manager *fsm.RaftManager[pb.RaftConfig]
	logger  dlog.ILogger

	allCache map[string]*pb.RaftConfig
	count    *atomic.Uint64
}

func NewRaftConfigServer(manager *fsm.RaftManager[pb.RaftConfig], logger dlog.ILogger) *RaftConfigServer {
	count := atomic.NewUint64(0)
	all, _ := manager.GetAll()
	if all != nil {
		count.Add(uint64(len(all)))
	}
	return &RaftConfigServer{manager: manager, logger: logger, count: count}
}

func (rcs *RaftConfigServer) PutConfiguration(ctx context.Context, payload *pb.PutRaftConfigRequest) (*pb.NewRaftConfigResponse, error) {
	if err := rcs.validatePutRaftConfig(payload); err != nil {
		errMsg := err.Error()
		return &pb.NewRaftConfigResponse{Name: &payload.Name, Error: &errMsg, Valid: false}, err
	}

	// prep the in-memory cache
	rcs.allCache[payload.Name] = payload.Config

	return &pb.NewRaftConfigResponse{Name: &payload.Name, Error: nil, Valid: true}, nil
}

func (rcs *RaftConfigServer) validatePutRaftConfig(payload *pb.PutRaftConfigRequest) error {
	if payload == nil {
		return errors.New("cannot validate an empty configuration")
	}

	config := payload.Config

	switch config.Type {
	case pb.IsConfigType_System:
		return rcs.validateSystemCluster(payload)
	}

	return nil
}

func (rcs *RaftConfigServer) validateSystemCluster(payload *pb.PutRaftConfigRequest) error {
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

	if config.MaxInMemLogSize > 0 && config.MaxInMemLogSize > maxInMemoryLogSize {
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

func (rcs *RaftConfigServer) GetConfiguration(ctx context.Context, payload *pb.GetRaftConfigRequest) (*pb.GetRaftConfigResponse, error) {
	if payload == nil {
		return nil, errors.New("cannot fetch an empty configuration")
	}

	if payload.Name == "" {
		return nil, errors.New("cannot fetch an empty configuration")
	}

	// try the fast path
	val, ok := rcs.allCache[payload.Name]
	if ok {
		return &pb.GetRaftConfigResponse{Configuration: val}, nil
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

func (rcs *RaftConfigServer) ListConfigurations(ctx context.Context, payload *pb.ListRaftConfigsRequest) (*pb.ListRaftConfigsResponse, error) {
	var err error
	all := make(map[string]*pb.RaftConfig)
	if len(rcs.allCache) == 0 {
		all, err = rcs.manager.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		all = rcs.allCache
	}

	return &pb.ListRaftConfigsResponse{AvailableConfigs: all}, nil
}

func (rcs *RaftConfigServer) mustEmbedUnimplementedRaftConfigServiceServer() {}
