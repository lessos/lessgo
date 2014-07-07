package setup

import ()

const (

	//FieldIndexPrimary     int = 1
	//FieldIndexAutoPrimary int = 2
	FieldIndexIndex  int = 3
	FieldIndexUnique int = 4

	FieldTypeInt       string = "int"       // Integer
	FieldTypeUint      string = "uint"      // Integer - Unsigned
	FieldTypeChar      string = "char"      // Char
	FieldTypeVarchar   string = "varchar"   // Varchar
	FieldTypeString    string = "text"      // Text
	FieldTypeTimestamp string = "timestamp" // Integer
)

type DataSet struct {
	Version uint `json:"version"`
	Tables  []Table
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
