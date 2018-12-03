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
	ar_int_mu sync.Mutex
)

type ArrayInt []int

func (ar *ArrayInt) Has(ui int) bool {

	for _, v := range *ar {

		if v == ui {
			return true
		}
	}

	return false
}

func (ar *ArrayInt) Set(ui int) bool {

	ar_int_mu.Lock()
	defer ar_int_mu.Unlock()

	for _, v := range *ar {

		if v == ui {
			return false
		}
	}

	*ar = append(*ar, ui)

	return true
}

func (ar *ArrayInt) Del(ui int) {

	ar_int_mu.Lock()
	defer ar_int_mu.Unlock()

	for i, v := range *ar {

		if v == ui {
			*ar = append((*ar)[:i], (*ar)[i+1:]...)
			return
		}
	}
}

func (ar *ArrayInt) Equal(ar2 ArrayInt) bool {

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

func (ar *ArrayInt) MatchAny(ar2 ArrayInt) bool {

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
func (ar ArrayInt) Len() int {
	return len(ar)
}

func (ar ArrayInt) Less(i, j int) bool {
	return ar[i] < ar[j]
}

func (ar ArrayInt) Swap(i, j int) {
	ar[i], ar[j] = ar[j], ar[i]
}
