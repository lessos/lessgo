package pagelet

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
	"time"
)

type View interface {
	Render(req *Request, resp *Response)
}

var (
	tplLocker     sync.Mutex
	templateFuncs = map[string]interface{}{
		"eq": Equal,
		// Skips sanitation on the parameter.  Do not use with dynamic data.
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},
		// Returns a copy of the string s with the old replaced by new
		"replace": func(s, old, new string) string {
			return strings.Replace(s, old, new, -1)
		},
		// Format a date according to the application's default date(time) format.
		"date": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"datetime": func(t time.Time) string {
			//t, _ := time.Parse("2006-01-02 15:04:05.000 -0700", fmttime)
			return t.Format("2006-01-02 15:04")
		},
		// upper returns a copy of the string s with all Unicode letters mapped to their upper case
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		// lower returns a copy of the string s with all Unicode letters mapped to their lower case
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		// Perform a message look-up for the given locale and message using the given arguments
		"T": func(lang map[string]interface{}, msg string, args ...interface{}) string {
			return i18nTranslate(lang["LANG"].(string), msg, args...)
		},
		// "set": func(renderArgs map[string]interface{}, key string, value interface{}) {
		// 	renderArgs[key] = value
		// },
	}
)

// This object handles loading and parsing of templates.
// Everything below the application's views directory is treated as a template.
type TemplateLoader struct {

	// Map from template name to the path from whence it was loaded.
	templatePaths map[string]string

	// This is the set of all templates under views
	templateSets map[string]*template.Template
}

type Template interface {
	Name() string
	Content() []string
	Render(wr io.Writer, arg interface{}) error
}

func (loader *TemplateLoader) initModule(module string) {

	tplLocker.Lock()
	defer tplLocker.Unlock()

	var cfgMod *ConfigModule
	for _, v := range Config.Module {

		if v.Name == module {
			cfgMod = &v
			break
		}
	}
	if cfgMod == nil {
		return
	}

	//
	loaderTemplateSet, _ := loader.templateSets[cfgMod.Name]

	for _, baseDir := range cfgMod.ViewPaths {

		_ = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			var fileStr string

			addTemplate := func(templateName string) (err error) {

				if _, ok := loader.templatePaths[cfgMod.Name+templateName]; ok {
					return nil
				}

				loader.templatePaths[cfgMod.Name+templateName] = path

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

						loaderTemplateSet = template.New(templateName).Funcs(templateFuncs)

						_, err = loaderTemplateSet.Parse(fileStr)

						loader.templateSets[cfgMod.Name] = loaderTemplateSet
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

func (loader *TemplateLoader) Template(module, name string) (Template, error) {

	set, ok := loader.templateSets[module]
	if !ok {
		loader.initModule(module)
	}

	set, ok = loader.templateSets[module]
	if !ok {
		return nil, fmt.Errorf("Template %s:%s not found.", module, name)
	}

	//Println("loader.templateSet", module, name)
	// This is necessary.
	// If a nil loader.compileError is returned directly, a caller testing against
	// nil will get the wrong result.  Something to do with casting *Error to error.
	var err error
	//if loader.compileError != nil {
	//	err = loader.compileError
	//}

	tmpl := set.Lookup(name)
	if tmpl == nil && err == nil {
		return nil, fmt.Errorf("Template %s not found.", name)
	}

	return goTemplate{tmpl, loader}, err
}

// Adapter for Go Templates.
type goTemplate struct {
	*template.Template
	loader *TemplateLoader
}

// return a 'revel.Template' from Go's template.
func (gotmpl goTemplate) Render(wr io.Writer, arg interface{}) error {
	return gotmpl.Execute(wr, arg)
}

func (gotmpl goTemplate) Content() []string {

	bytes, err := ioutil.ReadFile(gotmpl.loader.templatePaths[gotmpl.Name()])
	if err != nil {
		return []string{}
	}

	return strings.Split(string(bytes), "\n")
}
