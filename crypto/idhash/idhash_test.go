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

package idhash

import (
	"bytes"
	"testing"
)

func TestMain(t *testing.T) {

	if len(Rand(16)) != 16 {
		t.Fatal("Failed on Rand")
	}

	if bytes.Compare(Hash([]byte("123456"), 8), []byte{225, 10, 220, 57, 73, 186, 89, 171}) != 0 {
		t.Fatal("Failed on Hash")
	}

	if len(RandHexString(16)) != 16 {
		t.Fatal("Failed on RandHexString")
	}

	if HashToHexString([]byte("123456"), 16) != "e10adc3949ba59ab" {
		t.Fatal("Failed on HashStringToHexString")
	}

	if len(RandBase64String(0)) != 4 {
		t.Fatal("Failed on RandBase64String")
	}

	if len(RandBase64String(16)) != 16 {
		t.Fatal("Failed on RandBase64String")
	}

	if len(RandUUID()) != 36 {
		t.Fatal("Failed on RandUUID")
	}
}

func Benchmark_Rand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Rand(16)
	}
}

func Benchmark_Hash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash([]byte("123456"), 16)
	}
}

func Benchmark_RandHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandHexString(16)
	}
}

func Benchmark_HashToHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashToHexString([]byte("123456"), 16)
	}
}

func Benchmark_RandBase64String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandBase64String(16)
	}
}

func Benchmark_RandUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandUUID()
	}
}
