/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package services

// IServiceManager defines a basic way to interface with services.
type IServiceManager interface {
	Start(retry bool) error
	Stop(retry bool) error
	Restart(retry bool) error
}

type IStore[T any] interface {
	Get(key string) (*T, error)
	GetAll() (map[string]*T, error)
	Put(key string, payload *T) error
}
