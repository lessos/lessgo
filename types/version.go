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

type Version string

// Compare compares this version to another version. This
// returns -1, 0, or 1 if this version is smaller, equal,
// or larger than the compared version, respectively.
func (v *Version) Compare(other *Version) int {

	if *v == *other {
		return 0
	}

	vs, vs2 := v.parse(), other.parse()

	vslen := len(vs)
	if len(vs) > len(vs2) {
		vslen = len(vs2)
	}

	for i := 0; i < vslen; i++ {
		if lg := vs[i] - vs2[i]; lg > 0 {
			return 1
		} else if lg < 0 {
			return -1
		}
	}

	if lg := len(vs) - len(vs2); lg > 0 {
		return 1
	} else if lg < 0 {
		return -1
	}

	return 0
}

func (v *Version) String() string {
	return string(*v)
}

func (v *Version) parse() []int32 {

	var (
		segments = []int32{}
		num      = int32(-1)
	)

	for _, char := range *v {

		if char >= '0' && char <= '9' {

			if num > -1 {
				num *= 10
			} else {
				num = 0
			}

			num += char - '0'

		} else {

			if num > -1 {
				segments = append(segments, num)
			}

			if char >= 'a' && char <= 'z' {
				segments = append(segments, char-'a'+10)
			}

			num = -1
		}
	}

	if num > -1 {
		segments = append(segments, num)
	}

	return segments
}
