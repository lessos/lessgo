package base

type DialectInterface interface {
    Base() *Base
    Init(base *Base) error
    QuoteStr(str string) string

    SchemaIndexAdd(dbName, tableName string, index *Index) error
    SchemaIndexDel(dbName, tableName string, index *Index) error
    SchemaIndexSet(dbName, tableName string, index *Index) error
    SchemaIndexQuery(dbName, tableName string) ([]*Index, error)

    SchemaColumnAdd(dbName, tableName string, col *Column) error
    SchemaColumnDel(dbName, tableName string, col *Column) error
    SchemaColumnSet(dbName, tableName string, col *Column) error
    SchemaColumnQuery(dbName, tableName string) ([]*Column, error)
    SchemaColumnTypeSql(col *Column) string

    SchemaTableAdd(table *Table) error
    SchemaTableQuery(dbName string) ([]*Table, error)
    SchemaTableExist(dbName, tableName string) bool

    SchemaSync(dbName string, ds DataSet) error
    SchemaDataSet(dbName string) (DataSet, error)
}
