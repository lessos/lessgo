package base

const (
    IndexTypeEmpty          int = 0
    IndexTypeIndex          int = 1
    IndexTypeUnique         int = 2
    IndexTypePrimaryKey     int = 3
    IndexTypePrimaryKeyIncr int = 4
)

// database index
type Index struct {
    Name string
    Type int
    Cols []string
}

// add columns which will be composite index
func (index *Index) AddColumn(cols ...string) {
    for _, col := range cols {
        index.Cols = append(index.Cols, col)
    }
}

// new an index
func NewIndex(name string, indexType int) *Index {
    return &Index{name, indexType, []string{}}
}
