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

package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Universal Bytes
type Bytex []byte

func (bx Bytex) Bytes() []byte {
	return bx
}

func (bx Bytex) String() string {
	return string(bx)
}

func (bx Bytex) Bool() bool {

	if len(bx) < 1 {
		return false
	}

	if b, err := strconv.ParseBool(string(bx)); err == nil {
		return b
	}

	return false
}

// int
func (bx Bytex) Int() int {
	return int(bx.Int64())
}

func (bx Bytex) Int8() int8 {
	return int8(bx.Int64())
}

func (bx Bytex) Int16() int16 {
	return int16(bx.Int64())
}

func (bx Bytex) Int32() int32 {
	return int32(bx.Int64())
}

func (bx Bytex) Int64() int64 {

	if len(bx) < 1 {
		return 0
	}

	if i64, err := strconv.ParseInt(string(bx), 10, 64); err == nil {
		return i64
	}

	return 0
}

// unsigned int
func (bx Bytex) Uint() uint {
	return uint(bx.Uint64())
}

func (bx Bytex) Uint8() uint8 {
	return uint8(bx.Uint64())
}

func (bx Bytex) Uint16() uint16 {
	return uint16(bx.Uint64())
}

func (bx Bytex) Uint32() uint32 {
	return uint32(bx.Uint64())
}

func (bx Bytex) Uint64() uint64 {

	if len(bx) < 1 {
		return 0
	}

	if i64, err := strconv.ParseUint(string(bx), 10, 64); err == nil {
		return i64
	}

	return 0
}

// float
func (bx Bytex) Float32() float32 {
	return float32(bx.Float64())
}

func (bx Bytex) Float64() float64 {

	if len(bx) < 1 {
		return 0
	}

	if f64, err := strconv.ParseFloat(string(bx), 64); err == nil {
		return f64
	}

	return 0
}

func (bx Bytex) JsonDecode(v interface{}) error {

	if len(bx) < 2 {
		fmt.Errorf("json: invalid format")
	}

	return json.Unmarshal(bx, &v)
}
