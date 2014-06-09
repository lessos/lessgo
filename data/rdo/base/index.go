package base

const (
    IndexTypeEmpty      int = 0
    IndexTypeIndex      int = 1
    IndexTypeUnique     int = 2
    IndexTypePrimaryKey int = 3
    IndexTypeMultiple   int = 4
)

// database index
type Index struct {
    Name string   `json:"name"`
    Type int      `json:"type"`
    Cols []string `json:"cols"`
}

// add columns which will be composite index
func (index *Index) AddColumn(cols ...string) *Index {

    for _, col := range cols {

        exist := false
        for _, v := range index.Cols {
            if v == col {
                exist = true
            }
        }

        if !exist {
            index.Cols = append(index.Cols, col)
        }
    }

    return index
}

// new an index
func NewIndex(name string, indexType int) *Index {
    return &Index{name, indexType, []string{}}
}
