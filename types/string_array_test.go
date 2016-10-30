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

func TestStringArray(t *testing.T) {

	ar := StringArray("aa,bb")

	if ar.Has("aa") == false ||
		ar.Has("cc") == true {
		t.Fatal("Failed on Has")
	}

	ar2 := StringArray("aa,bb,cc")
	if ar2.Has("cc") != true {
		t.Fatal("Failed on Set")
	}

	ar2.Set("dd")
	ar2.Set("ee")

	if ar.Equal(ar) != true {
		t.Fatal("Failed on Equal")
	}

	if ar.Equal(ar2) != false {
		t.Fatal("Failed on Equal")
	}

	ar2.Del("aa")
	if ar2.Equal(StringArray("bb,cc,dd,ee")) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del("ee")
	if ar2.Equal(StringArray("bb,cc,dd")) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del("cc")
	if ar2.Equal(StringArray("bb,dd")) != true {
		t.Fatal("Failed on Del")
	}
}
