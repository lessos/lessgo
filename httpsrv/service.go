// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpsrv

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lessos/lessgo/deps/go.net/websocket"
	"github.com/lessos/lessgo/logger"
)

type Service struct {
	Config  Config
	Filters []Filter

	err            error
	modules        []Module
	server         *http.Server
	rpcRegs        map[string]*rpc.Server
	handlers       []reg_handler
	TemplateLoader *TemplateLoader
}

type reg_handler struct {
	base    string
	handler http.Handler
}

var (
	lock          sync.Mutex
	GlobalService = NewService()
)

func NewService() Service {

	return Service{

		Config: Config{
			HttpAddr:         "0.0.0.0",
			HttpPort:         8080,
			HttpTimeout:      30, // 30 seconds
			CookieKeyLocale:  "lang",
			CookieKeySession: "access_token",
		},

		Filters: DefaultFilters,

		modules: []Module{},

		rpcRegs:  map[string]*rpc.Server{},
		handlers: []reg_handler{},

		TemplateLoader: &TemplateLoader{
			templatePaths: map[string]string{},
			templateSets:  map[string]*template.Template{},
		},
	}
}

func (s *Service) HandlerRegisterPrev(baseuri string, h http.Handler) {

	baseuri = "/" + strings.Trim(filepath.Clean(baseuri), "/")

	for _, v := range s.handlers {
		if v.base == baseuri {
			return
		}
	}

	s.handlers = append(s.handlers, reg_handler{
		base:    baseuri,
		handler: h,
	})
}

func (s *Service) ModuleRegister(baseuri string, mod Module) {

	lock.Lock()
	defer lock.Unlock()

	set := Module{
		name:        mod.name,
		baseuri:     strings.Trim(baseuri, "/"),
		viewpaths:   mod.viewpaths,
		controllers: mod.controllers,
	}

	mod.routes = append(mod.routes, defaultRoute)

	for _, r := range mod.routes {

		if r.Type == RouteTypeStatic && r.StaticPath != "" {

			set.routes = append(set.routes, r)

		} else if r.Type == RouteTypeBasic {

			r.Path = strings.Trim(r.Path, "/")
			r.tree = strings.Split(r.Path, "/")
			r.treelen = len(r.tree)

			if r.treelen < 1 {
				continue
			}

			set.routes = append(set.routes, r)
		}
	}

	s.TemplateLoader.Set(mod.name, mod.viewpaths)

	s.modules = append(s.modules, set)
}

func (s *Service) Error() error {
	return s.err
}

func (s *Service) Start() error {

	network, localAddress := "tcp", s.Config.HttpAddr

	// If the port is zero, treat the address as a fully qualified local address.
	// This address must be prefixed with the network type followed by a colon,
	// e.g. unix:/tmp/app.socket or tcp6:::1 (equivalent to tcp6:0:0:0:0:0:0:0:1)
	if s.Config.HttpPort == 0 || strings.HasPrefix(s.Config.HttpAddr, "unix:") {
		parts := strings.SplitN(s.Config.HttpAddr, ":", 2)
		if len(parts) > 0 {
			network = parts[0]
		}
		if len(parts) > 1 {
			localAddress = parts[1]
		}
	} else {
		localAddress += fmt.Sprintf(":%d", s.Config.HttpPort)
	}

	if network != "unix" && network != "tcp" {
		logger.Printf("fatal", "lessgo/httpsrv: Unknown Network %s", network)
		return nil
	}

	//
	if network == "unix" {
		// TODO already in use
		os.Remove(localAddress)
	}

	//
	if s.Config.HttpTimeout < 3 {
		s.Config.HttpTimeout = 10
	}

	//
	srvmux := http.NewServeMux()

	//
	for rpcpath, rpcsrv := range s.rpcRegs {
		srvmux.Handle(rpcpath, rpcsrv)
	}

	//
	for _, v := range s.handlers {
		srvmux.Handle(v.base, v.handler)
	}

	srvmux.HandleFunc("/", s.handle)

	//
	s.server = &http.Server{
		Addr:    localAddress,
		Handler: srvmux,
		// ReadTimeout:    time.Duration(s.Config.HttpTimeout) * time.Second,
		// WriteTimeout:   time.Duration(s.Config.HttpTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//
	listener, err := net.Listen(network, localAddress)
	if err != nil {
		logger.Printf("fatal", "lessgo/httpsrv: net.Listen error %v", err)
		s.err = err
		return nil
	}
	fmt.Println("lessgo/httpsrv: Listening on", localAddress)

	if network == "unix" {
		os.Chmod(localAddress, 0770)
	}

	//
	if err := s.server.Serve(listener); err != nil {
		logger.Printf("fatal", "lessgo/httpsrv: server.Serve error %v", err)
		s.err = err
	} else {
		time.Sleep(100 * time.Millisecond)
		logger.Printf("info", "lessgo/httpsrv: Listening on %s ...", localAddress)
	}

	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) handle(w http.ResponseWriter, r *http.Request) {

	upgrade := r.Header.Get("Upgrade")

	if upgrade == "websocket" || upgrade == "Websocket" {

		websocket.Handler(func(ws *websocket.Conn) {
			r.Method = "WS"
			s.handleInternal(w, r, ws)
		}).ServeHTTP(w, r)

	} else {

		s.handleInternal(w, r, nil)
	}
}

func (s *Service) handleInternal(w http.ResponseWriter, r *http.Request, ws *websocket.Conn) {

	defer func() {

		// if err := recover(); err != nil {
		// 	logger.Printf("error", "handleInternal Panic on %s", err)
		// }

		r.Body.Close()
	}()

	var (
		req  = NewRequest(r)
		resp = NewResponse(w)
		c    = NewController(s, req, resp)
	)

	if ws != nil {
		req.WebSocket = WebSocket{
			conn: ws,
		}
	}

	for _, filter := range s.Filters {
		filter(c)
	}
}
