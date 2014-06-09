package base

type Table struct {
    Name        string    `json:"name"`
    Engine      string    `json:"engine"`
    Charset     string    `json:"charset"`
    PrimaryKeys []string  `json:"primary_keys"`
    Columns     []*Column `json:"columns"`
    Indexes     []*Index  `json:"indexes"`
    Comment     string    `json:"comment"`
    //AutoIncr    string
    //TableRows   string
}

func NewTable(name, engine, charset string) *Table {
    return &Table{
        Name:        name,
        Engine:      engine,
        Charset:     charset,
        PrimaryKeys: []string{},
        Columns:     []*Column{},
        Indexes:     []*Index{},
    }
}

func (table *Table) AddColumn(col *Column) {

    //table.Columns[col.Name] = col
    exist := false

    for k, v := range table.Columns {
        if v.Name == col.Name {
            table.Columns[k] = col
            exist = true
        }
    }

    if !exist {
        table.Columns = append(table.Columns, col)
    }

    if col.IsPrimaryKey() {
        table.PrimaryKeys = append(table.PrimaryKeys, col.Name)
    }
}

func (table *Table) AddIndex(index *Index) {
    for _, v := range table.Indexes {
        if v.Name == index.Name {
            return
        }
    }

    table.Indexes = append(table.Indexes, index)
}
