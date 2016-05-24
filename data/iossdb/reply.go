// Copyright 2013-2016 lessgo Author, All rights reserved.
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

package iossdb

import (
	"errors"

	"github.com/lessos/lessgo/types"
)

const (
	ReplyOK          string = "ok"
	ReplyNotFound    string = "not_found"
	ReplyError       string = "error"
	ReplyFail        string = "fail"
	ReplyClientError string = "client_error"
)

type Reply struct {
	State string
	Data  []types.Bytex
}

type Entry struct {
	Key   types.Bytex
	Value types.Bytex
}

func (r *Reply) Bytes() []byte {

	if len(r.Data) == 0 {
		return []byte{}
	}

	return r.Data[0].Bytes()
}

func (r *Reply) String() string {

	if len(r.Data) > 0 {
		return r.Data[0].String()
	}

	return ""
}

func (r *Reply) Int() int {
	return int(r.Int64())
}

func (r *Reply) Int64() int64 {

	if len(r.Data) < 1 {
		return 0
	}

	return r.Data[0].Int64()
}

func (r *Reply) Uint() uint {
	return uint(r.Uint64())
}

func (r *Reply) Uint64() uint64 {

	if len(r.Data) < 1 {
		return 0
	}

	return r.Data[0].Uint64()
}

func (r *Reply) Float64() float64 {

	if len(r.Data) < 1 {
		return 0
	}

	return r.Data[0].Float64()
}

func (r *Reply) Bool() bool {

	if len(r.Data) < 1 {
		return false
	}

	return r.Data[0].Bool()
}

func (r *Reply) List() []types.Bytex {
	return r.Data
}

func (r *Reply) Hash() []Entry {

	hs := []Entry{}

	if len(r.Data) < 2 {
		return hs
	}

	for i := 0; i < (len(r.Data) - 1); i += 2 {
		hs = append(hs, Entry{
			Key:   r.Data[i],
			Value: r.Data[i+1],
		})
	}

	return hs
}

func (r *Reply) Each(fn func(key, value types.Bytex)) int {

	if len(r.Data) < 2 {
		return 0
	}

	for i := 0; i < (len(r.Data) - 1); i += 2 {
		fn(r.Data[i], r.Data[i+1])
	}

	return len(r.Data) / 2
}

// Json returns the map that marshals from the reply bytes as json in response .
func (r *Reply) JsonDecode(v interface{}) error {

	if len(r.Data) < 1 {
		return errors.New("json: invalid format")
	}

	return r.Data[0].JsonDecode(&v)
}

func (r *Entry) JsonDecode(v interface{}) error {
	return r.JsonDecode(&v)
}
