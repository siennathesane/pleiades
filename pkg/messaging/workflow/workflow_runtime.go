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
	"encoding/json"
	"fmt"
	"time"

	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/mxplusb/pleiades/pkg/messaging/clients"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/execution"
	"github.com/mxplusb/pleiades/pkg/messaging/workflow/runtime/exporter"
	"github.com/cockroachdb/errors"
	"github.com/go-co-op/gocron"
	"github.com/nats-io/nats.go"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type WorkflowRuntime struct {
	flows           map[string]WorkflowDefinitionHandler
	stateStore      systemstore.StateStore
	dataStore       systemstore.DataStore
	logger          zerolog.Logger
	concurrency     int
	retryQueueCount int
	jobScheduler    *gocron.Scheduler
	updateJob       *gocron.Job
	streamClient    *clients.EmbeddedMessagingStreamClient
	taskQueues      map[string]*nats.Subscription
	done            chan struct{}
}

func NewWorkflowRuntime(flows map[string]WorkflowDefinitionHandler,
	stateStore systemstore.StateStore,
	dataStore systemstore.DataStore,
	concurrency int,
	retryQueueCount int,
	streamClient *clients.EmbeddedMessagingStreamClient,
	logger zerolog.Logger) *WorkflowRuntime {
	return &WorkflowRuntime{
		flows:           flows,
		stateStore:      stateStore,
		dataStore:       dataStore,
		logger:          logger,
		streamClient:    streamClient,
		concurrency:     concurrency,
		retryQueueCount: retryQueueCount,
	}
}

type Worker struct {
	ID          string   `json:"id"`
	Flows       []string `json:"flows"`
	Concurrency int      `json:"concurrency"`
}

type Task struct {
	FlowName    string              `json:"flow_name"`
	RequestID   string              `json:"request_id"`
	Body        string              `json:"body"`
	Header      map[string][]string `json:"header"`
	RawQuery    string              `json:"raw_query"`
	Query       map[string][]string `json:"query"`
	RequestType string              `json:"request_type"`
}

const (
	InternalRequestQueueInitial = "system.workflows"
	FlowKeyInitial              = "flow"
	WorkerKeyInitial            = "worker"

	GoFlowRegisterInterval = 4

	PartialRequest = "PARTIAL"
	NewRequest     = "NEW"
	PauseRequest   = "PAUSE"
	ResumeRequest  = "RESUME"
	StopRequest    = "STOP"
)

func (fr *WorkflowRuntime) Init() error {
	var err error

	err = fr.stateStore.Init()
	if err != nil {
		return errors.Newf("failed to initialize the stateStore, %v", err)
	}

	err = fr.dataStore.Init()
	if err != nil {
		return errors.Newf("failed to initialize the stateStore, %v", err)
	}

	return nil
}

// StartQueueWorker starts listening for request in queue
func (fr *WorkflowRuntime) StartQueueWorker(errorChan chan error) error {
	fr.taskQueues = make(map[string]*nats.Subscription)
	for flowName := range fr.flows {
		streamName := fr.internalRequestQueueId(flowName)
		_, err := fr.streamClient.AddStream(&nats.StreamConfig{
			Name:        streamName,
			Description: fmt.Sprintf("Message queue for %s", streamName),
		})

		index := 0
		for index < fr.concurrency {
			consumerName := fmt.Sprintf("%s-consumer-%d", flowName, index)
			_, err = fr.streamClient.AddConsumer(streamName, &nats.ConsumerConfig{
				Description: consumerName,
			})
			if err != nil {
				return errors.Newf("failed to add consumer, error %v", err)
			}

			sub, err := fr.streamClient.QueueSubscribe(streamName, streamName, fr.Consume, nats.Bind(streamName, consumerName))
			if err != nil {
				return errors.Newf("failed to register consumer, error %v", err)
			}
			fr.taskQueues[flowName] = sub
			index++
		}
	}

	fr.logger.Info().Msg("queue worker started successfully")
	return nil
}

// StartRuntime starts the runtime
func (fr *WorkflowRuntime) StartRuntime() error {
	worker := &Worker{
		ID:          getNewId(),
		Flows:       make([]string, 0, len(fr.flows)),
		Concurrency: fr.concurrency,
	}
	// Get the flow details for each flow
	flowDetails := make(map[string]string)
	for flowID, defHandler := range fr.flows {
		worker.Flows = append(worker.Flows, flowID)
		dag, err := getFlowDefinition(defHandler)
		if err != nil {
			return errors.Newf("failed to start runtime, dag export failed, error %v", err)
		}
		flowDetails[flowID] = dag
	}
	err := fr.saveWorkerDetails(worker)
	if err != nil {
		return errors.Newf("failed to register worker details, %v", err)
	}
	err = fr.saveFlowDetails(flowDetails)
	if err != nil {
		return errors.Newf("failed to register worker details, %v", err)
	}

	tz, err := time.LoadLocation("UTC")
	if err != nil {
		fr.logger.Error().Err(err).Msg("failed to load tz database")
		return err
	}

	fr.jobScheduler = gocron.NewScheduler(tz)
	fr.updateJob, err = fr.jobScheduler.Every(GoFlowRegisterInterval).Second().Do(func() {
		var err error
		err = fr.saveWorkerDetails(worker)
		if err != nil {
			fr.logger.Error().Err(err).Msg("failed to register worker details")
		}
		err = fr.saveFlowDetails(flowDetails)
		if err != nil {
			fr.logger.Error().Err(err).Msg("failed to register worker details")
		}
	})
	if err != nil {
		return errors.Newf("failed to start runtime, %v", err)
	}

	fr.jobScheduler.StartAsync()
	fr.logger.Info().Msg("runtime started")

	return nil
}

func (fr *WorkflowRuntime) CreateExecutor(req *runtime.Request) (*WorkflowExecutor, error) {
	flowHandler, ok := fr.flows[req.FlowName]
	if !ok {
		return nil, errors.Newf("could not find handler for flow %s", req.FlowName)
	}
	ex := &WorkflowExecutor{
		StateStore: fr.stateStore,
		DataStore:  fr.dataStore,
		Handler:    flowHandler,
		Logger:     fr.logger,
		Runtime:    fr,
	}
	err := ex.Init(req)
	return ex, err
}

func (fr *WorkflowRuntime) Execute(flowName string, request *runtime.Request) error {
	streamName := fr.internalRequestQueueId(flowName)
	task := &Task{
		FlowName:    flowName,
		RequestID:   request.RequestID,
		Body:        string(request.Body),
		Header:      request.Header,
		RawQuery:    request.RawQuery,
		Query:       request.Query,
		RequestType: NewRequest,
	}
	data, _ := json.Marshal(task)

	_, err := fr.streamClient.Publish(streamName, data)
	if err != nil {
		fr.logger.Error().Str("workflow", flowName).Interface("request", task).Err(err).Msg("failed to publish task")
		return errors.Wrap(err, "failed to publish task")
	}
	return nil
}

func (fr *WorkflowRuntime) Pause(flowName string, request *runtime.Request) error {
	streamName := fr.internalRequestQueueId(flowName)
	task := &Task{
		FlowName:    flowName,
		RequestID:   request.RequestID,
		Body:        string(request.Body),
		Header:      request.Header,
		RawQuery:    request.RawQuery,
		Query:       request.Query,
		RequestType: PauseRequest,
	}
	data, _ := json.Marshal(task)

	_, err := fr.streamClient.Publish(streamName, data)
	if err != nil {
		fr.logger.Error().Str("workflow", flowName).Interface("request", task).Err(err).Msg("failed to publish task")
		return errors.Wrap(err, "failed to publish task")
	}
	return nil
}

func (fr *WorkflowRuntime) Stop(flowName string, request *runtime.Request) error {
	streamName := fr.internalRequestQueueId(flowName)
	task := &Task{
		FlowName:    flowName,
		RequestID:   request.RequestID,
		Body:        string(request.Body),
		Header:      request.Header,
		RawQuery:    request.RawQuery,
		Query:       request.Query,
		RequestType: StopRequest,
	}
	data, _ := json.Marshal(task)

	_, err := fr.streamClient.Publish(streamName, data)
	if err != nil {
		fr.logger.Error().Str("workflow", flowName).Interface("request", task).Err(err).Msg("failed to publish task")
		return errors.Wrap(err, "failed to publish task")
	}
	return nil
}

func (fr *WorkflowRuntime) Resume(flowName string, request *runtime.Request) error {
	streamName := fr.internalRequestQueueId(flowName)
	task := &Task{
		FlowName:    flowName,
		RequestID:   request.RequestID,
		Body:        string(request.Body),
		Header:      request.Header,
		RawQuery:    request.RawQuery,
		Query:       request.Query,
		RequestType: ResumeRequest,
	}
	data, _ := json.Marshal(task)

	_, err := fr.streamClient.Publish(streamName, data)
	if err != nil {
		fr.logger.Error().Str("workflow", flowName).Interface("request", task).Err(err).Msg("failed to publish task")
		return errors.Wrap(err, "failed to publish task")
	}
	return nil
}

func (fr *WorkflowRuntime) EnqueuePartialRequest(pr *runtime.Request) error {
	streamName := fr.internalRequestQueueId(pr.FlowName)
	task := &Task{
		FlowName:    pr.FlowName,
		RequestID:   pr.RequestID,
		Body:        string(pr.Body),
		Header:      pr.Header,
		RawQuery:    pr.RawQuery,
		Query:       pr.Query,
		RequestType: PartialRequest,
	}

	data, _ := json.Marshal(task)

	_, err := fr.streamClient.Publish(streamName, data)
	if err != nil {
		fr.logger.Error().Str("workflow", pr.FlowName).Interface("request", task).Err(err).Msg("failed to publish task")
		return errors.Wrap(err, "failed to publish task")
	}
	return nil
}

// Consume messages from queue
func (fr *WorkflowRuntime) Consume(message *nats.Msg) {
	resend := func() {
		if err := message.Nak(); err != nil {
			fr.logger.Error().Err(err).Msg("failed to push message back to queue")
		}
	}

	var task Task
	if err := json.Unmarshal(message.Data, &task); err != nil {
		fr.logger.Error().Err(err).Msg("failed to unmarshal payload")
		resend()
	} else {
		if err = fr.handleRequest(makeRequestFromTask(task), task.RequestType); err != nil {
			fr.logger.Error().Err(err).Msg("rejecting task for failure")
			resend()
		}

		err = message.Ack()
		if err != nil {
			fr.logger.Error().Err(err).Msg("failed to acknowledge message")
			return
		}
	}
}

func (fr *WorkflowRuntime) handleRequest(request *runtime.Request, requestType string) error {
	l := fr.logger.With().Str("request-id", request.RequestID).Str("workflow", request.FlowName).Logger()

	flowExecutor, err := fr.CreateExecutor(request)
	if err != nil {
		l.Error().Err(err).Msg("failed to execute request")
		return errors.Newf("failed to execute request " + request.RequestID + ", error: " + err.Error())
	}

	rawRequest := &execution.RawRequest{}
	rawRequest.Data = request.Body
	rawRequest.Query = request.RawQuery
	if request.RequestID != "" {
		rawRequest.RequestId = request.RequestID
	}
	stateOption := execution.NewRequest(rawRequest)

	ex := execution.CreateFlowExecutor(flowExecutor, nil)
	_, err = ex.Execute(stateOption)
	if err != nil {
		l.Error().Err(err).Msg("failed to to execute")
	}

	switch requestType {
	case PartialRequest:
		_, err = ex.Execute(stateOption)
		break
	case NewRequest:
		_, err = ex.Execute(stateOption)
		break
	case PauseRequest:
		err = ex.Pause(request.RequestID)
		break
	case ResumeRequest:
		err = ex.Resume(request.RequestID)
		break
	case StopRequest:
		err = ex.Stop(request.RequestID)
		break
	default:
		return errors.Newf("invalid request %v received with type %s", request, requestType)
	}
	return errors.Wrap(err, "unable to handle request")
}

func (fr *WorkflowRuntime) internalRequestQueueId(flowName string) string {
	return fmt.Sprintf("%s.%s", InternalRequestQueueInitial, flowName)
}

func (fr *WorkflowRuntime) requestQueueId(flowName string) string {
	return flowName
}

func (fr *WorkflowRuntime) saveWorkerDetails(worker *Worker) error {
	key := fmt.Sprintf("%s:%s", WorkerKeyInitial, worker.ID)
	value := marshalWorker(worker)
	return errors.Wrap(fr.dataStore.Set(key, []byte(value)), "failed to save worker details")
}

func (fr *WorkflowRuntime) saveFlowDetails(flows map[string]string) error {
	for flowId, definition := range flows {
		key := fmt.Sprintf("%s:%s", FlowKeyInitial, flowId)
		errors.Wrap(fr.dataStore.Set(key, []byte(definition)), "failed to save flow details")
	}
	return nil
}

func marshalWorker(worker *Worker) string {
	jsonDef, _ := json.Marshal(worker)
	return string(jsonDef)
}

func makeRequestFromTask(task Task) *runtime.Request {
	request := &runtime.Request{
		FlowName:  task.FlowName,
		RequestID: task.RequestID,
		Body:      []byte(task.Body),
		Header:    task.Header,
		RawQuery:  task.RawQuery,
		Query:     task.Query,
	}
	return request
}

func getFlowDefinition(handler WorkflowDefinitionHandler) (string, error) {
	ex := &WorkflowExecutor{
		Handler: handler,
	}
	flowExporter := exporter.CreateFlowExporter(ex)
	resp, err := flowExporter.Export()
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func getNewId() string {
	guid := xid.New()
	return guid.String()
}
