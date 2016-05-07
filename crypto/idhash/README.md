## crypto/idhash

ID hashing utils

## Example
```go
package main

import (
    "fmt"
    "github.com/lessos/lessgo/crypto/idhash"
)

func main() {

    fmt.Println(idhash.Rand(16))

    fmt.Println(idhash.RandHexString(16))

    fmt.Println(idhash.HashToHexString([]byte("123456"), 16))

    fmt.Println(idhash.RandUUID())

    fmt.Println(idhash.RandBase64String(16))
}

```

build and run this example

```go
go run main.go

[123 71 115 5 82 134 54 161 92 1 113 33 138 138 152 101]
662a59b00ab4198f
e10adc3949ba59ab
c37ee956-7b82-4da6-291a-591ffe3b3c38
38uSeXAayTmoZjVh
```
