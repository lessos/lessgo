package pagelet

import (
    "fmt"
    "net/http"
    "strconv"
    //"strings"
    "time"
)

// APIs
//
//  pagelet.Config.UrlBasePath string
//  pagelet.Config.HttpPort int
//
//  pagelet.Config.ViewPath(module, path string)
//  pagelet.Config.RouteAppend(patten, action string)
//  pagelet.Config.RouteStaticAppend(patten, path string)
//
//  pagelet.RegisterController(module string, (*Object)(nil), []string)
//  pagelet.Run()

var (
    Config = ConfigBase{
        HttpPort: 1024,
    }
    MainRouter         = &Router{Routes: []Route{}, Modules: map[string][]Route{}}
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
    /*
       MainRouter.RouteStaticAppend("/static", "static")

       //
       route := Route{
           Type:    "std",
           Path:    "/:controller/:action",
           Tree:    []string{":controller", ":action"},
           TreeLen: 2,
       }
       MainRouter.Routes = append(MainRouter.Routes, route)
    */

    //Println("config", Config)
    //for k, v := range Config.ViewPaths {

    //}

    //paths := strings.Split(Config.ViewPaths["def"], ",")
    //paths = append(paths, "../src/views")
    //paths = append(paths, "src/views")

    //Println(Config)
    //
    MainTemplateLoader = NewTemplateLoader()

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
