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
	"github.com/rs/zerolog"
)

var (
	serverSingleton *server
)

func newServer(logger zerolog.Logger) (*server, error) {
	if serverSingleton != nil {
		return serverSingleton, nil
	}

	srv, err := messaging.NewEmbeddedMessagingWithDefaults()
	if err != nil {
		return nil, err
	}

	serverSingleton = &server{srv, logger.With().Str("component", "eventing").Logger()}

	return serverSingleton, nil
}

type server struct {
	*messaging.EmbeddedMessaging
	logger zerolog.Logger
}
