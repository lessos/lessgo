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

func TestArray(t *testing.T) {

	ar := ArrayString([]string{"aa", "bb"})

	if ar.Contain("aa") == false ||
		ar.Contain("cc") == true {
		t.Fatal("Failed on Contain")
	}

	ar2 := ar.Insert("cc")
	if ar2.Contain("cc") != true {
		t.Fatal("Failed on Insert")
	}

	ar2 = ar2.Insert("dd")
	ar2 = ar2.Insert("ee")

	if ar.Equal(ar) != true {
		t.Fatal("Failed on Equal")
	}

	if ar.Equal(ar2) != false {
		t.Fatal("Failed on Equal")
	}

	ar2 = ar2.Remove("aa")
	if ar2.Equal(ArrayString([]string{"bb", "cc", "dd", "ee"})) != true {
		t.Fatal("Failed on Remove")
	}

	ar2 = ar2.Remove("ee")
	if ar2.Equal(ArrayString([]string{"bb", "cc", "dd"})) != true {
		t.Fatal("Failed on Remove")
	}

	ar2 = ar2.Remove("cc")
	if ar2.Equal(ArrayString([]string{"bb", "dd"})) != true {
		t.Fatal("Failed on Remove")
	}
}
