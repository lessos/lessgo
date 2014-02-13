package pagelet

import (
    "reflect"
)

type Router struct {
    StaticDir map[string]string
}

type Route struct {
    Path           string
    ControllerName string
    MethodName     string
    Type           string
}

type RouteMatch struct {
    ControllerName string // e.g. App
    MethodName     string // e.g. Login
    //Params         map[string][]string // e.g. {id: 123}
}

func NewRouter() *Router {

    var r Router

    r.StaticDir = map[string]string{
        "/static": "static",
    }

    return &r
}

func (r *Router) SetStatic(route, path string) {
    r.StaticDir[route] = path
    Println("SetStatic", route, path)
}

func (r *Router) SetRoute(route, action string) {

}

func RouterFilter(c *Controller, fc []Filter) {

    defer func() {
        fc[0](c, fc[1:])
    }()

    c.Name = "Index"
    c.MethodName = "Index"

    /* var route *RouteMatch = MainRouter.Route(c.Request.Request)
       if route == nil {
           c.Result = c.NotFound("No matching route found")
           return
       } */

    /* for prefix, staticDir := range this.StaticDir {

        //fmt.Println(prefix, staticDir)

        if r.URL.Path == "/favicon.ico" {
            file := staticDir + r.URL.Path
            http.ServeFile(w, r, file)
            return
        }

        if strings.HasPrefix(r.URL.Path, prefix) {

            file := staticDir + r.URL.Path[len(prefix):]
            finfo, err := os.Stat(file)

            if err != nil {
                http.NotFound(w, r)
                return
            }

            // if the request is dir
            if finfo.IsDir() {
                http.NotFound(w, r)
                return
            }

            http.ServeFile(w, r, file)
            return
        }

        // TODO Compress
    } */

    ctrl, ok := controllers[c.Name]
    if !ok {
        return // TODO
    }

    var (
        appControllerPtr = reflect.New(ctrl.Type)
        appController    = appControllerPtr.Elem()
        cValue           = reflect.ValueOf(c)
    )

    for _, index := range ctrl.ControllerIndexes {
        appController.FieldByIndex(index).Set(cValue)
    }

    c.AppController = appControllerPtr.Interface()
}
