## lessgo/locker
lock toolkit for creating finer-grained locking to handle concurrent read/write, parallel processing or multi-tasking for business logic.


## Pool Locker
Pool lock allows multiple task to acquire fixed number of locks in the same time, and perform their own logic.

It can be used to specify the number of multi file downloading, limit number of connections, ...


```go
package main

import (
	"fmt"

	"github.com/lessos/lessgo/locker"
)
func main() {

	p := locker.NewPool(4)

	file_urls, num := []string{}, 100

	for i := 1; i <= num; i++ {
		file_urls = append(file_urls, fmt.Sprintf("http://www.example.com/file.%d", i))
	}

	done, done_num := make(chan bool, num), 0

	for _, url := range file_urls {

		go func(url string) {

			p.Lock()
			defer p.Unlock()

			// logic to download this file
			fmt.Printf("download %s\n", url)

			done <- true

		}(url)
	}

	for <-done {

		done_num++

		if done_num == num {
			fmt.Println("well done")
			break
		}
	}
}
```


## HashPool Locker
HashPool is a collection of multiple locks, locking and unlocking with a specific key. the key is hashed to a fixed lock in the pool.

It can be used to specify the number of concurrent transactions, multi-tasking, multi-shard data read and write, multi-queue processing, ...


```go
package main

import (
	"fmt"
	"runtime"

	"github.com/lessos/lessgo/locker"
)

func main() {

	hp := locker.NewHashPool(runtime.NumCPU())

	file_urls, num := []string{}, 100

	for i := 1; i <= num; i++ {
		file_urls = append(file_urls, fmt.Sprintf("http://www.example.com/file.%d", i))
	}

	done, done_num := make(chan bool, num), 0

	for _, url := range file_urls {

		go func(url string) {

			hp.Lock([]byte(url))
			defer hp.Unlock([]byte(url))

			// logic to download this file
			fmt.Printf("download %s\n", url)

			done <- true

		}(url)
	}

	for <-done {

		done_num++

		if done_num == num {
			fmt.Println("well done")
			break
		}
	}
}
```

