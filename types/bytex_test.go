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
	"testing"
)

type Object struct {
	Name string `json:"name"`
}

var (
	js_i0 = []byte(`{"name":"value-of-string"}`)
)

func TestBytex(t *testing.T) {

	bx := Bytex("123")

	if string(bx.Bytes()) != "123" {
		t.Fatal("Failed on Bytes")
	}

	if bx.String() != "123" {
		t.Fatal("Failed on String")
	}

	if bx.Int() != int(123) {
		t.Fatal("Failed on Int")
	}

	if bx.Int8() != int8(123) {
		t.Fatal("Failed on Int8")
	}
	if bx.Int16() != int16(123) {
		t.Fatal("Failed on Int16")
	}

	if bx.Int32() != int32(123) {
		t.Fatal("Failed on Int32")
	}

	if bx.Int64() != int64(123) {
		t.Fatal("Failed on Int64")
	}

	if bx.Uint() != uint(123) {
		t.Fatal("Failed on Uint")
	}

	if bx.Uint8() != uint8(123) {
		t.Fatal("Failed on Uint8")
	}

	if bx.Uint16() != uint16(123) {
		t.Fatal("Failed on Uint16")
	}

	if bx.Uint32() != uint32(123) {
		t.Fatal("Failed on Uint32")
	}

	if bx.Uint64() != uint64(123) {
		t.Fatal("Failed on Uint64")
	}

	if bx.Float32() != float32(123) {
		t.Fatal("Failed on Float32")
	}

	if bx.Float64() != float64(123) {
		t.Fatal("Failed on Float64")
	}

	if Bytex([]byte("true")).Bool() != true ||
		Bytex([]byte("false")).Bool() != false ||
		Bytex([]byte("1")).Bool() != true ||
		Bytex([]byte("0")).Bool() != false {
		t.Fatal("Failed on Bool")
	}

	var obj Object
	if err := Bytex(js_i0).JsonDecode(&obj); err != nil || obj.Name != "value-of-string" {
		t.Fatal("Failed on JsonDecode")
	}
}

func Benchmark_Bytex_JsonDecode(b *testing.B) {

	var obj Object
	bx := Bytex(js_i0)

	for i := 0; i < b.N; i++ {
		bx.JsonDecode(&obj)
	}
}
