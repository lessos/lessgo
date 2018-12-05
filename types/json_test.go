// Copyright 2013-2018 lessgo Author, All rights reserved.
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

	"github.com/lessos/lessgo/encoding/json"
)

const (
	jsonTest01 = `{
  "b_bool": true,
  "a_item": {
    "b": "vb",
    "a": "va"
  },
  "c_str": "vstr",
  "b_int": 1,
  "d_items": [
    {
      "b": "vb",
      "a": "va"
    }
  ],
  "b_float": 0.1
}`
	jsonTest01Sorted = `{
  "a_item": {
    "a": "va",
    "b": "vb"
  },
  "b_bool": true,
  "b_float": 0.1,
  "b_int": 1,
  "c_str": "vstr",
  "d_items": [
    {
      "a": "va",
      "b": "vb"
    }
  ]
}`
)

type jsonRawTestSub struct {
	A string `json:"a"`
	B string `json:"b"`
}

type jsonRawTest struct {
	AItem  jsonRawTestSub   `json:"a_item"`
	BBool  bool             `json:"b_bool"`
	BFloat float64          `json:"b_float"`
	BInt   int              `json:"b_int"`
	CStr   string           `json:"c_str"`
	DItems []jsonRawTestSub `json:"d_items"`
}

func TestJsonTypeless(t *testing.T) {

	//
	item, err := JsonTypelessItemDecode([]byte(jsonTest01))
	if err != nil {
		t.Fatal("Failed on JsonTypeless/Decode")
	}

	//
	item.Sort()
	if js, err := item.Encode("  "); err != nil {
		t.Fatal("Failed on JsonTypeless/Encode")
	} else if string(js) != jsonTest01Sorted {
		t.Fatal("Failed on JsonTypeless/Encode Diff")
	}

	if field := item.Get("c_str"); field == nil || field.String() != "vstr" {
		t.Fatal("Failed on JsonTypeless/Field.Get.String")
	}

	if field := item.Get("b_int"); field == nil || field.Int64() != 1 {
		t.Fatal("Failed on JsonTypeless/Field.Get.Int64")
	}

	if field := item.Get("b_bool"); field == nil || field.Bool() != true {
		t.Fatal("Failed on JsonTypeless/Field.Get.Bool")
	}

	if field := item.Get("b_float"); field == nil || field.Float64() != 0.1 {
		t.Fatal("Failed on JsonTypeless/Field.Get.Float64")
	}

	if field := item.Get("a_item"); field == nil {
		t.Fatal("Failed on JsonTypeless/Field.Get.Item")
	} else {
		item2 := field.Item()
		if field := item2.Get("a"); field == nil || field.String() != "va" {
			t.Fatal("Failed on JsonTypeless/Field.Get.Item.String")
		}
	}

	if field := item.Get("d_items"); field == nil || !field.IsArray() {
		t.Fatal("Failed on JsonTypeless/Field.Get.Items")
	} else {
		item2 := field.Array()
		if len(item2) < 1 {
			t.Fatal("Failed on JsonTypeless/Field.Get.Items LEN")
		}
		if field := item2[0].Get("a"); field == nil || field.String() != "va" {
			t.Fatal("Failed on JsonTypeless/Field.Get.Items[0].String")
		}
	}
}

func Benchmark_JsonTypeless_Encode(b *testing.B) {

	//
	item, err := JsonTypelessItemDecode([]byte(jsonTest01))
	if err != nil {
		b.Fatal("Failed on JsonTypeless/Decode")
	}

	for i := 0; i < b.N; i++ {
		if _, err := json.Encode(item, "  "); err != nil {
			b.Fatal("Failed on JsonTypeless/Encode")
		}
	}
}

func Benchmark_JsonTypeless_EncodeRaw(b *testing.B) {

	//
	var item jsonRawTest
	if err := json.Decode([]byte(jsonTest01), &item); err != nil {
		b.Fatal("Failed on JsonTypeless/Decode/Raw")
	}

	for i := 0; i < b.N; i++ {
		if _, err := json.Encode(item, "  "); err != nil {
			b.Fatal("Failed on JsonTypeless/Encode")
		}
	}
}

func Benchmark_JsonTypeless_Decode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var item JsonTypelessItem
		if err := json.Decode([]byte(jsonTest01), &item); err != nil {
			b.Fatal("Failed on JsonTypeless/Decode")
		}
	}
}

func Benchmark_JsonTypeless_DecodeRaw(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var item jsonRawTest
		if err := json.Decode([]byte(jsonTest01), &item); err != nil {
			b.Fatal("Failed on JsonTypeless/Decode/Raw")
		}
	}
}
