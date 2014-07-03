package base

var columnTypes = map[string]string{
    "bool":            "bool",
    "string":          "varchar(%v)",
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
    "float64-decimal": "numeric(%v, %v)",
}

// database column
type Column struct {
    Name     string   `json:"name"`
    Type     string   `json:"type"`
    Length   string   `json:"length"`
    NullAble bool     `json:"nullAble"`
    IncrAble bool     `json:"incrAble"`
    Default  string   `json:"default"`
    Comment  string   `json:"comment"`
    Extra    []string `json:"extra"`
}

func NewColumn(colName, colType, len string, null bool, def string) *Column {
    return &Column{
        Name:     colName,
        Type:     colType,
        Length:   len,
        NullAble: null,
        IncrAble: false,
        Default:  def,
        Comment:  "",
        Extra:    []string{},
    }
}
