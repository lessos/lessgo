package rdc

var sqliteFieldTypes = map[string]string{
    "auto":            "integer NOT NULL PRIMARY KEY AUTOINCREMENT",
    "pk":              "NOT NULL PRIMARY KEY",
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
