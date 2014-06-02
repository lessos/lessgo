package schema

type DataSet struct {
    Version uint `json:"version"`
    Tables  []Table
}
