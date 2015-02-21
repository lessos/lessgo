// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpsrv

import (
	"net/http"
)

type Response struct {
	Status      int
	ContentType string
	Out         http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{Out: w}
}

func (resp *Response) WriteHeader(status int, ctype string) {

	if resp.Status == 0 {
		resp.Status = status
		resp.Out.WriteHeader(resp.Status)
	}

	if resp.ContentType == "" {
		resp.ContentType = ctype
		resp.Out.Header().Set("Content-Type", resp.ContentType)
	}
}
