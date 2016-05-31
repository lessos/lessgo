# iossdb

iossdb is a minimalistic, connection pooling Go Client for SSDB (http://ssdb.io).

## APIs
* Connect. iossdb use iossdb.NewConnector(iossdb.Config{...}) to create connection with SSDB server. You can use iossdb.Config to set host, port, pool size, timeout, etc.

* Request. all SSDB operations go with ```iossdb.Connector.Cmd()```, it accepts variable arguments. The first argument of Cmd() is the SSDB command, for example "get", "set", etc. The rest arguments(maybe none) are the arguments of that command.

* Response. the iossdb.Connector.Cmd() method will return an Object of iossdb.Reply

	* State:  The element of iossdb.Reply.State is the response code, ```"ok"``` means the current command are valid results. The response code may be ```"not_found"``` if you are calling "get" on an non-exist key.

	* Data: The element of iossdb.Reply.Data is the response data. You can also use the following method to get a dynamic data struct what you want to need.
		* iossdb.Reply.String() string
		* iossdb.Reply.Bool() bool
		* iossdb.Reply.Int() int
		* iossdb.Reply.Int8() int8
		* iossdb.Reply.Int16() int16
		* iossdb.Reply.Int32() int32
		* iossdb.Reply.Int64() int64
		* iossdb.Reply.Uint() uint
		* iossdb.Reply.Uint8() uint8
		* iossdb.Reply.Uint16() uint16
		* iossdb.Reply.Uint32() uint32
		* iossdb.Reply.Uint64() uint64
		* iossdb.Reply.Float32() float32
		* iossdb.Reply.Float64() float64
		* iossdb.Reply.List() []types.Bytex
		* iossdb.Reply.KvEach(fn func(key, value types.Bytex)) int
		* iossdb.Reply.KvLen() int
		* iossdb.Reply.JsonDecode(obj interface{}) error

* Refer to [Official API documentation](http://ssdb.io/docs/) to checkout a complete list of all avilable commands.

## Example
```go
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

package main

import (
	"fmt"

	"github.com/lessos/lessgo/data/iossdb"
	"github.com/lessos/lessgo/types"
)

func main() {

	conn, err := iossdb.NewConnector(iossdb.Config{
		Host:    "127.0.0.1",
		Port:    6380,
		Timeout: 3,  // timeout in second, default to 10
		MaxConn: 10, // max connection number, default to 1
	})
	if err != nil {
		fmt.Println("Connect Error:", err)
		return
	}
	defer conn.Close()

	// API::Bool() bool
	conn.Cmd("set", "true", "True")
	if conn.Cmd("get", "true").Bool() {
		fmt.Println("set bool OK")
	}

	conn.Cmd("set", "aa", "val-aaaaaaaaaaaaaaaaaa")
	conn.Cmd("set", "bb", "val-bbbbbbbbbbbbbbbbbb")
	conn.Cmd("set", "cc", "val-cccccccccccccccccc")
	// API::String() string
	if rs := conn.Cmd("get", "aa"); rs.State == "ok" {
		fmt.Println("get OK\n\t", rs.String())
	}
	// API::KvEach()
	if rs := conn.Cmd("multi_get", "aa", "bb"); rs.State == "ok" {
		fmt.Println("multi_get OK")
		rs.KvEach(func(key, value types.Bytex) {
			fmt.Println("\t", key.String(), value.String())
		})
	}
	// API::KvEach() bytes
	bkeys := [][]byte{[]byte("aa"), []byte("bb")}
	if rs := conn.Cmd("multi_get", bkeys); rs.State == "ok" {
		fmt.Println("multi_get bytes OK")
		rs.KvEach(func(key, value types.Bytex) {
			fmt.Println("\t", key, value)
		})
	}

	if rs := conn.Cmd("scan", "aa", "cc", 10); rs.State == "ok" {
		fmt.Println("scan OK")
		n := rs.KvEach(func(key, value types.Bytex) {
			fmt.Println("\t", key.String(), value.String())
		})
		fmt.Println("\t got", n)
	}

	conn.Cmd("zset", "z", "a", 3)
	conn.Cmd("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	if rs := conn.Cmd("zrscan", "z", "", "", "", 10); rs.State == "ok" {
		fmt.Println("zrscan OK")
		rs.KvEach(func(key, value types.Bytex) {
			fmt.Println("\t", key.String(), value.Int())
		})
	}

	conn.Cmd("set", "key", 10)
	if rs := conn.Cmd("incr", "key", 1).Int(); rs > 0 {
		fmt.Println("incr OK\n\t", rs)
	}

	// API::Int() int
	// API::Int64() int64
	conn.Cmd("setx", "key", 123456, 300)
	if rs := conn.Cmd("ttl", "key").Int(); rs > 0 {
		fmt.Println("ttl OK\n\t", rs)
	}

	if rs := conn.Cmd("multi_hset", "zone", "c1", "v-01", "c2", "v-02"); rs.State == "ok" {
		fmt.Println("multi_hset OK")
	}
	if rs := conn.Cmd("multi_hget", "zone", "c1", "c2"); rs.State == "ok" {
		fmt.Println("multi_hget OK")
		rs.KvEach(func(key, value types.Bytex) {
			fmt.Println("\t", key.String(), value.String())
		})
	}

	// API::Float64() float64
	conn.Cmd("set", "float", 123.456)
	if rs := conn.Cmd("get", "float").Float64(); rs > 0 {
		fmt.Println("float OK\n\t", rs)
	}

	// API::List()
	conn.Cmd("qpush", "queue", "q-1111111111111")
	conn.Cmd("qpush", "queue", "q-2222222222222")
	if rs := conn.Cmd("qpop", "queue", 10); rs.State == "ok" {
		fmt.Println("qpop list OK")
		for i, value := range rs.List() {
			fmt.Println("\t", i, value.String())
		}
	}

	// iossdb.Reply.JsonDecode(obj interface{}) error
	conn.Cmd("set", "json_key", "{\"name\": \"test obj.name\", \"value\": \"test obj.value\"}")
	if rs := conn.Cmd("get", "json_key"); rs.State == "ok" {
		var rs_obj struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		if err := rs.JsonDecode(&rs_obj); err == nil {
			fmt.Println("JsonDecode OK")
			fmt.Println("\tname :", rs_obj.Name)
			fmt.Println("\tvalue:", rs_obj.Value)
		} else {
			fmt.Println("json_key ERR", err)
		}
	}
}
```

