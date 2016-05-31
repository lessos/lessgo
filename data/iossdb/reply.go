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
	ReplyOK          = "ok"
	ReplyNotFound    = "not_found"
	ReplyError       = "error"
	ReplyFail        = "fail"
	ReplyClientError = "client_error"
)

type Reply struct {
	State string
	Data  []types.Bytex
}

func (r *Reply) bytex() *types.Bytex {

	if len(r.Data) > 0 {
		return &r.Data[0]
	}

	return &types.Bytex{}
}

func (r *Reply) Bytes() []byte {
	return r.bytex().Bytes()
}

func (r *Reply) String() string {
	return r.bytex().String()
}

func (r *Reply) Int() int {
	return r.bytex().Int()
}

func (r *Reply) Int8() int8 {
	return r.bytex().Int8()
}

func (r *Reply) Int16() int16 {
	return r.bytex().Int16()
}

func (r *Reply) Int32() int32 {
	return r.bytex().Int32()
}

func (r *Reply) Int64() int64 {
	return r.bytex().Int64()
}

func (r *Reply) Uint() uint {
	return r.bytex().Uint()
}

func (r *Reply) Uint8() uint8 {
	return r.bytex().Uint8()
}

func (r *Reply) Uint16() uint16 {
	return r.bytex().Uint16()
}

func (r *Reply) Uint32() uint32 {
	return r.bytex().Uint32()
}

func (r *Reply) Uint64() uint64 {
	return r.bytex().Uint64()
}

func (r *Reply) Float32() float32 {
	return r.bytex().Float32()
}

func (r *Reply) Float64() float64 {
	return r.bytex().Float64()
}

func (r *Reply) Bool() bool {
	return r.bytex().Bool()
}

func (r *Reply) List() []types.Bytex {
	return r.Data
}

func (r *Reply) KvLen() int {
	return len(r.Data) / 2
}

func (r *Reply) KvEach(fn func(key, value types.Bytex)) int {

	if len(r.Data) < 2 {
		return 0
	}

	for i := 0; i < (len(r.Data) - 1); i += 2 {
		fn(r.Data[i], r.Data[i+1])
	}

	return r.KvLen()
}

// Json returns the map that marshals from the reply bytes as json in response .
func (r *Reply) JsonDecode(v interface{}) error {

	if len(r.Data) < 1 {
		return errors.New("json: invalid format")
	}

	return r.Data[0].JsonDecode(&v)
}
