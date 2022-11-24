/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package workflow

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/execution"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/exporter"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/graph"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

var _ execution.IExecutor = (*WorkflowExecutor)(nil)

type WorkflowDefinitionHandler func(flow *runtime.Workflow, context *execution.Context) error

type WorkflowExecutor struct {
	gateway          string
	flowName         string // the name of the function
	reqID            string // the request id
	CallbackURL      string // the callback url
	IsLoggingEnabled bool
	partialState     []byte
	rawRequest       *execution.RawRequest
	StateStore       systemstore.StateStore
	DataStore        systemstore.DataStore
	EventHandler     exporter.EventHandler
	Logger           zerolog.Logger
	Handler          WorkflowDefinitionHandler
	Runtime          *WorkflowRuntime
}

func (fe *WorkflowExecutor) GetFlowName() string {
	return fe.flowName
}

func (fe *WorkflowExecutor) HandleNextNode(partial *execution.PartialState) error {
	var err error
	request := &runtime.Request{}
	request.Body, err = partial.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode partial state, error %v")
	}
	request.RequestID = fe.reqID
	request.FlowName = fe.flowName
	request.Header = make(map[string][]string)

	err = fe.Runtime.EnqueuePartialRequest(request)
	if err != nil {
		return errors.Wrap(err, "failed to enqueue request, error %v")
	}
	return nil
}

func (fe *WorkflowExecutor) GetExecutionOption(operation execution.IOperation) map[string]interface{} {
	options := make(map[string]interface{})
	options["gateway"] = fe.gateway
	options["request-id"] = fe.reqID

	return options
}

func (fe *WorkflowExecutor) HandleExecutionCompletion(data []byte) error {
	if fe.CallbackURL == "" {
		return nil
	}

	fe.Logger.Info().Str("callback-url", fe.CallbackURL).Msg("executing callback")
	httpreq, _ := http.NewRequest(http.MethodPost, fe.CallbackURL, bytes.NewReader(data))
	httpreq.Header.Add("X-Faas-Flow-ReqiD", fe.reqID)
	client := &http.Client{}

	res, resErr := client.Do(httpreq)
	if resErr != nil {
		return resErr
	}
	defer res.Body.Close()
	resData, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return errors.Newf("failed to call callback %d: %s", res.StatusCode, string(resData))
	}

	return nil
}

func (fe *WorkflowExecutor) Configure(requestID string) {
	fe.reqID = requestID
}

func (fe *WorkflowExecutor) GetFlowDefinition(pipeline *graph.Pipeline, context *execution.Context) error {
	workflow := runtime.GetWorkflow(pipeline)
	return fe.Handler(workflow, context)
}

func (fe *WorkflowExecutor) ReqValidationEnabled() bool {
	return false
}

func (fe *WorkflowExecutor) GetValidationKey() (string, error) {
	return "", nil
}

func (fe *WorkflowExecutor) ReqAuthEnabled() bool {
	return false
}

func (fe *WorkflowExecutor) GetReqAuthKey() (string, error) {
	return "", nil
}

func (fe *WorkflowExecutor) MonitoringEnabled() bool {
	return false
}

func (fe *WorkflowExecutor) GetEventHandler() (exporter.EventHandler, error) {
	return fe.EventHandler, nil
}

func (fe *WorkflowExecutor) LoggingEnabled() bool {
	return fe.IsLoggingEnabled
}

func (fe *WorkflowExecutor) GetLogger() (zerolog.Logger, error) {
	return fe.Logger, nil
}

func (fe *WorkflowExecutor) GetStateStore() (systemstore.StateStore, error) {
	return fe.StateStore, nil
}

func (fe *WorkflowExecutor) GetDataStore() (systemstore.DataStore, error) {
	return fe.DataStore, nil
}

func (fe *WorkflowExecutor) Init(request *runtime.Request) error {
	fe.flowName = request.FlowName

	callbackURL := request.GetHeader("X-Faas-Flow-Callback-Url")
	fe.CallbackURL = callbackURL

	return nil
}
