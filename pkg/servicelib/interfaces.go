package servicelib

import (
	"gonum.org/v1/gonum/graph"
)

type LifecycleServiceType int

const (
	// TransientServiceType is used for things like RPC requests
	TransientServiceType LifecycleServiceType = 0
	// ScopedServiceType is used for things like individual processes
	ScopedServiceType LifecycleServiceType = 1
	// SingletonServiceType is used for globally unique things
	SingletonServiceType LifecycleServiceType = 2
)

type Service interface {
	graph.Node
	SetNodeID(nid int64)
	GetServiceName() string
	GetServiceType() LifecycleServiceType
	MarkDependencies(deps []Service) error
	GetDependencies() []Service
	PrepareToRun() error
	ReadyToRun() bool
	IsRunning() bool
	Start(retry bool) error
	Stop(retry, force bool) error
}
