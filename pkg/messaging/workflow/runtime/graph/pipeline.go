/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package graph

import (
	"encoding/json"
	"fmt"
)

const (
	DepthIncrement = 1
	DepthDecrement = -1
	DepthSame      = 0
)

// PipelineErrorHandler the error handler OnFailure() registration on pipeline
type PipelineErrorHandler func(error) ([]byte, error)

// PipelineHandler definition for the Finally() registration on pipeline
type PipelineHandler func(string)

type Pipeline struct {
	Dag *Dag `json:"-"` // Dag that will be executed

	ExecutionPosition map[string]string `json:"pipeline-execution-position"` // Denotes the node that is executing now
	ExecutionDepth    int               `json:"pipeline-execution-depth"`    // Denotes the depth of subgraph its executing

	CurrentDynamicOption map[string]string `json:"pipeline-dynamic-option"` // Denotes the current dynamic option mapped against the dynamic Node UQ id

	FailureHandler PipelineErrorHandler `json:"-"`
	Finally        PipelineHandler      `json:"-"`
}

// CreatePipeline creates a core pipeline
func CreatePipeline() *Pipeline {
	pipeline := &Pipeline{}
	pipeline.Dag = NewDag()

	pipeline.ExecutionPosition = make(map[string]string, 0)
	pipeline.ExecutionDepth = 0
	pipeline.CurrentDynamicOption = make(map[string]string, 0)

	return pipeline
}

// CountNodes counts the no of node added in the Pipeline Dag.
// It doesn't count subdags node
func (p *Pipeline) CountNodes() int {
	return len(p.Dag.nodes)
}

// GetAllNodesUniqueId returns a recursive list of all nodes that belongs to the pipeline
func (p *Pipeline) GetAllNodesUniqueId() []string {
	nodes := p.Dag.GetNodes("")
	return nodes
}

// GetInitialNodeId Get the very first node of the pipeline
func (p *Pipeline) GetInitialNodeId() string {
	node := p.Dag.GetInitialNode()
	if node != nil {
		return node.Id
	}
	return "0"
}

// GetNodeExecutionUniqueId provide a ID that is unique in an execution
func (p *Pipeline) GetNodeExecutionUniqueId(node *Node) string {
	depth := 0
	dag := p.Dag
	depthStr := ""
	optionStr := ""
	for depth < p.ExecutionDepth {
		depthStr = fmt.Sprintf("%d", depth)
		node := dag.GetNode(p.ExecutionPosition[depthStr])
		option := p.CurrentDynamicOption[node.GetUniqueId()]
		if node.subDag != nil {
			dag = node.subDag
		} else {
			dag = node.conditionalDags[option]
		}
		if optionStr == "" {
			optionStr = option
		} else {
			optionStr = option + "--" + optionStr
		}

		depth++
	}
	if optionStr == "" {
		return node.GetUniqueId()
	}
	return optionStr + "--" + node.GetUniqueId()
}

// GetCurrentNodeDag returns the current node and current dag based on execution position
func (p *Pipeline) GetCurrentNodeDag() (*Node, *Dag) {
	depth := 0
	dag := p.Dag
	depthStr := ""
	for depth < p.ExecutionDepth {
		depthStr = fmt.Sprintf("%d", depth)
		node := dag.GetNode(p.ExecutionPosition[depthStr])
		option := p.CurrentDynamicOption[node.GetUniqueId()]
		if node.subDag != nil {
			dag = node.subDag
		} else {
			dag = node.conditionalDags[option]
		}
		depth++
	}
	depthStr = fmt.Sprintf("%d", depth)
	node := dag.GetNode(p.ExecutionPosition[depthStr])
	return node, dag
}

// UpdatePipelineExecutionPosition updates pipeline execution position
// specified depthAdjustment and vertex denotes how the ExecutionPosition must be altered
func (p *Pipeline) UpdatePipelineExecutionPosition(depthAdjustment int, vertex string) {
	p.ExecutionDepth = p.ExecutionDepth + depthAdjustment
	depthStr := fmt.Sprintf("%d", p.ExecutionDepth)
	p.ExecutionPosition[depthStr] = vertex
}

// SetDag overrides the default dag
func (p *Pipeline) SetDag(dag *Dag) {
	p.Dag = dag
}

// decodePipeline decodes a json marshaled pipeline
func decodePipeline(data []byte) (*Pipeline, error) {
	pipeline := &Pipeline{}
	err := json.Unmarshal(data, pipeline)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// GetState get a state of a pipeline by encoding in JSON
func (p *Pipeline) GetState() string {
	encode, _ := json.Marshal(p)
	return string(encode)
}

// ApplyState apply a state to a pipeline by from encoded JSON pipeline
func (p *Pipeline) ApplyState(state string) {
	temp, _ := decodePipeline([]byte(state))
	p.ExecutionDepth = temp.ExecutionDepth
	p.ExecutionPosition = temp.ExecutionPosition
	p.CurrentDynamicOption = temp.CurrentDynamicOption
}
