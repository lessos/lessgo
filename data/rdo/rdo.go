package rdo

import (
    "./base"
    "./mysql"
    "errors"
)

var (
    clients = map[string]*base.Client{}
)

func NewClient(name string, c base.Config) (dc *base.Client, err error) {

    if _, ok := clients[name]; ok {
        return clients[name], nil
    }

    switch c.Driver {
    case "mysql":
        dc, err = mysql.NewClient(c)
    default:
        err = errors.New("No Driver Found")
    }

    if err != nil {
        return dc, err
    }

    clients[name] = dc

    return dc, nil
}
