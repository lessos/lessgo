## lessgo/logger
lessgo/logger is a Simplified implementation of Leveled logs ([google glog](https://github.com/google/glog)).

## Example

create a new test file main.go, and input the following codes

```go
package main

import (
    "flag"
    "github.com/lessos/lessgo/logger"
)

func main() {

    // init output dir
    flag.Parse()

    // API:Print
    logger.Print("error", "the error code/message: ", 400, "/", "bad request")

    // API::Printf
    logger.Printf("error", "the error code/message: %d/%s", 400, "bad request")

    select {}
}
```

build the main.go file, run it and output the log into stderr console

```shell
go build main.go
./main -logtostderr=true
2016-04-18 22:24:26.159220 test2.go:14] ERROR the error code/message: 404/bad request
2016-04-18 22:24:26.159236 test2.go:17] ERROR the error code/message: 404/bad request
```

or run it and output the log into file
```shell
./main -log_dir="/var/log/"
```

the output file name will formated like this

```
/var/log/{process name}.{hostname}.{current os username}.{date}.{pid}.log
```
