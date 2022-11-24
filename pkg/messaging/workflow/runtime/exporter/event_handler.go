/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package exporter

// EventHandler handle flow events
type EventHandler interface {
	// Configure the EventHandler with flow name and request ID
	Configure(flowName string, requestId string)
	// Initialize an EventHandler (called only once in a request span)
	Init() error
	// ReportRequestStart report a start of request
	ReportRequestStart(requestId string)
	// ReportRequestEnd reports an end of request
	ReportRequestEnd(requestId string)
	// ReportRequestFailure reports a failure of a request with error
	ReportRequestFailure(requestId string, err error)
	// ReportExecutionForward report that an execution is forwarded
	ReportExecutionForward(nodeId string, requestId string)
	// ReportExecutionContinuation report that an execution is being continued
	ReportExecutionContinuation(requestId string)
	// ReportNodeStart report a start of a Node execution
	ReportNodeStart(nodeId string, requestId string)
	// ReportNodeStart report an end of a node execution
	ReportNodeEnd(nodeId string, requestId string)
	// ReportNodeFailure report a Node execution failure with error
	ReportNodeFailure(nodeId string, requestId string, err error)
	// ReportOperationStart reports start of an operation
	ReportOperationStart(operationId string, nodeId string, requestId string)
	// ReportOperationEnd reports an end of an operation
	ReportOperationEnd(operationId string, nodeId string, requestId string)
	// ReportOperationFailure reports failure of an operation with error
	ReportOperationFailure(operationId string, nodeId string, requestId string, err error)
	// Flush flush the reports
	Flush()
}
