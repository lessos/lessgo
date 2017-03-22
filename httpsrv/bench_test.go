// Copyright 2016 lessOS.com, All rights reserved.
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
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	//. "github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/utils"
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
	err, port, _ := utils.NetFreePort(10000, 20000)
	if err != nil {
		t.Fatal(err)
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
	fmt.Printf("\tTime taken for tests:   %v\n", time.Since(start))
	fmt.Printf("\tConcurrency Level:      %d\n", max_cli)
	fmt.Printf("\tComplete requests:      %d\n", ok_no)
	fmt.Printf("\tFailed requests:        %d\n", err_no)
	fmt.Printf("\tRequests per second:    %.0f\n", rps)
}
