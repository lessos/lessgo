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

package locker

import (
	"hash/crc32"
	"runtime"
	_sync "sync"
)

type HashPool struct {
	mu    _sync.Mutex
	size  uint32
	pools map[uint16]*_sync.Mutex
}

func NewHashPool(size int) *HashPool {

	if size < 1 {
		size = runtime.NumCPU()
	}

	if size > 128 {
		size = 128
	}

	p := &HashPool{
		size:  uint32(size),
		pools: map[uint16]*_sync.Mutex{},
	}

	return p
}

func (p *HashPool) Lock(bs []byte) (key uint16) {

	key = uint16(crc32.ChecksumIEEE(bs) % p.size)

	p.mu.Lock()
	v, _ := p.pools[key]
	if v == nil {
		v = &_sync.Mutex{}
		p.pools[key] = v
	}
	p.mu.Unlock()

	v.Lock()

	return key
}

func (p *HashPool) Unlock(bs []byte) {
	p.UnlockWithKey(uint16(crc32.ChecksumIEEE(bs) % p.size))
}

func (p *HashPool) UnlockWithKey(key uint16) {

	p.mu.Lock()
	v, _ := p.pools[key]
	p.mu.Unlock()

	if v != nil {
		v.Unlock()
	}
}
