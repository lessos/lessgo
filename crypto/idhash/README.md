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

	fmt.Println(idhash.Hash([]byte("123456"), 8))

	fmt.Println(idhash.RandHexString(16))

	fmt.Println(idhash.HashToHexString([]byte("123456"), 16))

	fmt.Println(idhash.RandUUID())

	fmt.Println(idhash.RandBase64String(16))
}
```

build and run this example

```go
go run main.go

[73 235 46 128 5 198 240 75 220 102 146 115 226 166 6 5]
[225 10 220 57 73 186 89 171]
4e04ac45fd243f68
e10adc3949ba59ab
df8d9076-9f7e-445e-a703-32df8204c8ff
skNhSXFeTKKh0SBH
```

