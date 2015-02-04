package base

type Table struct {
	Name    string    `json:"name"`
	Engine  string    `json:"engine,omitempty"`
	Charset string    `json:"charset,omitempty"`
	Columns []*Column `json:"columns"`
	Indexes []*Index  `json:"indexes"`
	Comment string    `json:"comment,omitempty"`
}

func NewTable(name, engine, charset string) *Table {
	return &Table{
		Name:    name,
		Engine:  engine,
		Charset: charset,
		Columns: []*Column{},
		Indexes: []*Index{},
	}
}

func (table *Table) AddColumn(col *Column) {

	for k, v := range table.Columns {

		if v.Name != col.Name {
			continue
		}

		table.Columns[k] = col
		return
	}

	table.Columns = append(table.Columns, col)
}

func (table *Table) AddIndex(index *Index) {

	for k, v := range table.Indexes {

		if v.Name != index.Name {
			continue
		}

		table.Indexes[k] = index
		return
	}

	table.Indexes = append(table.Indexes, index)
}
