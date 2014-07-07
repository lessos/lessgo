package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var FieldTypes = map[string]string{
	"auto":            "int(11) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY",
	"pk":              "int(11) unsigned NOT NULL PRIMARY KEY",
	"pk-string":       "varchar(%d) NOT NULL",
	"bool":            "bool",
	"string":          "varchar(%d)",
	"string-text":     "longtext",
	"date":            "date",
	"datetime":        "datetime",
	"int8":            "tinyint",
	"int16":           "smallint",
	"int32":           "integer",
	"int64":           "bigint",
	"uint8":           "tinyint unsigned",
	"uint16":          "smallint unsigned",
	"uint32":          "integer unsigned",
	"uint64":          "bigint unsigned",
	"float64":         "double precision",
	"float64-decimal": "numeric(%d, %d)",
}

func Open(dirver, dsn string) (*sql.DB, error) {
	return sql.Open(dirver, dsn)
}
