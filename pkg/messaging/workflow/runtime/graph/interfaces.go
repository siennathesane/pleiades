/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package graph

type IOperation interface {
	GetId() string
	Encode() []byte
	GetProperties() map[string][]string
	// Execute executes an operation, executor can pass configuration
	Execute([]byte, map[string]interface{}) ([]byte, error)
}
