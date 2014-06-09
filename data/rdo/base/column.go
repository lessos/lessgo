package base

import (
//"fmt"
//"reflect"
//"strings"
//"errors"
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

// database column
type Column struct {
    Name            string `json:"name"`
    Type            string `json:"type"`
    Length          int    `json:"length"`
    Length2         int    `json:"length2"`
    Nullable        string `json:"nullable"`
    Default         string `json:"default"`
    IndexType       int    `json:"indexType"`
    IsAutoIncrement bool   `json:"isAutoIncrement"`
    Comment         string `json:"comment"`
}

func NewColumn(colName, colType string, len1, len2 int, null, def string) *Column {
    return &Column{
        Name:            colName,
        Type:            colType,
        Length:          len1,
        Length2:         len2,
        Nullable:        null,
        Default:         def,
        IndexType:       IndexTypeEmpty,
        IsAutoIncrement: false,
    }
}

func (col *Column) IsPrimaryKey() bool {
    return col.IndexType == IndexTypePrimaryKey
}

// generate column description string according dialect
func (col *Column) String(d DialectInterface) string {

    sql := d.QuoteStr(col.Name) + " "

    sql += d.SchemaColumnTypeSql(col) + " "

    /* if col.IsPrimaryKey {
           sql += "PRIMARY KEY "
           if col.IsAutoIncrement {
               sql += "AUTO_INCREMENT "
           }
       }

       if col.Nullable {
           sql += "NULL "
       } else {
           sql += "NOT NULL "
       }

       if col.Default != "" {
           sql += "DEFAULT " + col.Default + " "
       } */

    return sql
}

func (col *Column) StringNoPk(d DialectInterface) string {

    sql := d.QuoteStr(col.Name) + " "

    sql += d.SchemaColumnTypeSql(col) + " "

    /*
       if col.Nullable {
           sql += "NULL "
       } else {
           sql += "NOT NULL "
       }

       if col.Default != "" {
           sql += "DEFAULT " + col.Default + " "
       }
    */

    return sql
}
