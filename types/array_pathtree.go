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
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	ar_pt_mu sync.Mutex
)

type ArrayPathTree []string

func (ar *ArrayPathTree) Has(s string) bool {

	for _, v := range *ar {

		if v == s {
			return true
		}
	}

	return false
}

func (ar *ArrayPathTree) Set(s string) {

	ar_pt_mu.Lock()
	defer ar_pt_mu.Unlock()

	ps := strings.Split(strings.Trim(s, ","), ",")

	for _, p := range ps {

		var (
			pns    = strings.Split(strings.Trim(filepath.Clean(p), "/"), "/")
			pnbase = ""
		)

		for i, pn := range pns {

			pn = strings.TrimSpace(pn)
			pnps := []string{}

			if i+1 == len(pns) {
				pnps = strings.Split(pn, "|")
			} else {
				pnps = []string{pn}
			}

			for _, pnp := range pnps {

				pnp = strings.TrimSpace(pnp)

				if pnbase != "" {
					pnp = pnbase + "/" + pnp
				}

				hit := false

				for _, v := range *ar {

					if v == pnp {
						hit = true
						break
					}
				}

				if !hit {
					*ar = append(*ar, pnp)
				}
			}

			if pnbase == "" {
				pnbase = pn
			} else {
				pnbase += "/" + pn
			}
		}
	}
}

func (ar *ArrayPathTree) Equal(ar2 ArrayPathTree) bool {

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

func (ar *ArrayPathTree) Sort() {

	ar_pt_mu.Lock()
	defer ar_pt_mu.Unlock()

	sort.Strings(*ar)
}
