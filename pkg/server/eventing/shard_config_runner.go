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
	"strings"
	"time"

	raftv1 "github.com/mxplusb/pleiades/pkg/api/raft/v1"
	"github.com/mxplusb/pleiades/pkg/fsm"
	"github.com/lni/dragonboat/v3"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewShardConfigRunner(nodeHost *dragonboat.NodeHost, logger zerolog.Logger) (*shardConfigRunner, error) {
	l := logger.With().Str("component", "shard-config").Logger()
	srv, err := newServer(logger)
	if err != nil {
		return nil, err
	}

	store, err := fsm.NewShardStore(logger)
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	runner := &shardConfigRunner{logger: l, msgServ: srv, store: store, done: done}
	client, err := runner.msgServ.GetStreamClient()
	if err != nil {
		return nil, err
	}

	streamName := strings.Split(ShardConfigStream, ".")
	_, err = client.AddStream(&nats.StreamConfig{
		Name:      streamName[len(streamName)-1],
		Subjects:  []string{ShardConfigStream},
		Retention: nats.WorkQueuePolicy,
		Discard:   nats.DiscardOld,
		Storage:   nats.MemoryStorage,
	})
	if err != nil {
		return nil, err
	}

	go runner.run()

	return runner, nil
}

type shardConfigRunner struct {
	logger  zerolog.Logger
	msgServ *server
	store   *fsm.ShardStore
	nh *dragonboat.NodeHost
	done    chan struct{}
}

func (s *shardConfigRunner) run() {
	client, err := s.msgServ.GetStreamClient()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("can't talk to nats")
	}

	listener := make(chan *nats.Msg)
	sub, err := client.ChanSubscribe(ShardConfigStream, listener)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("can't subscribe to stream")
	}
	defer sub.Unsubscribe()

	for msg := range listener {
		go func(msg *nats.Msg) {
			payload := &raftv1.ShardStateEvent{}
			err := payload.UnmarshalVT(msg.Data)
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
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				memberState, err := s.nh.SyncGetClusterMembership(ctx, payload.GetEvent().GetShardId())
				if err != nil {
					s.logger.Error().Err(err).Msg("can't get shard state")
					return
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
						for k, _ := range memberState.Removed {
							m[k] = ""
						}
						return m
					}(),
					Type: 0,
				}

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

			err = msg.Ack()
			if err != nil {
				s.logger.Error().Err(err).Msg("error acknowledging message processing")
			}
		}(msg)
	}
}
