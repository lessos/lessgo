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
)

type View interface {
    Render(req *Request, resp *Response)
}

var (
    TemplateFuncs = map[string]interface{}{}
)

// This object handles loading and parsing of templates.
// Everything below the application's views directory is treated as a template.
type TemplateLoader struct {
    // This is the set of all templates under views
    templateSet *template.Template
    // If an error was encountered parsing the templates, it is stored here.
    //compileError *Error
    // Paths to search for templates, in priority order.
    paths []string
    // Map from template name to the path from whence it was loaded.
    templatePaths map[string]string
}

type Template interface {
    Name() string
    Content() []string
    Render(wr io.Writer, arg interface{}) error
}

func NewTemplateLoader(paths []string) *TemplateLoader {

    loader := &TemplateLoader{
        paths:         paths,
        templatePaths: map[string]string{},
        templateSet:   nil,
    }
    //fmt.Println("NewTemplateLoader", loader)

    var splitDelims []string
    ViewsPath := ""
    //var templateSet *template.Template = nil

    for _, baseDir := range loader.paths {

        _ = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {

            if err != nil {
                return nil
            }

            if info.IsDir() {
                return nil
            }

            fmt.Println(info.Name())

            var fileStr string

            addTemplate := func(templateName string) (err error) {

                if _, ok := loader.templatePaths[templateName]; ok {
                    return nil
                }

                loader.templatePaths[templateName] = path

                // Load the file if we haven't already
                if fileStr == "" {

                    fileBytes, err := ioutil.ReadFile(path)
                    if err != nil {
                        return nil
                    }

                    fileStr = string(fileBytes)
                }

                if loader.templateSet == nil {

                    var funcError error

                    func() {

                        defer func() {
                            if err := recover(); err != nil {
                                funcError = errors.New("Panic (Template Loader)")
                            }
                        }()

                        loader.templateSet = template.New(templateName).Funcs(TemplateFuncs)
                        // If alternate delimiters set for the project, change them for this set
                        if splitDelims != nil && baseDir == ViewsPath {
                            loader.templateSet.Delims(splitDelims[0], splitDelims[1])
                        } else {
                            // Reset to default otherwise
                            loader.templateSet.Delims("", "")
                        }

                        //fmt.Println("fileStr", templateName)
                        _, err = loader.templateSet.Parse(fileStr)
                    }()

                    if funcError != nil {
                        return funcError
                    }

                } else {

                    if splitDelims != nil && baseDir == ViewsPath {
                        loader.templateSet.Delims(splitDelims[0], splitDelims[1])
                    } else {
                        loader.templateSet.Delims("", "")
                    }
                    _, err = loader.templateSet.New(templateName).Parse(fileStr)

                }

                return err
            }

            templateName := path[len(baseDir)+1:]

            // Lower case the file name for case-insensitive matching
            lowerCaseTemplateName := strings.ToLower(templateName)

            //fmt.Println("templateName", path, baseDir, templateName)

            _ = addTemplate(templateName)

            _ = addTemplate(lowerCaseTemplateName)

            return nil
        })

        //fmt.Println(funcErr)
    }

    //fmt.Println("loader.templateSet", loader.templateSet)

    return loader
}

func (loader *TemplateLoader) Template(name string) (Template, error) {

    //fmt.Println(tmpl)
    // This is necessary.
    // If a nil loader.compileError is returned directly, a caller testing against
    // nil will get the wrong result.  Something to do with casting *Error to error.
    var err error
    //if loader.compileError != nil {
    //	err = loader.compileError
    //}

    tmpl := loader.templateSet.Lookup(name)
    //fmt.Println("loader.templateSet.Lookup", tmpl)
    if tmpl == nil && err == nil {
        return nil, fmt.Errorf("Template %s not found.", name)
    }

    return GoTemplate{tmpl, loader}, err
}

// Adapter for Go Templates.
type GoTemplate struct {
    *template.Template
    loader *TemplateLoader
}

// return a 'revel.Template' from Go's template.
func (gotmpl GoTemplate) Render(wr io.Writer, arg interface{}) error {
    return gotmpl.Execute(wr, arg)
}

func (gotmpl GoTemplate) Content() []string {
    content, _ := ReadLines(gotmpl.loader.templatePaths[gotmpl.Name()])
    return content
}
