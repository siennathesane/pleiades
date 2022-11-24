/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package messaging

import (
	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/mxplusb/pleiades/pkg/messaging/clients"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var workflowSingleton *EmbeddedWorkflowServer

func NewEmbeddedWorkflowServer(streamClient *clients.EmbeddedMessagingStreamClient, logger zerolog.Logger) (*EmbeddedWorkflowServer, error) {
	if workflowSingleton != nil {
		return workflowSingleton, nil
	}

	systemStore, err := systemstore.NewSystemStore(logger)
	if err != nil {
		return nil, err
	}

	dataStore, err := runtime.NewWorkflowStore(systemStore, logger)
	if err != nil {
		return nil, err
	}

	stateStore, err := runtime.NewWorkflowStateStore(systemStore, logger)
	if err != nil {
		return nil, err
	}

	workflowSingleton = &EmbeddedWorkflowServer{
		logger:     logger.With().Str("component", "workflow-server").Logger(),
		dataStore:  dataStore,
		stateStore: stateStore,
	}

	return workflowSingleton, nil
}

type EmbeddedWorkflowServer struct {
	workerConcurrency int
	retryCount        int
	flows             map[string]workflow.WorkflowDefinitionHandler
	dataStore         systemstore.DataStore
	stateStore        systemstore.StateStore
	logger            zerolog.Logger
	runtime           *workflow.WorkflowRuntime
	streamClient      *clients.EmbeddedMessagingStreamClient
}

func (fs *EmbeddedWorkflowServer) WorkerConcurrency() int {
	return fs.workerConcurrency
}

func (fs *EmbeddedWorkflowServer) RetryCount() int {
	return fs.retryCount
}

func (fs *EmbeddedWorkflowServer) Flows() map[string]workflow.WorkflowDefinitionHandler {
	return fs.flows
}

type Request struct {
	Body      []byte
	RequestId string
	Query     map[string][]string
	Header    map[string][]string
}

const (
	defaultWorkerConcurrency = 2
	defaultRetryCount        = 2
)

func (fs *EmbeddedWorkflowServer) Execute(flowName string, req *Request) error {
	if flowName == "" {
		return errors.Newf("flowName must be provided to execute flow")
	}

	fs.configureDefault()
	fs.runtime = &workflow.WorkflowRuntime{}

	request := &runtime.Request{
		Header:    req.Header,
		RequestID: req.RequestId,
		Body:      req.Body,
		Query:     req.Query,
	}

	err := fs.runtime.Execute(flowName, request)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	return nil
}

func (fs *EmbeddedWorkflowServer) Pause(flowName string, requestId string) error {
	if flowName == "" {
		return errors.Newf("flowName must be provided")
	}

	if requestId == "" {
		return errors.Newf("request Id must be provided")
	}

	fs.configureDefault()
	fs.runtime = &workflow.WorkflowRuntime{}

	request := &runtime.Request{
		RequestID: requestId,
	}

	err := fs.runtime.Pause(flowName, request)
	if err != nil {
		return errors.Newf("failed to pause request, %v", err)
	}

	return nil
}

func (fs *EmbeddedWorkflowServer) Resume(flowName string, requestId string) error {
	if flowName == "" {
		return errors.Newf("flowName must be provided")
	}

	if requestId == "" {
		return errors.Newf("request Id must be provided")
	}

	fs.configureDefault()
	fs.runtime = &workflow.WorkflowRuntime{}

	request := &runtime.Request{
		RequestID: requestId,
	}

	err := fs.runtime.Resume(flowName, request)
	if err != nil {
		return errors.Newf("failed to resume request, %v", err)
	}

	return nil
}

func (fs *EmbeddedWorkflowServer) Stop(flowName string, requestId string) error {
	if flowName == "" {
		return errors.Newf("flowName must be provided")
	}

	if requestId == "" {
		return errors.Newf("request Id must be provided")
	}

	fs.configureDefault()
	fs.runtime = &workflow.WorkflowRuntime{}

	request := &runtime.Request{
		RequestID: requestId,
	}

	err := fs.runtime.Stop(flowName, request)
	if err != nil {
		return errors.Newf("failed to stop request, %v", err)
	}

	return nil
}

func (fs *EmbeddedWorkflowServer) Register(flowName string, handler workflow.WorkflowDefinitionHandler) error {
	if flowName == "" {
		return errors.Newf("flow-name must not be empty")
	}
	if handler == nil {
		return errors.Newf("handler must not be nil")
	}

	if fs.flows == nil {
		fs.flows = make(map[string]workflow.WorkflowDefinitionHandler)
	}

	if fs.flows[flowName] != nil {
		return errors.Newf("flow-name must be unique for each flow")
	}

	fs.flows[flowName] = handler

	return nil
}

func (fs *EmbeddedWorkflowServer) Start() error {
	if len(fs.flows) == 0 {
		return errors.Newf("must register atleast one flow")
	}
	fs.configureDefault()
	fs.runtime = workflow.NewWorkflowRuntime(fs.flows, fs.stateStore, fs.dataStore, defaultWorkerConcurrency, defaultRetryCount, fs.streamClient, fs.logger)

	errorChan := make(chan error)
	defer close(errorChan)
	if err := fs.initRuntime(); err != nil {
		return err
	}

	go fs.runtimeWorker(errorChan)
	go fs.queueWorker(errorChan)
	err := <-errorChan
	return err
}

func (fs *EmbeddedWorkflowServer) StartServer() error {
	fs.configureDefault()
	fs.runtime = workflow.NewWorkflowRuntime(fs.flows, fs.stateStore, fs.dataStore, defaultWorkerConcurrency, defaultRetryCount, fs.streamClient, fs.logger)

	errorChan := make(chan error)
	defer close(errorChan)
	if err := fs.initRuntime(); err != nil {
		return err
	}

	go fs.runtimeWorker(errorChan)
	err := <-errorChan
	return errors.Newf("server has stopped, error: %v", err)
}

func (fs *EmbeddedWorkflowServer) StartWorker() error {
	fs.configureDefault()
	fs.runtime = workflow.NewWorkflowRuntime(fs.flows, fs.stateStore, fs.dataStore, defaultWorkerConcurrency, defaultRetryCount, fs.streamClient, fs.logger)

	errorChan := make(chan error)
	defer close(errorChan)
	if err := fs.initRuntime(); err != nil {
		return err
	}

	go fs.runtimeWorker(errorChan)
	go fs.queueWorker(errorChan)
	err := <-errorChan
	return errors.Newf("worker has stopped, error: %v", err)
}

func (fs *EmbeddedWorkflowServer) configureDefault() {
	if fs.workerConcurrency == 0 {
		fs.workerConcurrency = defaultWorkerConcurrency
	}
	if fs.retryCount == 0 {
		fs.retryCount = defaultRetryCount
	}
}

func (fs *EmbeddedWorkflowServer) initRuntime() error {
	err := fs.runtime.Init()
	if err != nil {
		return err
	}
	return nil
}

func (fs *EmbeddedWorkflowServer) runtimeWorker(errorChan chan error) {
	err := fs.runtime.StartRuntime()
	errorChan <- errors.Newf("runtime has stopped, error: %v", err)
}

func (fs *EmbeddedWorkflowServer) queueWorker(errorChan chan error) {
	err := fs.runtime.StartQueueWorker(errorChan)
	errorChan <- errors.Newf("worker has stopped, error: %v", err)
}
