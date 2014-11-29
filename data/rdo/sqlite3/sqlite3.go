package sqlite3

import (
	"database/sql"
	"errors"
	"github.com/lessos/lessgo/data/rdo/base"
	//"fmt"
	//_ "code.google.com/p/gosqlite/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func NewClient(c base.Config) (*base.Client, error) {

	client := &base.Client{}

	if c.Socket == "" {
		return client, errors.New("Incorrect configuration")
	}

	c.Dsn = c.Socket

	db, err := sql.Open(c.Driver, c.Dsn)
	if err != nil {
		return client, err
	}

	client.Base, _ = base.BaseInit(c, db)
	for k, v := range sqlite3Stmt {
		client.Base.BaseStmt[k] = v
	}

	client.Dialect = &sqlite3Dialect{}
	err = client.Dialect.Init(client.Base)

	return client, err
}
