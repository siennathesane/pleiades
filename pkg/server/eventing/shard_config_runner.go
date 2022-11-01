/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package eventing

import (
	"context"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/fsm"
	"github.com/cockroachdb/errors"
	"github.com/lni/dragonboat/v3"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewShardConfigRunner(nodeHost *dragonboat.NodeHost, logger zerolog.Logger) (*shardConfigRunner, error) {
	l := logger.With().Str("component", "shard-config").Logger()
	srv, err := NewServer(logger)
	if err != nil {
		return nil, err
	}

	store, err := fsm.NewShardStore(logger)
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	runner := &shardConfigRunner{logger: l, msgServ: srv, store: store, done: done, nh: nodeHost}
	go runner.run()

	return runner, nil
}

type shardConfigRunner struct {
	logger  zerolog.Logger
	msgServ *Server
	store   *fsm.ShardStore
	nh      *dragonboat.NodeHost
	done    chan struct{}
}

func (s *shardConfigRunner) run() {
	client, err := s.msgServ.GetPubSubClient()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("can't talk to nats")
	}

	sub, err := client.SubscribeSync(ShardConfigStream)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("can't subscribe to stream")
	}
	defer sub.Unsubscribe()

	for {
		msg, err := sub.Fetch(1)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			s.logger.Error().Err(err).Msg("cant fetch message")
		}
		if len(msg) == 0 {
			continue
		}
		go s.handleMsg(msg[0])
		select {
		case <-s.done:
			return
		}
	}
}

func (s *shardConfigRunner) handleMsg(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		s.logger.Error().Err(err).Msg("error acknowledging message processing")
	}

	payload := &raftv1.ShardStateEvent{}
	err = payload.UnmarshalVT(msg.Data)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't unmarshal shard message payload")
		return
	}

	if payload.GetEvent() == nil {
		s.logger.Error().Msg("empty event")
		return
	}

	switch payload.GetCmd() {
	case raftv1.ShardStateEvent_CMD_TYPE_PUT:

		event := payload.GetEvent()
		var memberState *dragonboat.Membership

		for memberState == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			memberState, err = s.nh.SyncGetClusterMembership(ctx, event.GetShardId())
			if err != nil {
				s.logger.Error().Err(err).Msg("can't get shard state")
				cancel()
				continue
			}
			cancel()
		}

		shardState := &raftv1.ShardState{
			LastUpdated:    timestamppb.Now(),
			ShardId:        payload.GetEvent().GetShardId(),
			ConfigChangeId: memberState.ConfigChangeID,
			Replicas:       memberState.Nodes,
			Observers:      memberState.Observers,
			Witnesses:      memberState.Witnesses,
			Removed: func() map[uint64]string {
				m := make(map[uint64]string)
				for k := range memberState.Removed {
					m[k] = ""
				}
				return m
			}(),
			Type: 0,
		}

		s.logger.Trace().
			Int64("last-updated", shardState.LastUpdated.Seconds).
			Uint64("shard-id", shardState.GetShardId()).
			Uint64("config-change-id", shardState.GetConfigChangeId()).
			Msg("shard state")

		if err := s.store.Put(shardState); err != nil {
			s.logger.Error().Err(err).Msg("can't store shard event")
			break
		}

	case raftv1.ShardStateEvent_CMD_TYPE_DELETE:
		if err := s.store.Delete(payload.GetEvent().GetShardId()); err != nil {
			s.logger.Error().Err(err).Msg("can't store shard event")
			break
		}
	default:
		break
	}
}
