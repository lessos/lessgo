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

func RouterFilter(c *Controller) {

    urlpath := strings.Trim(filepath.Clean(c.Request.URL.Path), "/")

    if Config.UrlBasePath != "" {
        urlpath = strings.TrimPrefix(strings.TrimPrefix(urlpath, Config.UrlBasePath), "/")
    }

    if urlpath == "favicon.ico" {
        return
    }

    for _, mod := range Config.Module {

        if !strings.HasPrefix(urlpath, mod.Name) && mod.Name != "default" {
            continue
        }

        urlpath = strings.TrimPrefix(strings.TrimPrefix(urlpath, mod.Name), "/")
        rt := strings.Split(urlpath, "/")
        rtlen := len(rt)

        //println("Router MAT", mod)

        for _, route := range mod.Routes {

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

        c.ModuleName = mod.Name

        //println("Router controllers added")

        ctrl, ok := controllers[c.ModuleName+c.Name]
        if !ok {
            return // TODO
        }

        //println("Router controllers added", c)

        var (
            appControllerPtr = reflect.New(ctrl.Type)
            appController    = appControllerPtr.Elem()
            cValue           = reflect.ValueOf(c)
        )

        for _, index := range ctrl.ControllerIndexes {
            appController.FieldByIndex(index).Set(cValue)
        }

        c.AppController = appControllerPtr.Interface()

        break
    }

}
