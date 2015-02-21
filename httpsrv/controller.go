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
	"io"
	"net/http"
	"reflect"

	"github.com/lessos/lessgo/utils"
)

type Controller struct {
	Name       string // The controller name, e.g. "App"
	ActionName string // The action name, e.g. "Index"
	Request    *Request
	Response   *Response
	Params     *Params // Parameters from URL and form (including multipart).
	Session    Session // Session, stored in cookie, signed.
	AutoRender bool
	Data       map[string]interface{}

	appController interface{} // The controller that was instantiated.
	module        string
	service       *Service
}

type controllerType struct {
	Type              reflect.Type
	Methods           []string
	ControllerIndexes [][]int
}

var (
	controllerPtrType = reflect.TypeOf(&Controller{})
)

func NewController(srv *Service, req *Request, resp *Response) *Controller {

	return &Controller{
		Name:       "Index",
		ActionName: "Index",
		service:    srv,
		Request:    req,
		Response:   resp,
		Params:     newParams(),
		AutoRender: true,
		Data:       map[string]interface{}{},
	}
}

func ActionInvoker(c *Controller) {

	//
	if c.appController == nil {
		return
	}

	execController := reflect.ValueOf(c.appController).MethodByName(c.ActionName + "Action")

	args := []reflect.Value{}
	if execController.Type().IsVariadic() {
		execController.CallSlice(args)
	} else {
		execController.Call(args)
	}

	if c.AutoRender {
		c.Render()
	}
}

func (c *Controller) Render(args ...interface{}) {

	c.AutoRender = false

	module, templatePath := c.module, c.Name+"/"+c.ActionName+".tpl"

	if len(args) == 2 &&
		reflect.TypeOf(args[0]).Kind() == reflect.String &&
		reflect.TypeOf(args[1]).Kind() == reflect.String {

		module, templatePath = args[0].(string), args[1].(string)

	} else if len(args) == 1 &&
		reflect.TypeOf(args[0]).Kind() == reflect.String {

		templatePath = args[0].(string)
	}

	// Handle panics when rendering templates.
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	template, err := c.service.templateLoader.Template(module, templatePath)
	if err != nil {
		return //c.RenderError(err)
	}

	// If it's a HEAD request, throw away the bytes.
	out := io.Writer(c.Response.Out)

	c.Response.WriteHeader(http.StatusOK, "text/html; charset=utf-8")

	if err = template.Render(out, c.Data); err != nil {
		println(err)
	}
}

func (c *Controller) RenderError(status int, msg string) {
	c.AutoRender = false
	c.Response.WriteHeader(status, "text/html; charset=utf-8")
	io.WriteString(c.Response.Out, msg)
}

func (c *Controller) UrlRedirect(url string) {
	c.AutoRender = false
	c.Response.Out.Header().Set("Location", url)
	c.Response.Out.WriteHeader(http.StatusFound)
}

func (c *Controller) RenderJSON(obj interface{}) {

	c.AutoRender = false

	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Set("Content-type", "application/json")

	if js, err := utils.JsonEncode(obj); err == nil {
		io.WriteString(c.Response.Out, js)
	}
}

func (c *Controller) Translate(msg string, args ...interface{}) string {
	return i18nTranslate(c.Request.Locale, msg, args...)
}

func findControllers(appControllerType reflect.Type) (indexes [][]int) {

	// It might be a multi-level embedding. To find the controllers, we follow
	// every anonymous field, using breadth-first search.
	type nodeType struct {
		val   reflect.Value
		index []int
	}

	var (
		appControllerPtr = reflect.New(appControllerType)
		queue            = []nodeType{{appControllerPtr, []int{}}}
	)

	for len(queue) > 0 {
		// Get the next value and de-reference it if necessary.
		var (
			node     = queue[0]
			elem     = node.val
			elemType = elem.Type()
		)

		if elemType.Kind() == reflect.Ptr {
			elem = elem.Elem()
			elemType = elem.Type()
		}

		queue = queue[1:]

		// Look at all the struct fields.
		for i := 0; i < elem.NumField(); i++ {
			// If this is not an anonymous field, skip it.
			structField := elemType.Field(i)
			if !structField.Anonymous {
				continue
			}

			fieldValue := elem.Field(i)
			fieldType := structField.Type

			// If it's a Controller, record the field indexes to get here.
			if fieldType == controllerPtrType {
				indexes = append(indexes, append(node.index, i))
				continue
			}

			queue = append(queue, nodeType{fieldValue,
				append(append([]int{}, node.index...), i)})
		}
	}

	return
}
