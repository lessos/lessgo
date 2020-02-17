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
	"errors"
	"fmt"
	"sync"
)

var (
	kvp_mu sync.RWMutex

	kvp_ErrKeyEmpty   = errors.New("key cannot be empty")
	kvp_ErrKeyLength  = errors.New("length of the key must be less than 100")
	kvp_ErrKeyInvalid = errors.New("invalid key")
)

// KvPairs are key value pairs that may be used to scope and select individual items.
type KvPairs []*KvPair

// KvPair is a key value . It implements KvPairs
type KvPair struct {
	Key   string `json:"key" toml:"key"`
	Value string `json:"value,omitempty" toml:"value,omitempty"`
}

// Set create or update the key-value pair entry for "key" to "value".
func (ls *KvPairs) Set(key string, value interface{}) error {

	kvp_mu.Lock()
	defer kvp_mu.Unlock()

	if len(key) < 1 {
		return kvp_ErrKeyEmpty
	}

	if len(key) > 100 {
		return kvp_ErrKeyLength
	}

	svalue := fmt.Sprintf("%v", value)

	for i, prev := range *ls {

		if prev.Key == key {

			if prev.Value != svalue {
				(*ls)[i].Value = svalue
			}

			return nil
		}
	}

	*ls = append(*ls, &KvPair{
		Key:   key,
		Value: svalue,
	})

	return nil
}

// Get fetch the key-value pair "value" (if any) for "key".
func (ls KvPairs) Get(key string) Bytex {

	kvp_mu.RLock()
	defer kvp_mu.RUnlock()

	for _, prev := range ls {

		if prev.Key == key {
			return Bytex(prev.Value)
		}
	}

	return nil
}

// Del remove the key-value pair (if any) for "key".
func (ls *KvPairs) Del(key string) {

	kvp_mu.Lock()
	defer kvp_mu.Unlock()

	for i, prev := range *ls {

		if prev.Key == key {
			*ls = append((*ls)[:i], (*ls)[i+1:]...)
			break
		}
	}
}

func (ls *KvPairs) Equal(items KvPairs) bool {

	kvp_mu.RLock()
	defer kvp_mu.RUnlock()

	if len(*ls) != len(items) {
		return false
	}

	for _, v := range *ls {

		hit := false

		for _, v2 := range items {

			if v.Key != v2.Key {
				continue
			}

			if v.Value != v2.Value {
				return false
			}

			hit = true
			break
		}

		if !hit {
			return false
		}
	}

	return true
}
