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
	"encoding/json"
)

// ExecutorRequest defines the body of async forward request to core
type ExecutorRequest struct {
	Sign        string `json: "sign"`         // request signature
	ID          string `json: "id"`           // request ID
	Query       string `json: "query"`        // query string
	CallbackUrl string `json: "callback-url"` // callback url

	ExecutionState string `json: "state"` // Execution State (execution position / execution vertex)

	Data []byte `json: "data"` // Partial execution data
	// (empty if intermediate_storage enabled

	ContextStore map[string][]byte `json: "store"` // Context State for default dataStore
	// (empty if external Store is used)
}

func buildRequest(id string,
	state string,
	query string,
	data []byte,
	contextState map[string][]byte,
	sign string) *ExecutorRequest {

	request := &ExecutorRequest{
		Sign:           sign,
		ID:             id,
		ExecutionState: state,
		Query:          query,
		Data:           data,
		ContextStore:   contextState,
	}
	return request
}

func decodeRequest(data []byte) (*ExecutorRequest, error) {
	request := &ExecutorRequest{}
	err := json.Unmarshal(data, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (req *ExecutorRequest) encode() ([]byte, error) {
	return json.Marshal(req)
}

func (req *ExecutorRequest) getData() []byte {
	return req.Data
}

func (req *ExecutorRequest) getID() string {
	return req.ID
}

func (req *ExecutorRequest) getExecutionState() string {
	return req.ExecutionState
}

func (req *ExecutorRequest) getContextStore() map[string][]byte {
	return req.ContextStore
}

func (req *ExecutorRequest) getQuery() string {
	return req.Query
}
