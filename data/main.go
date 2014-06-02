package main

import (
    "../utils"
    "./rdo"
    "./rdo/base"
    "fmt"
)

func main() {

    cfg := base.Config{
        Driver: "mysql",
        Host:   "127.0.0.1",
        Port:   "3306",
        User:   "root",
        Pass:   "123456",
        Dbname: "lessids",
    }

    dc, err := rdo.NewClient("def", cfg)
    if err != nil {
        fmt.Println(err)
        return
    }

    set := map[string]interface{}{
        "id":      100,
        "version": 2,
        "action":  "get",
        "created": "2014-01-01",
    }
    _, err = dc.Base.InsertIgnore("less_dataset_version", set)
    if err == nil {
        fmt.Println("InsertIgnore OK")
    }

    cols, err := dc.Dialect.SchemaColumns(cfg.Dbname, "less_dataset_version")
    jc, _ := utils.JsonEncode(cols)
    fmt.Println("SchemaColumns", jc)

    tables, err := dc.Dialect.SchemaTables(cfg.Dbname)
    jt, _ := utils.JsonEncode(tables)
    fmt.Println("SchemaTables", jt)

    indexes, err := dc.Dialect.SchemaIndexes(cfg.Dbname, "less_dataset_version")
    ji, _ := utils.JsonEncode(indexes)
    fmt.Println("SchemaIndexes", ji)
}
