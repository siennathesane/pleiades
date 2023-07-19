/*
 * Copyright (c) 2022-2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package eventing

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/mxplusb/pleiades/pkg/fsm"
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/mxplusb/pleiades/pkg/raftpb"
	"github.com/mxplusb/pleiades/pkg/server/runtime"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var (
	evSingleton *messaging.RaftEventHandler
)

type LifecycleManagerBuilderParams struct {
	fx.In

	StreamClient *messaging.EmbeddedMessagingStreamClient
	PubSubClient *messaging.EmbeddedMessagingPubSubClient
	ShardManager runtime.IShardManager
	RaftHost     runtime.IHost
	Logger       zerolog.Logger
}

type LifecycleManagerBuilderResults struct {
	fx.Out

	Runner *LifecycleManager
}

func NewLifecycleManager(lc fx.Lifecycle, params LifecycleManagerBuilderParams) LifecycleManagerBuilderResults {
	l := params.Logger.With().Str("component", "lifecycle-manager").Logger()

	store, err := fsm.NewSystemStore(l)
	if err != nil {
		l.Fatal().Err(err).Msg("can't create shard storage")
	}

	evSingleton = messaging.NewRaftEventHandler(params.PubSubClient, params.StreamClient, l)

	runner := &LifecycleManager{
		logger:       l,
		store:        store,
		shardManager: params.ShardManager,
		pubSubClient: params.PubSubClient,
		eventHandler: evSingleton,
		raftHost:     params.RaftHost,
	}
	runner.registerCallbacks()

	if err := runner.startHook(context.Background()); err != nil {
		l.Fatal().Err(err).Msg("can't start shards")
	}

	lc.Append(fx.Hook{
		OnStop: runner.StopHook,
	})

	return LifecycleManagerBuilderResults{Runner: runner}
}

type LifecycleManager struct {
	logger       zerolog.Logger
	store        *fsm.SystemStore
	eventHandler *messaging.RaftEventHandler
	shardManager runtime.IShardManager
	raftHost     runtime.IHost
	pubSubClient *messaging.EmbeddedMessagingPubSubClient
}

func (l *LifecycleManager) startHook(ctx context.Context) error {
	// we start the shards before we start the listener to prevent random startup issues
	err := l.StartShards()
	if err != nil {
		l.logger.Error().Err(err).Msg("can't start shards")
		return err
	}
	l.logger.Debug().Msg("started shards")

	go evSingleton.Run()
	return nil
}

func (l *LifecycleManager) StopHook(ctx context.Context) error {
	evSingleton.Stop()
	return l.StopShards()
}

// StartShards will attempt to boot any replicas that were running on the node
func (l *LifecycleManager) StartShards() error {
	all, err := l.store.GetAll()
	if err != nil {
		if !errors.Is(err, fsm.ErrNoShards) {
			l.logger.Error().Err(err).Msg("failed to get all shards")
			return err
		}
	} else {
		if len(all) == 0 {
			l.logger.Info().Msg("no existing shards found")
			return nil
		}
	}
	l.logger.Trace().Interface("shards", all).Msgf("found %d shards, will attempt to start", len(all))

	raftAddr := l.raftHost.RaftAddress()

	// nb (sienna): I hate this, rewrite it later...
	for _, shard := range all {
		for replicaId, addr := range shard.GetReplicas() {
			if addr != raftAddr {
				continue
			}

			l.logger.Trace().Interface("shard", shard).Msgf("attempting to start shard %d", shard.GetShardId())
			err = l.shardManager.StartReplica(&raftpb.StartReplicaRequest{
				ShardId:   shard.GetShardId(),
				ReplicaId: replicaId,
				Type:      shard.GetType(),
				Restart:   true,
			})
			if err != nil {
				l.logger.Error().
					Err(err).
					Uint64("shard-id", shard.GetShardId()).
					Uint64("replica-id", replicaId).
					Msg("can't start replica")
			}

			l.logger.Debug().Interface("shard", shard).Msgf("started shard %d", shard.GetShardId())
		}
		for replicaId, addr := range shard.GetObservers() {
			if addr != raftAddr {
				continue
			}

			l.logger.Trace().Interface("shard", shard).Msgf("attempting to start shard %d observer", shard.GetShardId())
			err = l.shardManager.StartReplicaObserver(&raftpb.StartReplicaObserverRequest{
				ShardId:   shard.GetShardId(),
				ReplicaId: replicaId,
				Type:      shard.GetType(),
				Restart:   true,
			})
			if err != nil {
				l.logger.Error().
					Err(err).
					Uint64("shard-id", shard.GetShardId()).
					Uint64("replica-id", replicaId).
					Msg("can't start replica observer")
			}

			l.logger.Debug().Interface("shard", shard).Msgf("started shard %d observer", shard.GetShardId())
		}
	}

	return nil
}

// StopShards will attempt to stop any replicas that were running on the node
func (l *LifecycleManager) StopShards() error {
	all, err := l.store.GetAll()
	if !errors.Is(err, fsm.ErrNoShards) {
		l.logger.Error().Err(err).Msg("failed to get all shards")
		return err
	}

	raftAddr := l.raftHost.RaftAddress()

	for _, shard := range all {
		for replicaId, addr := range shard.GetReplicas() {
			if addr != raftAddr {
				continue
			}
			_, err = l.shardManager.StopReplica(shard.GetShardId(), replicaId)
			if err != nil {
				l.logger.Error().
					Err(err).
					Uint64("shard-id", shard.GetShardId()).
					Uint64("replica-id", replicaId).
					Msg("can't start replica")
			}
		}
	}

	err = l.store.Close()
	if err != nil {
		l.logger.Error().Err(err).Msg("can't safely close system config store")
	}

	return nil
}

func (l *LifecycleManager) registerCallbacks() {
	l.eventHandler.RegisterCallback("leader-update", raftpb.Event_EVENT_LEADER_UPDATED, l.handleLeaderUpdate)
}

func (l *LifecycleManager) handleLeaderUpdate(event *raftpb.RaftEvent) {
	l.logger.Debug().Interface("payload", event).Msg("leader update recieved")

	// safety check
	if event.Typ != raftpb.EventType_EVENT_TYPE_RAFT {
		l.logger.Error().Msg("event type mismatched")
		return
	}

	lu := event.GetLeaderUpdate()
	if lu == nil {
		l.logger.Error().Msg("leader event is empty")
		return
	}

	members, err := l.shardManager.GetShardMembers(lu.GetShardId())
	if err != nil {
		l.logger.Error().Err(err).Uint64("shard-id", lu.GetShardId()).Msg("can't get shard members")
		return
	}

	update := &raftpb.ShardState{
		LastUpdated:    event.GetTimestamp(),
		ShardId:        lu.GetShardId(),
		ConfigChangeId: members.ConfigChangeId,
		Replicas:       members.Replicas,
		Observers:      members.Observers,
		Witnesses:      members.Witnesses,
		Removed: func() map[uint64]string {
			m := make(map[uint64]string)
			for k := range members.Removed {
				m[k] = ""
			}
			return m
		}(),
		Type: raftpb.StateMachineType_STATE_MACHINE_TYPE_KV,
	}

	err = l.store.Put(update)
	if err != nil {
		l.logger.Error().Err(err).Msg("can't put the update")
	}

	return
}
