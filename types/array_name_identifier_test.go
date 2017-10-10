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

func TestArrayNameIdentifier(t *testing.T) {

	ar := ArrayNameIdentifier([]NameIdentifier{"ni/aa", "ni/bb"})

	if ar.Has("ni/aa") == false ||
		ar.Has("ni/cc") == true {
		t.Fatal("Failed on Has")
	}

	ar2 := ArrayNameIdentifier([]NameIdentifier{"ni/aa", "ni/bb", "ni/cc"})
	if ar2.Has("ni/cc") != true {
		t.Fatal("Failed on Set")
	}

	ar2.Set("ni/dd")
	ar2.Set("ni/ee")

	if ar.Equal(ar) != true {
		t.Fatal("Failed on Equal")
	}

	if ar.Equal(ar2) != false {
		t.Fatal("Failed on Equal")
	}

	ar2.Del("ni/aa")
	if ar2.Equal(ArrayNameIdentifier([]NameIdentifier{"ni/bb", "ni/cc", "ni/dd", "ni/ee"})) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del("ni/ee")
	if ar2.Equal(ArrayNameIdentifier([]NameIdentifier{"ni/bb", "ni/cc", "ni/dd"})) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del("ni/cc")
	if ar2.Equal(ArrayNameIdentifier([]NameIdentifier{"ni/bb", "ni/dd"})) != true {
		t.Fatal("Failed on Del")
	}
}
