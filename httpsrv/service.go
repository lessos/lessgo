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
	"os"
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
	templateLoader *templateLoader
}

var (
	lock           sync.Mutex
	DefaultService = Service{

		Config: Config{
			HttpPort:         0,
			CookieKeyLocale:  "lang",
			CookieKeySession: "access_token",
		},

		Filters: DefaultFilters,

		modules: []Module{},

		templateLoader: &templateLoader{
			templatePaths: map[string]string{},
			templateSets:  map[string]*template.Template{},
		},
	}
)

func NewService() Service {
	return DefaultService
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

	for _, r := range mod.routes {

		if r.Type == RouteTypeStatic && len(r.Tree) > 0 {

			set.routes = append(set.routes, r)

		} else if r.Type == RouteTypeBasic {

			r.Path = strings.Trim(r.Path, "/")
			r.Tree = strings.Split(r.Path, "/")
			r.treelen = len(r.Tree)

			if r.treelen < 1 {
				continue
			}

			set.routes = append(set.routes, r)
		}
	}

	s.templateLoader.init(mod)

	s.modules = append(s.modules, set)
}

func (s *Service) Err() error {
	return s.err
}

func (s *Service) Start() error {

	network, localAddress := "tcp", s.Config.HttpAddr

	// If the port is zero, treat the address as a fully qualified local address.
	// This address must be prefixed with the network type followed by a colon,
	// e.g. unix:/tmp/app.socket or tcp6:::1 (equivalent to tcp6:0:0:0:0:0:0:0:1)
	if s.Config.HttpPort == 0 {
		parts := strings.SplitN(s.Config.HttpAddr, ":", 2)
		network = parts[0]
		localAddress = parts[1]
	} else {
		localAddress += fmt.Sprintf(":%d", s.Config.HttpPort)
	}

	if network != "unix" && network != "tcp" {
		logger.Printf("fatal", "lessgo/httpsrv: Unknown Network %s", network)
		return nil
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		logger.Printf("info", "lessgo/httpsrv: Listening on %s ...", localAddress)
	}()

	s.server = &http.Server{
		Addr:           localAddress,
		Handler:        http.HandlerFunc(s.handle),
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
		logger.Printf("fatal", "lessgo/httpsrv: net.Listen error %v", err)
		return nil
	}

	if err := s.server.Serve(listener); err != nil {
		logger.Printf("fatal", "lessgo/httpsrv: server.Serve error %v", err)
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

		if err := recover(); err != nil {
			logger.Printf("error", "handleInternal Panic on %s", err)
		}

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
