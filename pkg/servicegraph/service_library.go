/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package servicegraph

type IServiceLibrary interface {
	AddService(svc Service) error
	AddServices(svcs []Service) error
	GetService(svc Service) (Service, error)
	StartService(svc Service, retry bool) error
	StopService(retry, force bool, svc Service) error
}
