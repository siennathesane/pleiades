/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package kv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mxplusb/pleiades/pkg/api/v1/database"
	"github.com/mxplusb/pleiades/pkg/configuration"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.etcd.io/bbolt"
)

type op int

const (
	get op = 1
	put op = 2
)

var (
	_ statemachine.IOnDiskStateMachine = &BBoltStateMachine{}

	ErrBadUpdate = errors.New("bad update")
)

type BBoltStateMachine struct {
	ShardId   uint64
	ReplicaId uint64

	logger zerolog.Logger
	store  *bboltStore
}

func NewBBoltStateMachine(shardId, replicaId uint64) *BBoltStateMachine {
	return &BBoltStateMachine{
		ShardId:   shardId,
		ReplicaId: replicaId,
		logger: configuration.NewRootLogger().
			With().
			Str("component", "bbolt-fsm").
			Uint64("shard", shardId).
			Uint64("replica", replicaId).
			Logger(),
	}
}

// Open the bbolt backend and read the last index.
func (b *BBoltStateMachine) Open(stopc <-chan struct{}) (uint64, error) {

	basePath := viper.GetString("datastore.basePath")
	basePath = filepath.Join(basePath, fmt.Sprintf("shard-%d-replica-%d.db", b.ShardId, b.ReplicaId))

	store, err := newBboltStore(b.ShardId, b.ReplicaId, basePath, b.logger)
	if err != nil {
		return 0, err
	}
	b.store = store

	return store.GetMonotonicLogIndex()
}

func (b *BBoltStateMachine) Update(entries []statemachine.Entry) ([]statemachine.Entry, error) {
	applied := make([]statemachine.Entry, 0)

	if len(entries) == 0 {
		return applied, nil
	}

	for _, entry := range entries {
		wrapper := &database.KVStoreWrapper{}
		if err := wrapper.UnmarshalVT(entry.Cmd); err != nil {
			b.logger.Error().Err(err).Msg("can't unmarshal kv store wrapper")
			return applied, errors.Wrap(err, "can't unmarshal kv store wrapper")
		}

		switch wrapper.Typ {
		case database.KVStoreWrapper_CREATE_ACCOUNT_REQUEST:
			req := wrapper.GetCreateAccountRequest()

			resp, err := b.store.CreateAccountBucket(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.Wrap(err, "can't create account bucket")
			}
			wrapper.Typ = database.KVStoreWrapper_CREATE_ACCOUNT_REPLY
			wrapper.Payload = &database.KVStoreWrapper_CreateAccountReply{
				CreateAccountReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.Wrap(err, "can't serialize response wrapper")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break

		case database.KVStoreWrapper_DELETE_ACCOUNT_REQUEST:
			req := wrapper.GetDeleteAccountRequest()

			resp, err := b.store.DeleteAccountBucket(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't delete account")
				return applied, errors.Wrap(err, "can't delete account")
			}
			wrapper.Typ = database.KVStoreWrapper_DELETE_ACCOUNT_REPLY
			wrapper.Payload = &database.KVStoreWrapper_DeleteAccountReply{
				DeleteAccountReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't delete account bucket")
				return applied, errors.Wrap(err, "can't delete account bucket")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break

		case database.KVStoreWrapper_CREATE_BUCKET_REQUEST:
			req := wrapper.GetCreateBucketRequest()

			resp, err := b.store.CreateBucket(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create bucket")
				return applied, errors.Wrap(err, "can't create bucket")
			}
			wrapper.Typ = database.KVStoreWrapper_CREATE_BUCKET_REPLY
			wrapper.Payload = &database.KVStoreWrapper_CreateBucketReply{
				CreateBucketReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.Wrap(err, "can't serialize response wrapper")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break

		case database.KVStoreWrapper_DELETE_BUCKET_REQUEST:
			req := wrapper.GetDeleteBucketRequest()

			resp, err := b.store.DeleteBucket(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.Wrap(err, "can't create account bucket")
			}
			wrapper.Typ = database.KVStoreWrapper_DELETE_BUCKET_REPLY
			wrapper.Payload = &database.KVStoreWrapper_DeleteBucketReply{
				DeleteBucketReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't delete bucket")
				return applied, errors.Wrap(err, "can't delete bucket")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break

		case database.KVStoreWrapper_PUT_KEY_REQUEST:
			req := wrapper.GetPutKeyRequest()

			resp, err := b.store.PutKey(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.New("can't create account bucket")
			}
			wrapper.Typ = database.KVStoreWrapper_PUT_KEY_REPLY
			wrapper.Payload = &database.KVStoreWrapper_PutKeyReply{
				PutKeyReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.Wrap(err, "can't serialize response wrapper")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break

		case database.KVStoreWrapper_DELETE_KEY_REQUEST:
			req := wrapper.GetDeleteKeyRequest()

			resp, err := b.store.DeleteKey(req)
			if err != nil {
				b.logger.Error().Err(err).Msg("can't create account bucket")
				return applied, errors.New("can't create account bucket")
			}
			wrapper.Typ = database.KVStoreWrapper_DELETE_KEY_REPLY
			wrapper.Payload = &database.KVStoreWrapper_DeleteKeyReply{
				DeleteKeyReply: resp,
			}
			serialized, err := wrapper.MarshalVT()
			if err != nil {
				b.logger.Error().Err(err).Msg("can't delete account bucket")
				return applied, errors.Wrap(err, "can't delete response wrapper")
			}

			applied = append(applied, statemachine.Entry{
				Index:  entry.Index,
				Result: statemachine.Result{
					Value: entry.Index,
					Data: serialized,
				},
			})

			break
		default:
			return applied, errors.New("unsupported method")
		}
	}

	return applied, nil
}

func (b *BBoltStateMachine) Lookup(i interface{}) (interface{}, error) {
	wrapper := &database.KVStoreWrapper{}

	err := wrapper.UnmarshalVT(i.([]byte))
	if err != nil {
		b.logger.Error().Err(err).Msg("can't unmarshal payload")
		return []byte{}, errors.Wrap(err, "bad request")
	}

	switch wrapper.Typ {
	case database.KVStoreWrapper_GET_KEY_REQUEST:
		req := wrapper.GetGetKeyRequest()
		var resp *database.GetKeyReply
		resp, err = b.store.GetKey(req)
		wrapper.Typ = database.KVStoreWrapper_GET_KEY_REPLY
		wrapper.Payload = &database.KVStoreWrapper_GetKeyReply{GetKeyReply: resp}
		break

	default:
		return []byte{}, errors.New("unsupported method")
	}

	buf, err := wrapper.MarshalVT()
	if err != nil {
		b.logger.Error().Err(err).Msg("error fetching data")
	}
	return buf, err
}

func (b *BBoltStateMachine) Sync() error {
	return b.store.db.Sync()
}

func (b *BBoltStateMachine) PrepareSnapshot() (interface{}, error) {
	return nil, nil
}

func (b *BBoltStateMachine) SaveSnapshot(ctx interface{}, writer io.Writer, done <-chan struct{}) error {
	return b.store.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.WriteTo(writer)
		return err
	})
}

func (b *BBoltStateMachine) RecoverFromSnapshot(reader io.Reader, i <-chan struct{}) error {
	basePath := viper.GetString("datastore.basePath")
	basePath = filepath.Join(basePath, fmt.Sprintf("shard-%d-replica-%d.db", b.ShardId, b.ReplicaId))

	fn := func(r io.Reader) error {
		target, err := os.Create(basePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(target, reader)
		if err != nil {
			return err
		}
		return nil
	}

	// verify the existing database is closed
	err := b.store.db.Close()
	if err != nil {
		return err
	}

	_, err = os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fn(reader)
		}
		return err
	}

	err = os.Remove(basePath)
	if err != nil {
		return err
	}

	return fn(reader)
}

func (b *BBoltStateMachine) Close() error {
	if b.store != nil {
		if b.store.db != nil {
			err := b.store.db.Close()
			if err != nil {
				return err
			}
		} else {
			b.logger.Panic().Msg("no reference to bbolt")
		}
	} else {
		b.logger.Panic().Msg("no reference to bbolt")
	}
	return nil
}
