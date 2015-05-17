// Copyright 2015 lessOS.com, All rights reserved.
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

func (ls LabelListMeta) Insert(key, value string) LabelListMeta {

	for i, prev := range ls {

		if prev.Key == key {

			if prev.Val != value {
				ls[i].Val = value
			}

			return ls
		}
	}

	ls = append(ls, LabelMeta{Key: key, Val: value})

	return ls
}

func (ls LabelListMeta) Fetch(key string) (string, bool) {

	for _, prev := range ls {

		if prev.Key == key {
			return prev.Val, true
		}
	}

	return "", false
}

func (ls LabelListMeta) Remove(key string) {

	for i, prev := range ls {

		if prev.Key == key {
			ls = append(ls[:i], ls[i:]...)
			break
		}
	}
}
