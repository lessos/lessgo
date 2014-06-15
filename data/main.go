package main

import (
    "../utils"
    "./rdo"
    "./rdo/base"
    "fmt"
)

func main() {

    cfg := base.Config{
        Driver:  "mysql",
        Host:    "127.0.0.1",
        Port:    "3306",
        User:    "root",
        Pass:    "123456",
        Dbname:  "test",
        Engine:  "InnoDB",
        Charset: "utf8",
    }

    dc, err := rdo.NewClient("def", cfg)
    if err != nil {
        fmt.Println(err)
        return
    }

    ds, err := base.LoadDataSetFromFile("./test.db.json")
    fmt.Println("db.json\n", err, ds)

    jds, _ := utils.JsonEncode(ds)
    fmt.Println("\nDataSetJSON\n", jds)

    /*
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
    */

    // cols, err := dc.Dialect.SchemaColumns(cfg.Dbname, "less_dataset_version")
    // jc, _ := utils.JsonEncode(cols)
    // fmt.Println("\nSchemaColumns", jc)

    tables, err := dc.Dialect.SchemaDataSet(cfg.Dbname)
    jt, _ := utils.JsonEncode(tables)
    fmt.Println("\nSchemaTables\n", jt)

    //return
    err = dc.Dialect.SchemaSync(cfg.Dbname, ds)
    fmt.Println("\nSchemaSync\n", err)

    rs, _ := dc.Base.QueryRaw("select * from ids_login")
    fmt.Println("rs len", len(rs))
    for _, v := range rs {
        fmt.Println("v", v.Field("name").String())
    }

    // indexes, err := dc.Dialect.SchemaIndexes(cfg.Dbname, "less_dataset_version")
    // ji, _ := utils.JsonEncode(indexes)
    // fmt.Println("\nSchemaIndexes", ji)

    /*
       newtable := base.NewTable("feed", "", "")

       colid := base.NewColumn("id", "uint64", 8, 0, "", "")
       colid.IndexType = base.IndexTypePrimaryKeyIncr
       //colid.IsPrimaryKey = true
       //colid.IsAutoIncrement = true

       colcreated := base.NewColumn("created", "uint32", 11, 0, "", "")
       colcreated.IndexType = base.IndexTypeIndex

       colname := base.NewColumn("name", "string", 200, 0, "", "")

       newtable.AddColumn(colid)
       newtable.AddColumn(colcreated)
       newtable.AddColumn(colname)

       newtable.AddColumn(base.NewColumn("rel0", "string", 10, 0, "", ""))
       newtable.AddColumn(base.NewColumn("rel1", "string", 10, 0, "", ""))

       idx := base.NewIndex("rel", base.IndexTypeUnique)
       idx.AddColumn("rel0", "rel1")

       newtable.AddIndex(idx)
       newtable.Comment = "TEST COMMENT !!!"

       jtab, _ := utils.JsonEncode(newtable)
       fmt.Println("\nNewTable", jtab)

       tblsql, err := dc.Dialect.SchemaTableCreateSql(newtable)
       fmt.Println("\nSchemaTableCreateSql", err, tblsql)
    */
}
