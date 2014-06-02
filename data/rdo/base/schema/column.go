package schema

import (
    //"fmt"
    //"reflect"
    //"strings"
    "errors"
)

var fieldTypes = map[string]string{
    "bool":            "bool",
    "string":          "varchar",
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
    "float64-decimal": "numeric",
}

// database column
type Column struct {
    Name string
    Type string
    //FieldName string
    //SQLType         SQLType
    Length   int
    Length2  int
    Nullable bool
    Default  string
    //Indexes         map[string]bool
    IsPrimaryKey    bool
    IsAutoIncrement bool
    DefaultIsEmpty  bool
    //MapType         int
    //IsCreated       bool
    //IsUpdated       bool
    //IsCascade       bool
    //IsVersion       bool
    //fieldPath       []string
    //EnumOptions     map[string]int
}

func (c *Column) TypeDialect() (string, error) {

    str, ok := fieldTypes[c.Type]
    if !ok {
        return "", errors.New("Unsupported column type `" + c.Type + "`")
    }

    return str, nil
}
