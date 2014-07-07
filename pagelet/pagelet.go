package pagelet

import (
	"fmt"
	"net/http"
	"strconv"
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
		HttpPort:        1024,
		LocaleCookieKey: "lang",
	}
	MainRouter         = &Router{Routes: []Route{}, Modules: map[string][]Route{}}
	MainTemplateLoader *TemplateLoader
	Server             *http.Server
)

func println(args ...interface{}) {
	fmt.Println(args...)
}
func printf(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}

func Run() {

	//println(Config)
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
		printf("lessgo/pagelet: Listening on port %d ...", Config.HttpPort)
	}()
}

func handle(w http.ResponseWriter, r *http.Request) {

	defer func() {

		if err := recover(); err != nil {
			println("handle", err)
		}

		r.Body.Close()
	}()

	var (
		req  = NewRequest(r)
		resp = NewResponse(w)
		c    = NewController(req, resp)
	)

	for _, filter := range Filters {
		filter(c)
	}
}
