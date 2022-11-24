/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

import (
	"fmt"

	"github.com/mxplusb/pleiades/pkg/fsm/systemstore"
	"github.com/rs/zerolog"
	"github.com/s8sg/goflow/core/sdk"
)

var (
	_ sdk.DataStore  = (*WorkflowStore)(nil)
	_ sdk.StateStore = (*WorkflowStateStore)(nil)
)

func NewWorkflowStore(store *systemstore.SystemStore, logger zerolog.Logger) (*WorkflowStore, error) {
	return &WorkflowStore{
		logger: logger.With().Str("component", "workflow-store").Logger(),
		store:  store,
	}, nil
}

type WorkflowStore struct {
	logger     zerolog.Logger
	store      *systemstore.SystemStore
	bucketName string
}

func (w *WorkflowStore) Configure(flowName string, requestId string) {
	w.bucketName = fmt.Sprintf(WorkflowStateBucketFormat, flowName, requestId)
}

func (w *WorkflowStore) Init() error {
	return nil
}

func (w *WorkflowStore) Set(key string, value []byte) error {
	return w.store.Put(w.bucketName, key, value)
}

func (w *WorkflowStore) Get(key string) ([]byte, error) {
	return w.store.Get(w.bucketName, key)
}

func (w *WorkflowStore) Del(key string) error {
	return w.store.Delete(w.bucketName, key)
}

func (w *WorkflowStore) Cleanup() error {
	return w.store.DeleteBucket(w.bucketName)
}

func NewWorkflowStateStore(store *systemstore.SystemStore, logger zerolog.Logger) (*WorkflowStateStore, error) {
	return &WorkflowStateStore{
		logger: logger.With().Str("component", "workflow-store").Logger(),
		store:  store,
	}, nil
}

type WorkflowStateStore struct {
	logger     zerolog.Logger
	store      *systemstore.SystemStore
	bucketName string
}

func (w *WorkflowStateStore) Configure(flowName string, requestId string) {
	w.bucketName = fmt.Sprintf(WorkflowStateBucketFormat, flowName, requestId)
}

func (w *WorkflowStateStore) Init() error {
	return nil
}

func (w *WorkflowStateStore) Set(key string, value string) error {
	return w.store.Put(w.bucketName, key, []byte(value))
}

func (w *WorkflowStateStore) Get(key string) (string, error) {
	resp, err := w.store.Get(w.bucketName, key)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (w *WorkflowStateStore) Update(key string, oldValue string, newValue string) error {
	return w.store.Put(w.bucketName, key, []byte(newValue))
}

func (w *WorkflowStateStore) Cleanup() error {
	return w.store.DeleteBucket(w.bucketName)
}
