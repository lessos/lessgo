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

package idhash

import (
	"encoding/base64"
)

func RandBase64String(length int) string {

	if length < 4 {
		length = 1
	} else if length%4 > 0 {
		length = length/4 + 1
	} else {
		length = length / 4
	}

	return base64.RawStdEncoding.EncodeToString(Rand(3 * length))
}

func HashToBase64String(alg int, bs []byte, length int) string {

	if length < 4 {
		length = 1
	} else if length%4 > 0 {
		length = length/4 + 1
	} else {
		length = length / 4
	}

	return base64.RawStdEncoding.EncodeToString(HashSum(alg, bs, 3*length))
}
