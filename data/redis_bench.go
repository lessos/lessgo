package main

import (
	"./redis"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"
)

var (
	host     *string = flag.String("h", "127.0.0.1", "Server hostname (default 127.0.0.1)")
	port     *string = flag.String("p", "6379", "Server port (default 6379)")
	maxconn  *int    = flag.Int("c", 50, "Number of parallel connections (default 50)")
	requests *int    = flag.Int("n", 10000, "Total number of requests (default 10000)")
	dsize    *int    = flag.Int("d", 2, "Data size of SET/GET value in bytes (default 2)")
)

func benchmark(c *redis.Connector, title, command string, params ...interface{}) {

	fmt.Printf("====== %s ======\n", title)

	reqs := make(chan int, *requests)
	for i := 1; i <= *requests; i++ {
		reqs <- i
	}
	maxcs := *maxconn * 2
	done := make(chan int, 2)

	start := time.Now()

	for i := 0; i < maxcs; i++ {

		go func() {
			for {
				req := <-reqs
				c.Cmd(command, params...)
				if req == *requests {
					done <- 1
					break
				}
			}
		}()
	}

	select {
	case <-done:
		duration := time.Now().Sub(start)
		fmt.Printf("\t%d Requests completed in %v\n", *requests, time.Since(start))
		rps := float64(*requests) / duration.Seconds()
		fmt.Printf("\tRequests per second %.0f\n", rps)
	}
}

func main() {

	var data string
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *requests < 1000 {
		*requests = 1000
	}

	for i := 0; i < *dsize; i++ {
		data += "x"
	}

	cfg := redis.Config{
		Host:    *host,
		Port:    *port,
		MaxConn: *maxconn,
	}

	c, e := redis.NewConnector(cfg)
	if e != nil {
		log.Println(e)
	}

	fmt.Printf("### Parallel connections %d ###\n", *maxconn)

	c.Cmd("flushall")

	benchmark(c, "PING", "PING")
	benchmark(c, "SET", "SET", "keys:rand:000000000000", data)
	benchmark(c, "GET", "GET", "keys:rand:000000000000")
	benchmark(c, "INCR", "keys:counter:rand:000000000000")
	benchmark(c, "INCR", "LPUSH", "keys:rand:list", data)
	benchmark(c, "LPOP", "LPOP", "keys:rand:list")
	benchmark(c, "SADD", "SADD", "keys:rand:sets", "counter:rand:000000000000")
	benchmark(c, "SPOP", "SPOP", "keys:rand:sets")

	benchmark(c, "LPUSH", "LPUSH", "keys:rand:list", data)
	benchmark(c, "LRANGE 100", "LRANGE", "keys:rand:list", 0, 99)
	benchmark(c, "LRANGE 300", "LRANGE", "keys:rand:list", 0, 299)
	benchmark(c, "LRANGE 600", "LRANGE", "keys:rand:list", 0, 599)

	args := make([]interface{}, 20)
	for i := 0; i < 20; i += 2 {
		args[i] = "keys:rand:000000000000"
		args[i+1] = data
	}
	benchmark(c, "MSET (10 keys)", "MSET", args...)

	for i := 0; i < 100000; i++ {
		c.Cmd("ZADD", "keys:rand:zset", i, i)
	}
	benchmark(c, "ZRANGE 10w100", "ZRANGE", "keys:rand:zset", 50000, 50100)

	c.Cmd("flushall")
}
