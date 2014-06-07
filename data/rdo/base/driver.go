package base

type DialectInterface interface {
    Base() *Base
    Init(base *Base) error
    QuoteStr(str string) string
    SchemaTableCreateSql(table *Table) (string, error)
    SchemaTableExist(dbName, tableName string) bool
    SchemaTables(dbName string) (map[string]*Table, error)
    SchemaColumnTypeSql(col *Column) string
    SchemaColumns(dbName, tableName string) (map[string]*Column, error)
    SchemaIndexes(dbName, tableName string) (map[string]*Index, error)
}
