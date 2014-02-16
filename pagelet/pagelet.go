package pagelet

import (
    "fmt"
    "net/http"
    "time"
)

var (
    MainRouter         = &Router{Routes: []Route{}}
    MainTemplateLoader *TemplateLoader
    Server             *http.Server
)

func Println(args ...interface{}) {
    fmt.Println(args...)
}

func Run(port string) {

    //
    MainRouter.RouteStaticAppend("/static", "static")

    //
    route := Route{
        Type:    "std",
        Path:    "/:controller/:action",
        Tree:    []string{":controller", ":action"},
        TreeLen: 2,
    }
    MainRouter.Routes = append(MainRouter.Routes, route)

    fmt.Println(MainRouter.Routes)

    //
    MainTemplateLoader = NewTemplateLoader([]string{"../src/views", "src/views"})

    go func() {

        Server = &http.Server{
            Addr:           ":" + port,
            Handler:        http.HandlerFunc(handle),
            ReadTimeout:    30 * time.Second,
            WriteTimeout:   30 * time.Second,
            MaxHeaderBytes: 1 << 20,
        }

        Server.ListenAndServe()
    }()

    go func() {
        time.Sleep(100 * time.Millisecond)
        Println("lessgo/pagelet: Listening on port %s ...", port)
    }()
}

func handle(w http.ResponseWriter, r *http.Request) {

    defer func() {

        if err := recover(); err != nil {
            Println("handle", err)
        }
    }()

    var (
        req  = NewRequest(r)
        resp = NewResponse(w)
        c    = NewController(req, resp)
    )

    Filters[0](c, Filters[1:])
}
