package pagelet

import (
    "net/http"
    "os"
    "path/filepath"
    "reflect"
    "strings"
)

type Router struct {
    Routes  []Route
    Modules map[string][]Route
}

type Route struct {
    Type           string
    Path           string
    Tree           []string
    TreeLen        int
    ControllerName string
    MethodName     string
}

type RouteMatch struct {
    ControllerName string            // e.g. App
    MethodName     string            // e.g. Login
    Params         map[string]string // e.g. {id: 123}
}

func (r *Router) RouteStaticAppend(path, pathto string) {

    route := Route{
        Type: "static",
        Path: strings.Trim(path, "/"),
        Tree: []string{pathto},
    }

    r.Routes = append(r.Routes, route)
}

func (r *Router) RouteAppend(path, action string) {

    actions := strings.Split(action, ".")

    if len(actions) != 2 {
        return
    }

    tree := strings.Split(strings.Trim(path, "/"), "/")
    if len(tree) < 1 {
        return
    }

    route := Route{
        Type:           "std",
        Path:           path,
        Tree:           tree,
        TreeLen:        len(tree),
        ControllerName: actions[0],
        MethodName:     actions[1],
    }

    r.Routes = append(r.Routes, route)
}

func RouterFilter(c *Controller, fc []Filter) {

    defer func() {
        fc[0](c, fc[1:])
    }()

    urlpath := strings.Trim(filepath.Clean(c.Request.URL.Path), "/")

    if Config.UrlBasePath != "" {
        urlpath = strings.TrimLeft(strings.TrimLeft(urlpath, Config.UrlBasePath), "/")
    }

    if urlpath == "favicon.ico" {
        return
    }

    rt := strings.Split(urlpath, "/")
    rtlen := len(rt)

    for _, route := range MainRouter.Routes {

        if route.Type == "static" && strings.HasPrefix(urlpath, route.Path) {

            file := route.Tree[0] + "/" + urlpath[len(route.Path):]
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
        if route.Type != "std" {
            continue
        }

        if rtlen < route.TreeLen {
            continue
        }

        matRoute := 0
        ctrlName := ""
        methodName := ""

        for i, node := range route.Tree {

            if node == ":controller" {
                ctrlName = rt[i]
                matRoute++
                continue
            }

            if node == ":action" {
                methodName = rt[i]
                matRoute++
                continue
            }

            if node == rt[i] {
                matRoute++
            }
        }

        if matRoute == route.TreeLen {

            if len(ctrlName) > 0 {
                c.Name = strings.Replace(strings.Title(ctrlName), "-", "", -1)
            } else {
                c.Name = "Index"
            }

            if len(methodName) > 0 {
                c.MethodName = strings.Replace(strings.Title(methodName), "-", "", -1)

            } else {
                c.MethodName = "Index"
            }

            break
        }
    }

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
