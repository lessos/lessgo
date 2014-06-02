package base

import (
    "./schema"
)

type DialectInterface interface {
    Base() *Base
    Init(base *Base) error
    SchemaTables(dbName string) (map[string]*schema.Table, error)
    SchemaColumns(dbName, tableName string) (map[string]*schema.Column, error)
    SchemaIndexes(dbName, tableName string) (map[string]*schema.Index, error)
}
