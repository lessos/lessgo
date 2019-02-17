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
	"reflect"
	"sync"
)

var iter_object_mu sync.RWMutex

type IterObject interface {
	IterKey() string
	// IterEqual(it IterObject) bool
}

/**
func IterObjectsEqual(ls, ls2 interface{}) bool {

	vs := reflect.ValueOf(ls)
	vs2 := reflect.ValueOf(ls2)
	if vs.Len() != vs2.Len() {
		return false
	}

	hit := false
	for i := 0; i < vs.Len(); i++ {
		o, ok := (vs.Index(i).Interface()).(IterObject)
		if !ok {
			return false
		}
		hit = false
		for j := 0; j < vs2.Len(); j++ {
			o2, ok2 := (vs.Index(i).Interface()).(IterObject)
			if !ok2 {
				return false
			}

			if o.IterKey() != o2.IterKey() {
				continue
			}

			if !o.IterEqual(o2) {
				return false
			}
			hit = true
		}
		if !hit {
			return false
		}
	}

	return true
}
*/

func IterObjectLookup(items interface{}, key string, fn func(i int)) {

	iter_object_mu.Lock()
	defer iter_object_mu.Unlock()

	vs := reflect.ValueOf(items)
	for i := 0; i < vs.Len(); i++ {
		o, ok := (vs.Index(i).Interface()).(IterObject)
		if !ok {
			break
		}
		if o.IterKey() == key {
			fn(i)
			return
		}
	}

	fn(-1)
}

func IterObjectGet(items interface{}, key string) IterObject {

	iter_object_mu.RLock()
	defer iter_object_mu.RUnlock()

	vs := reflect.ValueOf(items)

	for i := 0; i < vs.Len(); i++ {
		o, ok := (vs.Index(i).Interface()).(IterObject)
		if !ok {
			return nil
		}
		if o.IterKey() == key {
			return o
		}
	}

	return nil
}
