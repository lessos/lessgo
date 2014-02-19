package pagelet

import (
    "strings"
)

type ConfigBase struct {
    UrlBasePath string
    HttpPort    int
    Module      []ConfigModule
}

type ConfigModule struct {
    Name      string
    ViewPaths []string
    Routes    []Route
}

func (c *ConfigBase) moduleInit(module string) {

    for _, v := range c.Module {

        if v.Name == module {
            return
        }
    }

    m := ConfigModule{
        Name:      module,
        ViewPaths: []string{},
        Routes:    []Route{},
    }

    c.Module = append(c.Module, m)
}

func (c *ConfigBase) RouteStaticAppend(module, path, pathto string) {

    c.moduleInit(module)

    for i, v := range c.Module {

        if v.Name != module {
            continue
        }

        route := Route{
            Type: "static",
            Path: strings.Trim(path, "/"),
            Tree: []string{pathto},
        }
        //Println(route)
        c.Module[i].Routes = append(v.Routes, route)

        break
    }
}

func (c *ConfigBase) RouteAppend(module, path string) {

    c.moduleInit(module)

    path = strings.Trim(path, "/")
    tree := strings.Split(path, "/")
    if len(tree) < 1 {
        return
    }

    route := Route{
        Type:    "std",
        Path:    path,
        Tree:    tree,
        TreeLen: len(tree),
    }

    for i, v := range c.Module {

        if v.Name != module {
            continue
        }

        c.Module[i].Routes = append(v.Routes, route)

        break
    }
}

func (c *ConfigBase) ViewPath(module, path string) {

    c.moduleInit(module)

    for i, v := range c.Module {

        if v.Name != module {
            continue
        }

        gotPath := false
        for _, v2 := range v.ViewPaths {

            if v2 != path {
                continue
            }

            gotPath = true
        }

        if !gotPath {
            c.Module[i].ViewPaths = append(v.ViewPaths, path)
        }

        break
    }
}
