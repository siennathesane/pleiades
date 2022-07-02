
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package servicegraph

import (
	"fmt"

	"gonum.org/v1/gonum/graph/simple"
)

var (
	_ IServiceLibrary = (*ServiceLibrary)(nil)
)

type ServiceLibrary struct {
	graph *simple.DirectedGraph
}

func NewServiceLibrary() (*ServiceLibrary, error) {
	sl := &ServiceLibrary{
		simple.NewDirectedGraph(),
	}
	return sl, nil
}

func (sl *ServiceLibrary) AddService(svc Service) error {
	if svc.GetServiceName() == "" {
		return fmt.Errorf("cannot add service without a name")
	}

	svc.SetNodeID(sl.graph.NewNode().ID())
	sl.graph.AddNode(svc)
	if err := sl.registerDeps(svc); err != nil {
		return err
	}

	return nil
}

func (sl *ServiceLibrary) registerDeps(svc Service) error {
	if len(svc.GetDependencies()) == 0 {
		return nil
	}

	for _, dep := range svc.GetDependencies() {
		dep.SetNodeID(sl.graph.NewNode().ID())
		sl.graph.NewEdge(dep, svc)
		if err := sl.registerDeps(dep); err != nil {
			return err
		}
	}

	return nil
}

func (sl *ServiceLibrary) AddServices(svcs []Service) error {
	for _, svc := range svcs {
		if err := sl.AddService(svc); err != nil {
			return err
		}
	}
	return nil
}

func (sl *ServiceLibrary) GetService(svc Service) (Service, error) {
	if svc.GetServiceName() == "" || svc.ID() == 0 {
		return nil, fmt.Errorf("cannot find a service without a name or node id")
	}

	nodes := sl.graph.Nodes()
	for nodes.Len() > 0 {
		if nodes.Next() {
			// we know this will always be a service, no need to type check
			target := nodes.Node().(Service)
			if target.GetServiceName() == svc.GetServiceName() || target.ID() == svc.ID() {
				return target, nil
			}
		}
	}

	return nil, fmt.Errorf("cannot locate service %s with service id of %d", svc.GetServiceName(), svc.ID())
}

func (sl *ServiceLibrary) StartService(svc Service, retry bool) error {
	edges := sl.graph.Edges()
	for edges.Len() > 0 {
		if edges.Next() {
			child := edges.Edge().From()
			var currentChild Service
			if child != nil {
				currentChild = child.(Service)
			}

			// slow ass recursion
			if len(currentChild.GetDependencies()) > 0 {
				if err := sl.StartService(currentChild, retry); err != nil {
					return fmt.Errorf("error starting dependent service %s: %w", currentChild.GetServiceName(), err)
				}
			}

			if !currentChild.ReadyToRun() {
				if err := currentChild.PrepareToRun(); err != nil {
					return err
				}
			}
			if err := currentChild.Start(retry); err != nil {
				return err
			}
			if !currentChild.IsRunning() {
				return fmt.Errorf("service %s is not running", currentChild.GetServiceName())
			}
		}
	}
	return nil
}

func (sl *ServiceLibrary) StopService(retry, force bool, svc Service) error {
	return nil
}
