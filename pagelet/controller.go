package pagelet

import (
    "io"
    "net/http"
    "reflect"
)

type Controller struct {
    Name          string // The controller name, e.g. "App"
    MethodName    string // The method name, e.g. "Index"
    ModuleName    string
    Request       *Request
    Response      *Response
    AppController interface{} // The controller that was instantiated.
    AutoRender    bool
    ViewData      map[string]interface{}
}

type ControllerType struct {
    Type              reflect.Type
    Methods           []string
    ControllerIndexes [][]int
}

var (
    controllerPtrType = reflect.TypeOf(&Controller{})
)

var controllers = make(map[string]*ControllerType)

func NewController(req *Request, resp *Response) *Controller {

    return &Controller{
        Request:    req,
        Response:   resp,
        AutoRender: true,
    }
}

func ActionInvoker(c *Controller, _ []Filter) {

    //
    if c.AppController == nil {
        return
    }
    //Println("AAA")

    execController := reflect.ValueOf(c.AppController).MethodByName(c.MethodName)

    args := []reflect.Value{}
    if execController.Type().IsVariadic() {
        execController.CallSlice(args)
    } else {
        execController.Call(args)
    }

    if c.AutoRender {
        c.AutoRender = false
        c.Render()
    }

    //Println("ActionInvoker DONE")
}

func (c *Controller) Render(args ...interface{}) {

    templatePath := c.Name + "/" + c.MethodName + ".tpl"
    if len(args) == 1 && reflect.TypeOf(args[0]).Kind() == reflect.String {
        templatePath = args[1].(string)
    }

    //Println(c.ModuleName, templatePath)
    // Handle panics when rendering templates.
    defer func() {
        if err := recover(); err != nil {

        }
    }()

    template, err := MainTemplateLoader.Template(c.ModuleName, templatePath)
    if err != nil {
        return //c.RenderError(err)
    }

    // If it's a HEAD request, throw away the bytes.
    out := io.Writer(c.Response.Out)

    c.Response.WriteHeader(http.StatusOK, "text/html; charset=utf-8")

    if err = template.Render(out, c.ViewData); err != nil {
        //
    }
}

func RegisterController(module string, c interface{}, methods []string) {

    v := reflect.ValueOf(c)
    if !v.IsValid() {
        return
    }

    t := reflect.TypeOf(c)
    elem := t.Elem()

    cm := &ControllerType{
        Type:              elem,
        Methods:           []string{},
        ControllerIndexes: findControllers(elem),
    }

    for _, method := range methods {

        if m := v.MethodByName(method); m.IsValid() {
            cm.Methods = append(cm.Methods, method)
        }
    }

    controllers[module+elem.Name()] = cm
}

func findControllers(appControllerType reflect.Type) (indexes [][]int) {

    // It might be a multi-level embedding. To find the controllers, we follow
    // every anonymous field, using breadth-first search.
    type nodeType struct {
        val   reflect.Value
        index []int
    }

    appControllerPtr := reflect.New(appControllerType)
    queue := []nodeType{{appControllerPtr, []int{}}}

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

            queue = append(queue,
                nodeType{fieldValue, append(append([]int{}, node.index...), i)})
        }
    }

    return
}
