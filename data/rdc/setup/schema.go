package setup

import (
)

const (
    FieldIndexIndex   int = 1
    FieldIndexPrimary int = 2
    FieldIndexUnique  int = 3

    FieldTypeInt       string = "int"       // Integer
    FieldTypeUint      string = "uint"      // Integer - Unsigned
    FieldTypeUinti     string = "uinti"     // Integer - Auto Increment
    FieldTypeVarchar   string = "varchar"   // Varchar
    FieldTypeString    string = "text"      // Text
    FieldTypeTimestamp string = "timestamp" // Integer
)

type DataSet struct {
    Tables []Table
}

type Table struct {
    Name   string  `json:"name"`
    Fields []Field `json:"fields"`
}

type Field struct {
    Name string `json:"name"`
    Type string `json:"type"`
    Len  int    `json:"len"`
    Idx  int    `json:"idx"`
}

func NewTable(name string) Table {
    return Table{Name: name}
}

func (t *Table) FieldAdd(nm, tp string, ln, ix int) {
    t.Fields = append(t.Fields, Field{nm, tp, ln, ix})
}

func NewDataSet() DataSet {
    return DataSet{}
}

func (s *DataSet) TableAdd(t Table) {
    s.Tables = append(s.Tables, t)
}


