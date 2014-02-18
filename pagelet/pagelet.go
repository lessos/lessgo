package pagelet

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var (
    Config             = ConfigStruct{HttpPort: 1024}
    MainRouter         = &Router{Routes: []Route{}}
    MainTemplateLoader *TemplateLoader
    Server             *http.Server
)

func Println(args ...interface{}) {
    fmt.Println(args...)
}
func Printf(str string, args ...interface{}) {
    fmt.Printf(str+"\n", args...)
}

func Run() {

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

    paths := strings.Split(Config.ViewPaths, ",")
    paths = append(paths, "../src/views")
    paths = append(paths, "src/views")

    //
    MainTemplateLoader = NewTemplateLoader(paths)

    go func() {

        Server = &http.Server{
            Addr:           ":" + strconv.Itoa(Config.HttpPort),
            Handler:        http.HandlerFunc(handle),
            ReadTimeout:    30 * time.Second,
            WriteTimeout:   30 * time.Second,
            MaxHeaderBytes: 1 << 20,
        }

        Server.ListenAndServe()
    }()

    go func() {
        time.Sleep(100 * time.Millisecond)
        Printf("lessgo/pagelet: Listening on port %d ...", Config.HttpPort)
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
