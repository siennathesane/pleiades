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
	"github.com/mxplusb/pleiades/pkg/messaging"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

const (
	SystemStreamName = "system"
)

var (
	serverSingleton *Server
)

func NewServer(logger zerolog.Logger) (*Server, error) {
	if serverSingleton != nil {
		return serverSingleton, nil
	}

	srv, err := messaging.NewEmbeddedMessagingWithDefaults()
	if err != nil {
		return nil, err
	}

	serverSingleton = &Server{srv, logger.With().Str("component", "eventing").Logger()}

	client, err := serverSingleton.GetStreamClient()
	if err != nil {
		return nil, err
	}

	_, err = client.AddStream(&nats.StreamConfig{
		Name: SystemStreamName,
		Description: "All internal system streams",
		Subjects: []string{"system.>"},
		Retention: nats.WorkQueuePolicy,
		Discard:   nats.DiscardOld,
		Storage:   nats.MemoryStorage,
	})
	if err != nil {
		return nil, err
	}

	return serverSingleton, nil
}

type Server struct {
	*messaging.EmbeddedMessaging
	logger zerolog.Logger
}
