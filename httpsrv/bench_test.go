// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package httpsrv

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/lessos/lessgo/net/httpclient"
)

const (
	test_bench_string = "hello world"
)

type Index struct {
	*Controller
}

func (c Index) IndexAction() {
	c.RenderString(test_bench_string)
}

func TestBench(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//
	module := NewModule("default")
	module.ControllerRegister(new(Index))
	GlobalService.ModuleRegister("/", module)

	//
	port := net_free_port()
	if port == 0 {
		t.Fatal(errors.New("Listen failed"))
	}
	GlobalService.Config.HttpPort = uint16(port)
	go GlobalService.Start()

	time.Sleep(1e9)

	//
	max_cli := runtime.NumCPU() * 20
	if max_cli > 1000 {
		max_cli = 1000
	}

	max_req := runtime.NumCPU() * 5000
	if max_req > 100000 {
		max_req = 100000
	}

	clients := make(chan int, max_cli)
	for i := 0; i < max_cli; i++ {
		clients <- 1
	}

	var (
		mu     sync.Mutex
		start  = time.Now()
		ok_no  = 0
		err_no = 0
		url    = fmt.Sprintf("http://127.0.0.1:%d/", port)
	)

	for i := 0; i <= max_req; i++ {

		if i == max_req {

			for {
				if len(clients) >= max_cli {
					break
				}

				time.Sleep(5e7)
			}

			break
		}

		_ = <-clients

		go func() {

			s, _ := httpclient.Get(url).ReplyString()

			mu.Lock()
			if s == test_bench_string {
				ok_no++
			} else {
				err_no++
			}
			mu.Unlock()

			clients <- 1
		}()
	}

	//
	duration := time.Now().Sub(start)
	rps := float64(max_req) / duration.Seconds()

	//
	t.Logf("\tTime taken for tests:   %v\n", time.Since(start))
	t.Logf("\tConcurrency Level:      %d\n", max_cli)
	t.Logf("\tComplete requests:      %d\n", ok_no)
	t.Logf("\tFailed requests:        %d\n", err_no)
	t.Logf("\tRequests per second:    %.0f\n", rps)
}

func net_free_port() int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {

		iport := 10000 + r.Intn(50000)

		port := strconv.Itoa(iport)
		ln, err := net.Listen("tcp", ":"+port)
		if err == nil {
			ln.Close()
			return iport
		}
	}

	return 0
}
