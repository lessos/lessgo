package rdc

import (
	"database/sql"
)

type Conn struct {
	driver string
	db     *sql.DB
	cfg    Config
	dsn    string
}

func (cn *Conn) Close() {
	cn.db.Close()
}
