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
	ar_uint16_mu sync.Mutex
)

type ArrayUint16 []uint16

func (ar *ArrayUint16) Has(ui uint16) bool {

	for _, v := range *ar {

		if v == ui {
			return true
		}
	}

	return false
}

func (ar *ArrayUint16) Set(ui uint16) bool {

	ar_uint16_mu.Lock()
	defer ar_uint16_mu.Unlock()

	for _, v := range *ar {

		if v == ui {
			return false
		}
	}

	*ar = append(*ar, ui)

	return true
}

func (ar *ArrayUint16) Del(ui uint16) {

	ar_uint16_mu.Lock()
	defer ar_uint16_mu.Unlock()

	for i, v := range *ar {

		if v == ui {
			*ar = append((*ar)[:i], (*ar)[i+1:]...)
			return
		}
	}
}

func (ar *ArrayUint16) Equal(ar2 ArrayUint16) bool {

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

func (ar *ArrayUint16) MatchAny(ar2 ArrayUint16) bool {

	for _, v2 := range ar2 {

		for _, v := range *ar {

			if v == v2 {
				return true
			}
		}
	}

	return false
}

// sort.Interface
func (ar ArrayUint16) Len() int {
	return len(ar)
}

func (ar ArrayUint16) Less(i, j int) bool {
	return ar[i] < ar[j]
}

func (ar ArrayUint16) Swap(i, j int) {
	ar[i], ar[j] = ar[j], ar[i]
}
