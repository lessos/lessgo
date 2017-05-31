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
	"net/url"
	"strconv"
)

type Params struct {
	url.Values // A unified view of all the individual param maps below

	// Set by the ParamsFilter
	Query url.Values // Parameters from the query string, e.g. /index?limit=10
	Form  url.Values // Parameters from the request body.
}

func newParams() *Params {

	return &Params{
		Values: make(url.Values, 0),
	}
}

func ParamsFilter(c *Controller) {

	c.Params.Query = c.Request.URL.Query()
	for k, v := range c.Params.Query {
		if _, ok := c.Params.Values[k]; !ok {
			c.Params.Values[k] = v
		}
	}

	if c.Request.ContentType == "application/x-www-form-urlencoded" {
		// Typical form.
		if err := c.Request.ParseForm(); err != nil {
			// Error parsing request body
		} else {

			c.Params.Form = c.Request.Form

			for k, v := range c.Params.Form {
				if _, ok := c.Params.Values[k]; !ok {
					c.Params.Values[k] = v
				}
			}
		}
	}
}

func (p *Params) String(key string) string {
	return p.Values.Get(key)
}

func (p *Params) Uint64(key string) uint64 {
	ui64, _ := strconv.ParseUint(p.Values.Get(key), 10, 64)
	return ui64
}

func (p *Params) Int64(key string) int64 {
	i64, _ := strconv.ParseInt(p.Values.Get(key), 10, 64)
	return i64
}
