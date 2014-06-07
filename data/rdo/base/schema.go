package base

type DataSet struct {
    Version uint `json:"version"`
    Tables  []Table
}
