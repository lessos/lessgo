package base

type DialectInterface interface {
    Base() *Base
    Init(base *Base) error
    QuoteStr(str string) string
    SchemaIndexAdd(dbName, tableName string, index *Index) error
    SchemaIndexDel(dbName, tableName, indexName string) error
    SchemaColumnAddSql(dbName, tableName string, col *Column) (string, error)
    SchemaColumnSetSql(dbName, tableName string, col *Column) (string, error)
    SchemaTableCreateSql(table *Table) (string, error)
    SchemaTableExist(dbName, tableName string) bool
    SchemaSync(dbName string, ds DataSet) error
    SchemaDataSet(dbName string) (DataSet, error)
    SchemaTables(dbName string) ([]*Table, error)
    SchemaColumnTypeSql(col *Column) string
    SchemaColumns(dbName, tableName string) ([]*Column, error)
    SchemaIndexes(dbName, tableName string) ([]*Index, error)
}
