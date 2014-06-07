package mysql

import (
    "../base"
)

const (
    mysqlQuote = "`"
)

var columnTypes = map[string]string{
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

var mysqlStmt = map[string]string{
    "insertIgnore": "INSERT IGNORE INTO `%s` (`%s`) VALUES (%s)",
}

func (dc *mysqlDialect) QuoteStr(str string) string {
    return mysqlQuote + str + mysqlQuote
}

func SqlType(col base.Column) string {

    return ""
}
