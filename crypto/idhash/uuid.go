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
	"fmt"
)

const (
	// Set version (4) and variant (2).
	uuid_version = 4 << 4
	uuid_variant = 2 << 4
)

// RandUUID generates a new UUID based on version 4.
func RandUUID() string {

	bs := Rand(16)

	bs[6] = uuid_version | (bs[6] & 15)
	bs[8] = uuid_variant | (bs[8] & 15)

	return fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
}
