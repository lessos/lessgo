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
	"crypto/md5"
	"crypto/rand"
	mrand "math/rand"
	"time"
)

const (
	rand_bytes_max = 1024 * 1024
)

func init() {
	mrand.Seed(time.Now().UTC().UnixNano())
}

func Rand(size int) []byte {

	if size < 1 {
		size = 1
	} else if size > rand_bytes_max {
		size = rand_bytes_max
	}

	bs := make([]byte, size)

	// rand.Read() is a helper function that calls rand.Reader.Read using io.ReadFull.
	// On return, n == len(b) if and only if err == nil.
	//
	// rand.Reader is a global, shared instance of a cryptographically
	// strong pseudo-random generator.
	//
	// On Unix-like systems, rand.Reader reads from /dev/urandom.
	// On Linux, rand.Reader uses getrandom(2) if available, /dev/urandom otherwise.
	// On Windows systems, rand.Reader uses the CryptGenRandom API.
	if _, err := rand.Read(bs); err != nil {
		for i := range bs {
			bs[i] = uint8(mrand.Intn(256))
		}
	}

	return bs
}

func Hash(bs []byte, bytelen int) []byte {

	if bytelen < 1 {
		bytelen = 1
	} else if bytelen > 16 {
		bytelen = 16
	}

	bs_hash := md5.Sum(bs)

	return bs_hash[:bytelen]
}
