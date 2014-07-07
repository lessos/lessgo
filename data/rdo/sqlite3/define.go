package sqlite3

const (
	sqlite3Quote = "`"
)

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

var sqlite3Stmt = map[string]string{
	"insertIgnore": "INSERT OR IGNORE INTO `%s` (`%s`) VALUES (%s)",
}

func (dc *sqlite3Dialect) QuoteStr(str string) string {
	return sqlite3Quote + str + sqlite3Quote
}
