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

package sync

import (
	"runtime"
)

type PermitPool struct {
	size uint16
	free chan uint16
}

func NewPermitPool(size int) *PermitPool {

	if size < 1 {
		size = runtime.NumCPU()
	}

	if size > 65535 {
		size = 65535
	}

	p := &PermitPool{
		size: uint16(size),
		free: make(chan uint16, size),
	}

	for i := uint16(0); i < p.size; i++ {
		p.free <- i
	}

	return p
}

func (p *PermitPool) Pull() uint16 {
	return <-p.free
}

func (p *PermitPool) Push(num uint16) {
	p.free <- num
}
