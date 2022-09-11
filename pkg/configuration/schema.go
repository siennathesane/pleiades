/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package configuration

type Configuration struct {
	Datastore *Storage `json:"datastore,omitempty" yaml:"datastore,omitempty"`
}

type Storage struct {
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`

}
