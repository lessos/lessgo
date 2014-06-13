package rdc

import (
    "./mysql"
    "./sqlite3"
    "errors"
    "fmt"
)

var (
    databases = map[string]*Conn{}
)

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

func NewConnect(name string, c ConfigSource) (*Conn, error) {

    if _, ok := databases[name]; ok {
        return databases[name], nil
    }

    var dn Conn
    var err error

    dsn := ""
    if c.Driver == "sqlite3" {
        dsn = c.Socket
    } else if c.Host != "" {
        dsn = fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`,
            c.User, c.Pass, c.Host, c.Port, c.Dbname)
    } else if c.Socket != `` {
        dsn = fmt.Sprintf(`%s:%s@unix(%s)/%s?charset=%s`,
            c.User, c.Pass, c.Socket, c.Dbname, c.Charset)
    }

    switch c.Driver {
    case "mysql":
        dn.db, err = mysql.Open(c.Driver, dsn)
    case "sqlite3":
        dn.db, err = sqlite3.Open(c.Driver, dsn)
    default:
        err = errors.New("No Driver Found")
    }

    if err != nil {
        return &dn, nil
    }

    dn.driver = c.Driver
    dn.dsn = dsn

    databases[name] = &dn

    return &dn, nil
}
