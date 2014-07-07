package rdc

import (
//"database/sql"
)

type DataConn interface {

	//Data  *sql.DB

	// TryConnect verify can connect to the database
	TryConnect() error

	Test() bool
}

/*
func (cn *DataConn) Close() {
    //cn.Data.Close()
}
*/
