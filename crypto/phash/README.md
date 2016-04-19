## crypto/phash

High Security Password Hashing framework 

* Scrypt based
* random salt to against dictionary attacks
* high CPU cost to against brute force attacks
* high Memory cost to against GPU/FPGA attacks

## Example
```go
package main

import (
    "fmt"
    "github.com/lessos/lessgo/crypto/phash"
)

func main() {

    password := "123456"

    //
    hashtxt, _ := phash.Generate(password)
    fmt.Println("password      :", password)
    fmt.Println("hashed text   :", hashtxt)

    //
    if phash.Verify(password, hashtxt) {
        fmt.Println("verify", password, ": ok")
    } else {
        fmt.Println("verify", password, ": failed")
    }

    //
    password = "abcdef"
    if phash.Verify(password, hashtxt) {
        fmt.Println("verify", password, ": ok")
    } else {
        fmt.Println("verify", password, ": failed")
    }
}
```

build and run this example

```go
go run main.go

password      : 123456
hashed text   : L001f812OY+rtPErWnlos26MKGgmzs/U7a7VFXLqYa/PhZ/PkCINeUycCuJOpb588SRwcOQ03AL
verify 123456 : ok
verify abcdef : failed
```
