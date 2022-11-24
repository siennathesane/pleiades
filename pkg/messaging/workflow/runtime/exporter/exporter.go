/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package exporter

import (
	"fmt"

	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/execution"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/graph"
)

// Exporter
type Exporter interface {
	// GetFlowName get name of the flow
	GetFlowName() string
	// GetFlowDefinition get definition of the faas-flow
	GetFlowDefinition(*graph.Pipeline, *execution.Context) error
}

// FlowExporter core exporter
type FlowExporter struct {
	flow     *graph.Pipeline
	flowName string
	exporter Exporter // exporter
}

// createContext create a context from request handler
func (fexp *FlowExporter) createContext() *execution.Context {
	context := execution.CreateContext("export", "",
		fexp.flowName, nil)

	return context
}

// Export retrieve core definition
func (fexp *FlowExporter) Export() ([]byte, error) {

	// Init flow
	fexp.flow = graph.CreatePipeline()
	fexp.flowName = fexp.exporter.GetFlowName()

	context := fexp.createContext()

	// Get definition: Get Pipeline definition from user implemented Define()
	err := fexp.exporter.GetFlowDefinition(fexp.flow, context)
	if err != nil {
		return nil, fmt.Errorf("Failed to define flow, %v", err)
	}

	definition := graph.GetPipelineDefinition(fexp.flow)

	return []byte(definition), nil
}

// CreateFlowExporter initiate a FlowExporter with a provided IExecutor
func CreateFlowExporter(exporter Exporter) (fexp *FlowExporter) {
	fexp = &FlowExporter{}
	fexp.exporter = exporter

	return fexp
}
