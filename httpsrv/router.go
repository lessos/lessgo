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
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	RouteTypeBasic  = "basic"
	RouteTypeStatic = "static"
)

var (
	defaultRoute = Route{
		Type: RouteTypeBasic,
		Path: ":controller/:action",
	}
)

type Route struct {
	Type       string
	Path       string
	StaticPath string
	Params     map[string]string // e.g. {id: 123}
	tree       []string
	treelen    int
}

func RouterFilter(c *Controller) {

	urlpath := strings.Trim(filepath.Clean(c.Request.URL.Path), "/")

	if c.service.Config.UrlBasePath != "" {
		urlpath = strings.TrimPrefix(strings.TrimPrefix(urlpath, c.service.Config.UrlBasePath), "/")
	}

	if urlpath == "favicon.ico" {
		return
	}

	for _, mod := range c.service.modules {

		if !strings.HasPrefix(urlpath, mod.baseuri) && mod.name != "default" {
			continue
		}

		urlpath = strings.TrimPrefix(strings.TrimPrefix(urlpath, mod.baseuri), "/")

		var (
			rt    = strings.Split(urlpath, "/")
			rtlen = len(rt)
		)

		if rtlen == 1 {
			rt = append(rt, "Index")
			rtlen++
		}

		for _, route := range mod.routes {

			if route.Type == RouteTypeStatic && strings.HasPrefix(urlpath, route.Path) {

				file := route.StaticPath + "/" + urlpath[len(route.Path):]
				finfo, err := os.Stat(file)

				if err != nil {
					http.NotFound(c.Response.Out, c.Request.Request)
					return
				}

				if finfo.IsDir() {
					http.NotFound(c.Response.Out, c.Request.Request)
					return
				}

				http.ServeFile(c.Response.Out, c.Request.Request, file)
				return
			}

			// TODO
			if route.Type != RouteTypeBasic {
				continue
			}

			if rtlen < route.treelen {
				continue
			}

			matlen, ctrlName, actionName, params := 0, "", "", map[string]string{}

			for i, node := range route.tree {

				if node[0:1] == ":" {

					switch node[1:] {

					case "controller":
						ctrlName = rt[i]

					case "action":
						actionName = rt[i]

					default:
						params[node[1:]] = rt[i]
					}

					matlen++

				} else if node == rt[i] {

					matlen++
				}
			}

			if matlen == route.treelen {

				if len(ctrlName) > 0 {
					c.Name = strings.Replace(strings.Title(ctrlName), "-", "", -1)
				} else if val, ok := route.Params["controller"]; ok {
					c.Name = strings.Replace(strings.Title(val), "-", "", -1)
				}

				if len(actionName) > 0 {
					c.ActionName = strings.Replace(strings.Title(actionName), "-", "", -1)
				} else if val, ok := route.Params["action"]; ok {
					c.ActionName = strings.Replace(strings.Title(val), "-", "", -1)
				}

				// TODO sec
				for k, v := range route.Params {
					c.Params.Values[k] = append(c.Params.Values[k], v)
				}

				for k, v := range params {
					c.Params.Values[k] = append(c.Params.Values[k], v)
				}

				break
			}
		}

		c.mod_name = mod.name
		c.mod_urlbase = mod.baseuri

		ctrl, ok := mod.controllers[c.Name]
		if !ok {

			c.Name = "Index"

			if ctrl, ok = mod.controllers[c.Name]; !ok {

				c.mod_name = "default"

				if ctrl, ok = mod.controllers[c.Name]; !ok {
					return
				}
			}
		}

		var (
			appControllerPtr = reflect.New(ctrl.Type)
			appController    = appControllerPtr.Elem()
			cValue           = reflect.ValueOf(c)
		)

		for _, index := range ctrl.ControllerIndexes {
			appController.FieldByIndex(index).Set(cValue)
		}

		if mod.baseuri != "" {
			c.Request.RequestPath = mod.baseuri + "/" + urlpath
		} else {
			c.Request.RequestPath = urlpath
		}
		c.appController = appControllerPtr.Interface()

		break
	}
}
