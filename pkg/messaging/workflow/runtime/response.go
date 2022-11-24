/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package runtime

type Response struct {
	RequestID string
	Header    map[string][]string
	Body      []byte
}

func (response *Response) SetHeader(header string, value string) {
	response.Header[header] = []string{value}
}
