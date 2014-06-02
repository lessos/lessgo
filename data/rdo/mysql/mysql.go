package mysql

import (
    "../base"
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func NewClient(c base.Config) (*base.Client, error) {

    client := &base.Client{}

    if c.Host != "" {
        c.Dsn = fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s`,
            c.User, c.Pass, c.Host, c.Port, c.Dbname)
    } else if c.Socket != `` {
        c.Dsn = fmt.Sprintf(`%s:%s@unix(%s)/%s?charset=%s`,
            c.User, c.Pass, c.Socket, c.Dbname, c.Charset)
    } else {
        return client, errors.New("Incorrect configuration")
    }

    db, err := sql.Open(c.Driver, c.Dsn)
    if err != nil {
        return client, err
    }

    client.Base, _ = base.BaseInit(c, db)
    for k, v := range mysqlStmt {
        client.Base.BaseStmt[k] = v
    }

    client.Dialect = &mysqlDialect{}
    err = client.Dialect.Init(client.Base)

    return client, err
}
