package schema

type Table struct {
    Name      string
    Engine    string
    Charset   string
    Columns   map[string]*Column
    AutoIncr  string
    TableRows string
}
