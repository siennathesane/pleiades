/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package execution

import (
	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/graph"
	"github.com/rs/zerolog"
)

// ExecutionRuntime implements how operation executed and handle next nodes in async
type ExecutionRuntime interface {
	// HandleNextNode handles execution of next nodes based on partial state
	HandleNextNode(state *PartialState) (err error)
	// Provide an execution option that will be passed to the operation
	GetExecutionOption(operation IOperation) map[string]interface{}
	// Handle the completion of execution of data
	HandleExecutionCompletion(data []byte) error
}

// IExecutor implements a faas-flow executor
type IExecutor interface {
	// Configure configure an executor with request id
	Configure(requestId string)
	// GetFlowName get name of the flow
	GetFlowName() string
	// GetFlowDefinition get definition of the faas-flow
	GetFlowDefinition(*graph.Pipeline, *Context) error
	// ReqValidationEnabled check if request validation enabled
	ReqValidationEnabled() bool
	// GetValidationKey get request validation key
	GetValidationKey() (string, error)
	// ReqAuthEnabled check if request auth enabled
	ReqAuthEnabled() bool
	// GetReqAuthKey get the request auth key
	GetReqAuthKey() (string, error)
	// MonitoringEnabled check if request monitoring enabled
	MonitoringEnabled() bool
	// LoggingEnabled check if logging is enabled
	LoggingEnabled() bool
	// GetLogger get the logger
	GetLogger() (zerolog.Logger, error)
	// GetStateStore get the state store
	GetStateStore() (systemstore.StateStore, error)
	// GetDataStore get the data store
	GetDataStore() (systemstore.DataStore, error)

	ExecutionRuntime
}

type IOperation interface {
	GetId() string
	Encode() []byte
	GetProperties() map[string][]string
	// Execute executes an operation, executor can pass configuration
	Execute([]byte, map[string]interface{}) ([]byte, error)
}
