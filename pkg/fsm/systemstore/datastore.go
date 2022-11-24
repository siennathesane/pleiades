/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package systemstore

// DataStore for Storing Data
type DataStore interface {
	// Configure the DaraStore with flow name and request ID
	Configure(flowName string, requestId string)
	// Initialize the DataStore (called only once in a request span)
	Init() error
	// Set store a value for key, in failure returns error
	Set(key string, value []byte) error
	// Get retrieves a value by key, if failure returns error
	Get(key string) ([]byte, error)
	// Del deletes a value by a key
	Del(key string) error
	// Cleanup all the resources in DataStore
	Cleanup() error
}

// StateStore for saving execution state
type StateStore interface {
	// Configure the StateStore with flow name and request ID
	Configure(flowName string, requestId string)
	// Initialize the StateStore (called only once in a request span)
	Init() error
	// Set a value (override existing, or create one)
	Set(key string, value string) error
	// Get a value
	Get(key string) (string, error)
	// Compare and Update a value
	Update(key string, oldValue string, newValue string) error
	// Cleanup all the resources in StateStore (called only once in a request span)
	Cleanup() error
}
