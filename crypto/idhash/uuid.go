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
	"encoding/hex"
)

const (
	uuid_dash byte = '-'
)

// RandUUID generates a new UUID based on version 4.
func RandUUID() string {
	return uuid(Rand(16))
}

func HashUUID(bs []byte) string {
	return uuid(HashSum(AlgSha1, bs, 16))
}

func uuid(bs []byte) string {

	bs[6] = (bs[6] & 0x0F) | 0x40 // version 4
	bs[8] = (bs[8] & 0x3F) | 0x80 // variant rfc4122

	uuid := make([]byte, 36)

	hex.Encode(uuid[0:8], bs[0:4])
	uuid[8] = uuid_dash

	hex.Encode(uuid[9:13], bs[4:6])
	uuid[13] = uuid_dash

	hex.Encode(uuid[14:18], bs[6:8])
	uuid[18] = uuid_dash

	hex.Encode(uuid[19:23], bs[8:10])
	uuid[23] = uuid_dash

	hex.Encode(uuid[24:], bs[10:])

	return string(uuid)
}
