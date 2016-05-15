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
	"runtime"
)

type Pool struct {
	size uint16
	free chan bool
}

func NewPool(size int) *Pool {

	if size < 1 {
		size = runtime.NumCPU()
	}

	if size > 65535 {
		size = 65535
	}

	p := &Pool{
		size: uint16(size),
		free: make(chan bool, size),
	}

	for i := uint16(0); i < p.size; i++ {
		p.free <- true
	}

	return p
}

func (p *Pool) Lock() {
	<-p.free
}

func (p *Pool) Unlock() {
	p.free <- true
}
