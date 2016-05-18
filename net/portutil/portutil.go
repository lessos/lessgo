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

package portutil

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func FreeRange(port_start, limit uint16) (uint16, uint16, error) {

	if port_start > 60000 {
		port_start = 60000
	} else if port_start < 1024 {
		port_start = 1024
	}

	port_end := port_start

	if limit > 500 {
		limit = 500
	} else if limit < 1 {
		limit = 1
	}

	for ; port_end < 65535; port_end++ {

		if !IsFree(port_end) {
			port_start = port_end + 1
			continue
		}

		if port_start+limit <= port_end {
			return port_start, port_end, nil
		}
	}

	return 0, 0, errors.New("Not Enough Ports")
}

func IsFree(port uint16) bool {

	if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		ln.Close()
		return true
	}

	return false
}

func Free(port_start, limit uint16) (uint16, error) {

	if port_start > 60000 {
		port_start = 60000
	} else if port_start < 1 {
		port_start = 1
	}

	if limit < 100 {
		limit = 100
	} else if 65535-port_start < limit {
		limit = 65535 - port_start
	}

	for try := uint16(0); try < 100; try++ {

		if port := uint16(int(port_start) + r.Intn(int(limit))); IsFree(port) {
			return port, nil
		}
	}

	return 0, errors.New("No Free Port")
}
