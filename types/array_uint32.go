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
	"sync"
)

var (
	ar_uint32_mu sync.Mutex
)

type ArrayUint32 []uint32

func (ar *ArrayUint32) Has(ui uint32) bool {

	for _, v := range *ar {

		if v == ui {
			return true
		}
	}

	return false
}

func (ar *ArrayUint32) Set(ui uint32) {

	ar_uint32_mu.Lock()
	defer ar_uint32_mu.Unlock()

	for _, v := range *ar {

		if v == ui {
			return
		}
	}

	*ar = append(*ar, ui)
}

func (ar *ArrayUint32) Del(ui uint32) {

	ar_uint32_mu.Lock()
	defer ar_uint32_mu.Unlock()

	for i, v := range *ar {

		if v == ui {
			*ar = append((*ar)[:i], (*ar)[i+1:]...)
			return
		}
	}
}

func (ar *ArrayUint32) Equal(ar2 ArrayUint32) bool {

	if len(*ar) != len(ar2) {
		return false
	}

	for _, v := range *ar {

		hit := false

		for _, v2 := range ar2 {

			if v == v2 {
				hit = true
				break
			}
		}

		if !hit {
			return false
		}
	}

	return true
}

// sort.Interface
func (ar ArrayUint32) Len() int {
	return len(ar)
}

func (ar ArrayUint32) Less(i, j int) bool {
	return ar[i] < ar[j]
}

func (ar ArrayUint32) Swap(i, j int) {
	ar[i], ar[j] = ar[j], ar[i]
}
