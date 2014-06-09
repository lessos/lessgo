package base

type DialectInterface interface {
    Base() *Base
    Init(base *Base) error
    QuoteStr(str string) string
    SchemaTableCreateSql(table *Table) (string, error)
    SchemaTableExist(dbName, tableName string) bool
    SchemaDataSet(dbName string) (DataSet, error)
    SchemaTables(dbName string) ([]*Table, error)
    SchemaColumnTypeSql(col *Column) string
    SchemaColumns(dbName, tableName string) ([]*Column, error)
    SchemaIndexes(dbName, tableName string) ([]*Index, error)
}
