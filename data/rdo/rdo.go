package rdo

import (
    "./base"
    "./mysql"
    "./sqlite3"
    "errors"
)

var (
    clients = map[string]*base.Client{}
)

func ClientPull(name string) (dc *base.Client, err error) {

    if _, ok := clients[name]; ok {
        return clients[name], nil
    }

    return dc, errors.New("No Client Found")
}

func NewClient(name string, c base.Config) (dc *base.Client, err error) {

    if _, ok := clients[name]; ok {
        return clients[name], nil
    }

    switch c.Driver {
    case "mysql":
        dc, err = mysql.NewClient(c)
    case "sqlite3":
        dc, err = sqlite3.NewClient(c)
    default:
        err = errors.New("No Driver Found")
    }

    if err != nil {
        return dc, err
    }

    clients[name] = dc

    return dc, nil
}
