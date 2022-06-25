/*
 * Copyright (c) 2022 Sienna Lloyd <sienna.lloyd@hey.com>
 */

package fsm

import (
	"bytes"
	"errors"
	"reflect"

	"capnproto.org/go/capnp/v3"
	"github.com/rs/zerolog"
	v1 "r3t.io/pleiades/pkg/protocols/config/v1"
	"r3t.io/pleiades/pkg/servers/services"
)

var (
	_ services.IStore[v1.RaftConfiguration] = (*ConfigServiceStoreManager)(nil)
)

type ConfigServiceStoreManager struct {
	logger zerolog.Logger
	store  *services.StoreManager
}

func NewConfigServiceStoreManager(logger zerolog.Logger, store *services.StoreManager) (*ConfigServiceStoreManager, error) {
	cssm := &ConfigServiceStoreManager{logger: logger, store: store}
	err := store.Start(false)
	if err != nil {
		return nil, err
	}
	return cssm, nil
}

func (r *ConfigServiceStoreManager) Get(key string) (*v1.RaftConfiguration, error) {
	payload, err := r.store.Get(key, reflect.TypeOf(&v1.RaftConfiguration{}))
	if err != nil {
		r.logger.Err(err).Str("key", key).Msg("error fetching key from raft store")
		return &v1.RaftConfiguration{}, err
	}

	msgPayload, err := capnp.NewDecoder(bytes.NewReader(payload)).Decode()
	if err != nil {
		r.logger.Err(err).Str("key", key).Msg("error decoding payload")
		return &v1.RaftConfiguration{}, err
	}

	config, err := v1.ReadRootRaftConfiguration(msgPayload)
	if err != nil {
		r.logger.Err(err).Str("key", key).Msg("error reading payload")
		return &v1.RaftConfiguration{}, err
	}

	return &config, nil
}

func (r *ConfigServiceStoreManager) GetAll() (map[string]*v1.RaftConfiguration, error) {
	payloads, err := r.store.GetAll(reflect.TypeOf(&v1.RaftConfiguration{}))
	if err != nil {
		r.logger.Err(err).Msg("error fetching all from raft store")
		return nil, err
	}

	configs := make(map[string]*v1.RaftConfiguration)
	for key, payload := range payloads {
		msgPayload, err := capnp.NewDecoder(bytes.NewReader(payload)).Decode()
		if err != nil {
			r.logger.Err(err).Str("key", key).Msg("error decoding payload")
			return nil, err
		}

		config, err := v1.ReadRootRaftConfiguration(msgPayload)
		if err != nil {
			r.logger.Err(err).Str("key", key).Msg("error reading payload")
			return nil, err
		}

		configs[key] = &config
	}

	return configs, nil
}

func (r *ConfigServiceStoreManager) Put(key string, payload *v1.RaftConfiguration) error {
	if payload == nil {
		err := errors.New("payload is nil")
		r.logger.Err(err).Str("key", key).Msg("error putting nil payload")
		return err
	}

	var buf bytes.Buffer
	err := capnp.NewEncoder(&buf).Encode(payload.Message())
	if err != nil {
		r.logger.Err(err).Msg("error encoding payload")
		return err
	}

	err = r.store.Put(key, buf.Bytes(), reflect.TypeOf(&v1.RaftConfiguration{}))
	if err != nil {
		r.logger.Err(err).Str("key", key).Msg("error putting key into raft store")
		return err
	}

	return nil
}
