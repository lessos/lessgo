package base

import ()

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
    Name     string `json:"name"`
    Type     string `json:"type"`
    Length   int    `json:"length"`
    Length2  int    `json:"length2"`
    NullAble bool   `json:"nullAble"`
    IncrAble bool   `json:"incrAble"`
    Default  string `json:"default"`
    Comment  string `json:"comment"`
}

func NewColumn(colName, colType string, len1, len2 int, null bool, def string) *Column {
    return &Column{
        Name:     colName,
        Type:     colType,
        Length:   len1,
        Length2:  len2,
        NullAble: null,
        IncrAble: false,
        Default:  def,
        Comment:  "",
    }
}
