package sqlite3

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var FieldTypes = map[string]string{
    "auto":            "integer NOT NULL PRIMARY KEY AUTOINCREMENT",
    "pk":              "integer NOT NULL PRIMARY KEY",
    "pk-string":       "varchar(%d) NOT NULL PRIMARY KEY",
    "bool":            "bool",
    "string":          "varchar(%d)",
    "string-text":     "text",
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
    "float64":         "real",
    "float64-decimal": "decimal",
}

func Open(dirver, dsn string) (*sql.DB, error) {
    return sql.Open(dirver, dsn)
}
