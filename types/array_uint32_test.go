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
	"sort"
	"testing"
)

func TestArrayUint32(t *testing.T) {

	ar := ArrayUint32([]uint32{1, 2})

	if ar.Has(1) == false ||
		ar.Has(3) == true {
		t.Fatal("Failed on Has")
	}

	ar2 := ArrayUint32([]uint32{1, 2, 3})
	if ar2.Has(3) != true {
		t.Fatal("Failed on Set")
	}

	ar2.Set(4)
	ar2.Set(5)

	if ar.Equal(ar) != true {
		t.Fatal("Failed on Equal")
	}

	if ar.Equal(ar2) != false {
		t.Fatal("Failed on Equal")
	}

	ar2.Del(1)
	if ar2.Equal(ArrayUint32([]uint32{2, 3, 4, 5})) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del(5)
	if ar2.Equal(ArrayUint32([]uint32{2, 3, 4})) != true {
		t.Fatal("Failed on Del")
	}

	ar2.Del(3)
	if ar2.Equal(ArrayUint32([]uint32{2, 4})) != true {
		t.Fatal("Failed on Del")
	}
}

func TestArrayUint32Sort(t *testing.T) {

	ar := ArrayUint32([]uint32{1, 3, 2})

	sort.Sort(ar)

	if ar[0] != 1 || ar[1] != 2 || ar[2] != 3 {
		t.Fatal("Failed on sort.Interface")
	}
}
