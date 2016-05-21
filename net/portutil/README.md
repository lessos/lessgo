## net/portutil

network port utils

## Example
```go
package main

import (
	"fmt"

	"github.com/lessos/lessgo/net/portutil"
)

func main() {

	if port, err := portutil.Free(5000, 1000); err == nil {
		fmt.Println("free port ", port)
	}

	if start, end, err := portutil.FreeRange(5500, 50000); err == nil {
		fmt.Printf("free ports %d~%d\n", start, end)
	}

	fmt.Println("check   80", portutil.IsFree(80))
}
```

