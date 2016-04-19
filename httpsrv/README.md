## lessgo/httpsrv
lessgo/httpsrv is a Lightweight, Modular, High Performance MVC web framework.

## Quick Start

Install httpsrv framework

```shell
go get -u github.com/lessos/lessgo/httpsrv

```

first hello world demo

```go
package main

import (
    "github.com/lessos/lessgo/httpsrv"
)

type Index struct {
    *httpsrv.Controller
}

func (c Index) IndexAction() {
    c.RenderString("hello world")
}

func main() {

    // init one module
    module := httpsrv.NewModule("default")
    
    // register controller to module
    module.ControllerRegister(new(Index))

    // register module to httpsrv
    httpsrv.GlobalService.ModuleRegister("/", module)

    // listening on port 18080
    httpsrv.GlobalService.Config.HttpPort = 18080

    // start
    httpsrv.GlobalService.Start()
}
```

