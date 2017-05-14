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

package hissdb

import (
	"encoding/json"
	"errors"
	"strconv"
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
	Data  []string
}

type Entry struct {
	Key   string
	Value string
}

func (r *Reply) String() string {

	if len(r.Data) > 0 {
		return r.Data[0]
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

	if i64, err := strconv.ParseInt(r.Data[0], 10, 64); err == nil {
		return i64
	}

	return 0
}

func (r *Reply) Uint() uint {
	return uint(r.Uint64())
}

func (r *Reply) Uint64() uint64 {

	if len(r.Data) < 1 {
		return 0
	}

	if i64, err := strconv.ParseUint(r.Data[0], 10, 64); err == nil {
		return i64
	}

	return 0
}

func (r *Reply) Float64() float64 {

	if len(r.Data) < 1 {
		return 0
	}

	if f64, err := strconv.ParseFloat(r.Data[0], 64); err == nil {
		return f64
	}

	return 0
}

func (r *Reply) Bool() bool {

	if len(r.Data) < 1 {
		return false
	}

	if b, err := strconv.ParseBool(r.Data[0]); err == nil {
		return b
	}

	return false
}

func (r *Reply) List() []string {

	if len(r.Data) > 0 {
		return r.Data
	}

	return []string{}
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

// HMap returns a slice string reply  as a map[string]string or an error.
// it must have an even number of elements,
// they must be in a "key value key value..." order and
// keys must all be set, cannot be set empty strings.
// Nil values are returned as empty strings.
// Useful for hash commands.
func (r *Reply) HMap() (map[string]string, error) {
	if r.State != ReplyOK {
		return nil, errors.New(r.State)
	}
	if len(r.Data)%2 != 0 {
		return nil, errors.New("reply has odd number of elements")
	}
	rmap := map[string]string{}
	for i := 0; i < (len(r.Data) - 1); i += 2 {
		if 0 == len(r.Data[i]) {
			return nil, errors.New("reply map key has been set empty string")
		}
		if _, exist := rmap[r.Data[i]]; !exist {
			rmap[r.Data[i]] = r.Data[i+1]
		} else {
			return nil, errors.New("reply map key repeated")
		}
	}
	return rmap, nil
}

// to be remove
func (r *Reply) Json(v interface{}) error {
	return r.JsonDecode(v)
}

// to be remove
func (r *Entry) Json(v interface{}) error {
	return r.JsonDecode(v)
}

// Json returns the map that marshals from the reply bytes as json in response .
func (r *Reply) JsonDecode(v interface{}) error {
	return _json_decode(r.String(), v)
}

func (r *Entry) JsonDecode(v interface{}) error {
	return _json_decode(r.Value, v)
}

func _json_decode(src string, rs interface{}) (err error) {

	if len(src) < 2 {
		return errors.New("json: invalid format")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("json: invalid format")
		}
	}()

	return json.Unmarshal([]byte(src), &rs)
}
