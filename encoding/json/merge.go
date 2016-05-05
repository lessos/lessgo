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

package json

import (
	"reflect"
)

const (
	jsonMergeDepth = 32
)

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
//
// Refer
// 	https://github.com/peterbourgon/mergemap
// func Merge(dst, src map[string]interface{}) map[string]interface{} {
// 	return jsMerge(dst, src, 0)
// }

func jsMerge(dst, src map[string]interface{}, depth int) map[string]interface{} {

	if depth > jsonMergeDepth {
		return dst
		// panic("too deep!")
	}

	for key, srcVal := range src {

		if dstVal, ok := dst[key]; ok {

			srcMap, srcMapOk := jsMapify(srcVal)
			dstMap, dstMapOk := jsMapify(dstVal)

			if srcMapOk && dstMapOk {
				srcVal = jsMerge(dstMap, srcMap, depth+1)
			}
		}

		dst[key] = srcVal
	}

	return dst
}

func jsMapify(i interface{}) (map[string]interface{}, bool) {

	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Map {

		m := map[string]interface{}{}

		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}

		return m, true
	}

	return map[string]interface{}{}, false
}
