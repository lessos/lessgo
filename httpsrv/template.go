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
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var tlock sync.Mutex

// This object handles loading and parsing of templates.
// Everything below the application's views directory is treated as a template.
type TemplateLoader struct {
	// Map from template name to the path from whence it was loaded.
	templatePaths map[string]string

	// This is the set of all templates under views
	templateSets map[string]*template.Template
}

type iTemplate interface {
	Render(wr io.Writer, arg interface{}) error
}

func (loader *TemplateLoader) Clean(modname string) {

	tlock.Lock()
	defer tlock.Unlock()

	if _, ok := loader.templateSets[modname]; ok {
		delete(loader.templateSets, modname)
	}

	for k := range loader.templatePaths {

		if strings.HasPrefix(k, modname+".") {
			delete(loader.templatePaths, k)
		}
	}
}

func (loader *TemplateLoader) Set(modname string, viewpaths []string) {

	tlock.Lock()
	defer tlock.Unlock()

	loaderTemplateSet, ok := loader.templateSets[modname]
	if ok {
		return
	}

	for _, baseDir := range viewpaths {

		_ = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			var fileStr string

			addTemplate := func(templateName string) (err error) {

				if _, ok := loader.templatePaths[modname+"."+templateName]; ok {
					return nil
				}

				loader.templatePaths[modname+"."+templateName] = path

				// Load the file if we haven't already
				if fileStr == "" {

					fileBytes, err := ioutil.ReadFile(path)
					if err != nil {
						return nil
					}

					fileStr = string(fileBytes)
				}

				if loaderTemplateSet == nil {

					var funcError error

					func() {

						defer func() {
							if err := recover(); err != nil {
								funcError = errors.New("Panic (Template Loader)")
							}
						}()

						loaderTemplateSet = template.New(templateName).Funcs(TemplateFuncs)

						if _, err := loaderTemplateSet.Parse(fileStr); err == nil {
							loader.templateSets[modname] = loaderTemplateSet
						}
					}()

					if funcError != nil {
						return funcError
					}

				} else {

					_, err = loaderTemplateSet.New(templateName).Parse(fileStr)
				}

				return err
			}

			templateName := path[len(baseDir)+1:]

			// Lower case the file name for case-insensitive matching
			lowerCaseTemplateName := strings.ToLower(templateName)

			_ = addTemplate(templateName)
			_ = addTemplate(lowerCaseTemplateName)

			return nil
		})
	}
}

func (loader *TemplateLoader) Template(modname, tplname string) (iTemplate, error) {

	set, ok := loader.templateSets[modname]
	if !ok || set == nil {
		return nil, fmt.Errorf("Template %s not found.", tplname)
	}

	tmpl := set.Lookup(tplname)
	if tmpl == nil {
		return nil, fmt.Errorf("Template %s:%s not found.", modname, tplname)
	}

	return goTemplate{tmpl, loader}, nil
}

// Adapter for Go Templates.
type goTemplate struct {
	*template.Template
	loader *TemplateLoader
}

// return a 'httpsrv.goTemplate' from Go's template.
func (gotmpl goTemplate) Render(wr io.Writer, arg interface{}) error {

	defer func() {
		if err := recover(); err != nil {
			//
		}
	}()

	return gotmpl.Execute(wr, arg)
}
