package rdc

import (
    "errors"
)

var databases = map[string]*Conn{}

func InstanceRegister(name string, cn *Conn) {
    if _, ok := databases[name]; ok {
        return
    }

    databases[name] = cn
}

func InstancePull(name string) (*Conn, error) {

    if v, ok := databases[name]; ok {
        return v, nil
    }

    return nil, errors.New("No Instance Found")
}
