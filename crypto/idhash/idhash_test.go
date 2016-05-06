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
	"testing"
)

func TestMain(t *testing.T) {

	if len(Rand(16)) != 16 {
		t.Fatal("Failed on Rand")
	}

	if len(RandToHexString(16)) != 16 {
		t.Fatal("Failed on RandToHexString")
	}

	if HashToHexString("123456", 16) != "e10adc3949ba59ab" {
		t.Fatal("Failed on HashToHexString")
	}

	if len(RandToBase64String(0)) != 4 {
		t.Fatal("Failed on RandToBase64String")
	}

	if len(RandToBase64String(16)) != 16 {
		t.Fatal("Failed on RandToBase64String")
	}

	if len(RandUUID()) != 36 {
		t.Fatal("Failed on RandUUID")
	}
}
