package main

import (
	"./ssdb"
	"fmt"
)

func main() {

	conn, err := ssdb.NewConnector(ssdb.Config{
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
	if conn.Cmd("set", "aa", "val-aaaaaaaaaaaaaaaaaa").Bool() {
		fmt.Println("set OK")
	}
	// API::String() string
	if rs := conn.Cmd("get", "aa"); rs.State == "ok" {
		fmt.Println("get OK\n\t", rs.String())
	}
	// API::Hash() []Entry
	conn.Cmd("set", "bb", "val-bbbbbbbbbbbbbbbbbb")
	conn.Cmd("set", "cc", "val-cccccccccccccccccc")
	if rs := conn.Cmd("multi_get", "aa", "bb"); rs.State == "ok" {
		fmt.Println("multi_get OK")
		for _, v := range rs.Hash() {
			fmt.Println("\t", v.Key, v.Value)
		}
	}
	if rs := conn.Cmd("scan", "", "", 10); rs.State == "ok" {
		fmt.Println("scan OK")
		for _, v := range rs.Hash() {
			fmt.Println("\t", v.Key, v.Value)
		}
	}

	conn.Cmd("zset", "z", "a", 3)
	conn.Cmd("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	if rs := conn.Cmd("zrscan", "z", "", "", "", 10); rs.State == "ok" {
		fmt.Println("zrscan OK")
		for _, v := range rs.Hash() {
			fmt.Println("\t", v.Key, v.Value)
		}
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
		for _, v := range rs.Hash() {
			fmt.Println("\t", v.Key, v.Value)
		}
	}

	// API::List() []string
	conn.Cmd("qpush", "queue", "q-1111111111111")
	conn.Cmd("qpush", "queue", "q-2222222222222")
	if rs := conn.Cmd("qpop", "queue", 10); rs.State == "ok" {
		fmt.Println("qpop OK")
		for k, v := range rs.List() {
			fmt.Println("\t", k, v)
		}
	}
}
