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
	ar_ni_mu sync.Mutex
)

type ArrayNameIdentifier []NameIdentifier

func (ar *ArrayNameIdentifier) Has(s NameIdentifier) bool {

	for _, v := range *ar {

		if v == s {
			return true
		}
	}

	return false
}

func (ar *ArrayNameIdentifier) Set(s NameIdentifier) error {

	if err := s.Valid(); err != nil {
		return err
	}

	ar_ni_mu.Lock()
	defer ar_ni_mu.Unlock()

	for _, v := range *ar {

		if v == s {
			return nil
		}
	}

	*ar = append(*ar, s)
	return nil
}

func (ar *ArrayNameIdentifier) Del(s NameIdentifier) {

	ar_ni_mu.Lock()
	defer ar_ni_mu.Unlock()

	for i, v := range *ar {

		if v == s {
			*ar = append((*ar)[:i], (*ar)[i+1:]...)
			return
		}
	}
}

func (ar *ArrayNameIdentifier) Equal(ar2 ArrayNameIdentifier) bool {

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
