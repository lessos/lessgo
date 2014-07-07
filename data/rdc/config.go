package rdc

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	//"fmt"
	//"./mysql"
)

type Result sql.Result

type Config struct {
	Driver string
	DbPath string
}

type ConfigSource struct {
	// Database driver
	Driver string `json:"driver"`

	// Database server hostname or IP. Leave blank if using unix sockets.
	Host string `json:"host"`

	// Database server port. Leave blank if using unix sockets.
	Port string `json:"port"`

	// Username for authentication.
	User string `json:"user"`

	// Password for authentication.
	Pass string `json:"pass"`

	// A path of a UNIX socket file. Leave blank if using host and port.
	Socket string `json:"socket"`

	// Name of the database.
	Dbname string `json:"dbname"`

	// Database charset.
	Charset string `json:"charset"`
}

var configDrivers = map[string]bool{
	"sqlite3": true,
	"mysql":   true,
}

func NewConfig() Config {
	return Config{}
}

func (c Config) Instance() (*Conn, error) {

	var err error

	if !configDrivers[c.Driver] {
		return nil, errors.New("Driver can not found")
	}

	var cn Conn

	cn.db, err = sql.Open(c.Driver, c.DbPath)
	if err != nil {
		return nil, err
	}

	cn.cfg = c

	return &cn, nil
}

/*
func (c ConfigSource) NewConnection() (*mysql.DataConn, error) {

    var err error

    if !configDrivers[c.Driver] {
        return nil, errors.New("Driver can not found")
    }

    dsn := ""
    if c.Host != "" {
        dsn = fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=%s`,
            c.User, c.Pass, c.Host, c.Port, c.Dbname, c.Charset)
    } else if c.Socket != `` {
        dsn = fmt.Sprintf(`%s:%s@unix(%s)/%s?charset=%s`,
            c.User, c.Pass, c.Socket, c.Dbname, c.Charset)
    }

    dc, err := mysql.Open(c.Driver, dsn)
    if err != nil {
        return nil, err
    }

    return &dc, nil
}
*/
