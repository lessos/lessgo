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

package json

import (
	"os"
	"testing"
)

type Object struct {
	Name string `json:"name"`
}

var (
	js_i0   = []byte(`{"name":"value-of-string"}`)
	js_i1   = "{\n\t\"name\": \"value-of-string\"\n}"
	js_i2   = "{\n\t\t\"name\": \"value-of-string\"\n}"
	js_file = "/tmp/output.file.json"
)

func TestJson(t *testing.T) {

	var obj Object
	if err := Decode(js_i0, &obj); err != nil || obj.Name != "value-of-string" {
		t.Fatal("Failed on Decode")
	}

	if bs, err := Encode(obj, "\t"); err != nil || string(bs) != js_i1 {
		t.Fatal("Failed on Encode")
	}

	if bs2, err := Indent(js_i0, "\t\t"); err != nil || string(bs2) != js_i2 {
		t.Fatal("Failed on Indent")
	}

	if fp, err := os.OpenFile(js_file, os.O_RDWR|os.O_CREATE, 0644); err == nil {

		fp.Close()

		if err := EncodeToFile(obj, js_file, "\t"); err != nil {
			t.Fatal("Faile on EncodeToFile")
		}

		var obj2 Object
		if err := DecodeFile(js_file, &obj2); err != nil || obj2.Name != "value-of-string" {
			t.Fatal("Faile on DecodeFile")
		}

		os.Remove(js_file)
	}
}

func Benchmark_Decode(b *testing.B) {

	var obj Object

	for i := 0; i < b.N; i++ {
		Decode(js_i0, &obj)
	}
}

func Benchmark_Encode(b *testing.B) {

	obj := Object{
		Name: "value-of-string",
	}

	for i := 0; i < b.N; i++ {
		Encode(obj, "")
	}
}

func Benchmark_Indent(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Indent(js_i0, "\t")
	}
}
