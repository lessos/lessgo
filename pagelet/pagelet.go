package pagelet

import (
	"../deps/go.net/websocket"
	"../logger"
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

func Run() {

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
		logger.Printf("info", "lessgo/pagelet: Listening on port %d ...", Config.HttpPort)
	}()
}

func handle(w http.ResponseWriter, r *http.Request) {

	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" || upgrade == "Websocket" {
		websocket.Handler(func(ws *websocket.Conn) {
			r.Method = "WS"
			handleInternal(w, r, ws)
		}).ServeHTTP(w, r)
	} else {
		handleInternal(w, r, nil)
	}
}

func handleInternal(w http.ResponseWriter, r *http.Request, ws *websocket.Conn) {

	defer func() {

		if err := recover(); err != nil {
			logger.Printf("error", "handleInternal Panic on %s", err)
		}

		r.Body.Close()
	}()

	var (
		req  = NewRequest(r)
		resp = NewResponse(w)
		c    = NewController(req, resp)
	)

	if ws != nil {
		req.WebSocket = WebSocket{
			conn: ws,
		}
	}

	for _, filter := range Filters {
		filter(c)
	}
}
