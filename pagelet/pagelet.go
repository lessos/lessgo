package pagelet

import (
	"github.com/lessos/lessgo/deps/go.net/websocket"
	"github.com/lessos/lessgo/logger"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		HttpPort:         0,
		LocaleCookieKey:  "lang",
		SessionCookieKey: "access_token",
	}
	MainRouter = &Router{
		Routes:  []Route{},
		Modules: map[string][]Route{},
	}
	MainTemplateLoader = &TemplateLoader{
		templatePaths: map[string]string{},
		templateSets:  map[string]*template.Template{},
	}
	Server *http.Server
)

func Run() {

	network, localAddress := "tcp", Config.HttpAddr

	// If the port is zero, treat the address as a fully qualified local address.
	// This address must be prefixed with the network type followed by a colon,
	// e.g. unix:/tmp/app.socket or tcp6:::1 (equivalent to tcp6:0:0:0:0:0:0:0:1)
	if Config.HttpPort == 0 {
		parts := strings.SplitN(Config.HttpAddr, ":", 2)
		network = parts[0]
		localAddress = parts[1]
	} else {
		localAddress += ":" + strconv.Itoa(Config.HttpPort)
	}

	if network != "unix" && network != "tcp" {
		logger.Printf("fatal", "lessgo/pagelet: Unknown Network %s", network)
		return
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		logger.Printf("info", "lessgo/pagelet: Listening on port %d ...", Config.HttpPort)
	}()

	Server = &http.Server{
		Addr:           localAddress,
		Handler:        http.HandlerFunc(handle),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if network == "unix" {
		// TODO already in use
		os.Remove(localAddress)
	}

	listener, err := net.Listen(network, localAddress)
	if err != nil {
		logger.Printf("fatal", "lessgo/pagelet: net.Listen error %v", err)
		return
	}

	if err := Server.Serve(listener); err != nil {
		logger.Printf("fatal", "lessgo/pagelet: Server.Serve error %v", err)
	}
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
