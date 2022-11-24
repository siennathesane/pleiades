/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

type Request struct {
	FlowName  string
	RequestID string
	Header    map[string][]string
	RawQuery  string
	Query     map[string][]string
	Body      []byte
}

func (request *Request) GetHeader(header string) string {
	val := request.Header[header]
	if len(val) >= 1 {
		return val[0]
	}
	return ""
}
