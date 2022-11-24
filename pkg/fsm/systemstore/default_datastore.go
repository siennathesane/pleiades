/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package systemstore

import (
	"github.com/cockroachdb/errors"
)

// json to encode
type InMemoryDataStore struct {
	store map[string][]byte
}

func (rstore *InMemoryDataStore) Store() map[string][]byte {
	return rstore.store
}

func (rstore *InMemoryDataStore) SetStore(store map[string][]byte) {
	rstore.store = store
}

// NewInMemoryDataStore creates a new InMemoryDataStore
func NewInMemoryDataStore() *InMemoryDataStore {
	rstore := &InMemoryDataStore{}
	rstore.store = make(map[string][]byte)
	return rstore
}

// NewInMemoryDataStoreFromExisting creates a store manager from a map
func NewInMemoryDataStoreFromExisting(store map[string][]byte) *InMemoryDataStore {
	rstore := &InMemoryDataStore{}
	rstore.store = store
	return rstore
}

// Configure with requestId and flowname
func (rstore *InMemoryDataStore) Configure(flowName string, requestId string) {

}

// Init initialize the storemanager (called only once in a request span)
func (rstore *InMemoryDataStore) Init() error {
	return nil
}

// Set sets a value (implement dataStore)
func (rstore *InMemoryDataStore) Set(key string, value []byte) error {
	rstore.store[key] = value
	return nil
}

// Get gets a value (implement dataStore)
func (rstore *InMemoryDataStore) Get(key string) ([]byte, error) {
	value, ok := rstore.store[key]
	if !ok {
		return nil, errors.Newf("no field name %s", key)
	}
	return value, nil
}

// Del delets a value (implement dataStore)
func (rstore *InMemoryDataStore) Del(key string) error {
	if _, ok := rstore.store[key]; ok {
		delete(rstore.store, key)
	}
	return nil
}

// Cleanup
func (rstore *InMemoryDataStore) Cleanup() error {
	return nil
}
