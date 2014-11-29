package main

import (
	"fmt"
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
)

func main() {

	cfg := base.Config{
		Driver:  "sqlite3",
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "123456",
		Dbname:  "test",
		Engine:  "InnoDB",
		Socket:  "./sqlite3.db",
		Charset: "utf8",
	}

	dc, err := rdo.NewClient("def", cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	ds, err := base.LoadDataSetFromFile("./test.db.json")
	// fmt.Println("db.json\n", err, ds)

	// jds, _ := utils.JsonEncode(ds)
	// fmt.Println("\nDataSetJSON\n", jds)

	// cols, err := dc.Dialect.SchemaColumns(cfg.Dbname, "less_dataset_version")
	// jc, _ := utils.JsonEncode(cols)
	// fmt.Println("\nSchemaColumns", jc)

	// fmt.Println("CHECK Dialect.SchemaDataSet")
	// tables, err := dc.Dialect.SchemaDataSet(cfg.Dbname)
	// jt, _ := utils.JsonEncode(tables)
	// fmt.Println(jt)

	fmt.Println("CHECK Dialect.SchemaSync")
	err = dc.Dialect.SchemaSync(cfg.Dbname, ds)
	if err == nil {
		fmt.Println("\tOK")
	} else {
		fmt.Println("\tERROR", err)
	}

	fmt.Println("CHECK Base.Insert")
	set := map[string]interface{}{
		"id":      100,
		"version": 2,
		"action":  "get",
		"created": "2014-01-01",
	}
	_, err = dc.Base.Insert("less_dataset_version", set)
	if err == nil {
		fmt.Println("\tOK")
	} else {
		fmt.Println("\tERROR", err)
	}

	fmt.Println("CHECK Base.InsertIgnore")
	_, err = dc.Base.InsertIgnore("less_dataset_version", set)
	if err == nil {
		fmt.Println("\tOK")
	} else {
		fmt.Println("\tERROR", err)
	}

	fmt.Println("CHECK SchemaTableExist")
	if dc.Dialect.SchemaTableExist(cfg.Dbname, "less_dataset_version") {
		fmt.Println("\tless_dataset_version YES")
	}

	if !dc.Dialect.SchemaTableExist(cfg.Dbname, "less_dataset_version_NULL") {
		fmt.Println("\tless_dataset_version_NULL NO")
	}

	fmt.Println("CHECK SchemaColumnAdd")
	col := base.NewColumn("new_col", "uint32", "", true, "")
	if err := dc.Dialect.SchemaColumnAdd(cfg.Dbname, "less_dataset_version", col); err != nil {
		fmt.Println("\tERROR", err)
	} else {
		fmt.Println("\tOK")
	}

	fmt.Println("CHECK SchemaColumnDel")
	if err := dc.Dialect.SchemaColumnDel(cfg.Dbname, "less_dataset_version", col); err != nil {
		fmt.Println("\tERROR", err)
	} else {
		fmt.Println("\tOK")
	}

	// rs, _ := dc.Base.QueryRaw("select * from ids_login")
	// fmt.Println("rs len", len(rs))
	// for _, v := range rs {
	//     fmt.Println("v", v.Field("name").String())
	// }

	// indexes, err := dc.Dialect.SchemaIndexes(cfg.Dbname, "less_dataset_version")
	// ji, _ := utils.JsonEncode(indexes)
	// fmt.Println("\nSchemaIndexes", ji)
}
