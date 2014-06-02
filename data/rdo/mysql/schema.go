package mysql

import (
    "../base/schema"
    "strconv"
    "strings"
    //"fmt"
)

func (dc *mysqlDialect) SchemaTables(dbName string) (map[string]*schema.Table, error) {

    tables := map[string]*schema.Table{}

    q := "SELECT `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, `AUTO_INCREMENT` from `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=?"

    rows, err := dc.Base().Conn.Query(q, dbName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {

        var name, engine, tableRows, autoIncr string
        err = rows.Scan(&name, &engine, &tableRows, &autoIncr)
        if err != nil {
            return nil, err
        }

        //cols, _ := dc.SchemaColumns(dbName, name)
        tables[name] = &schema.Table{
            Name:      name,
            Engine:    engine,
            AutoIncr:  autoIncr,
            TableRows: tableRows,
            //Columns:   cols,
        }
    }

    return tables, nil
}

func (dc *mysqlDialect) SchemaColumns(dbName, tableName string) (map[string]*schema.Column, error) {

    cols := map[string]*schema.Column{}

    q := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`," +
        " `COLUMN_KEY`, `EXTRA` FROM `INFORMATION_SCHEMA`.`COLUMNS` " +
        " WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    res, err := dc.Base().QueryRaw(q, dbName, tableName)
    if err != nil {
        return cols, err
    }

    for _, record := range res {

        col := &schema.Column{}

        for name, v := range record {
            content := v.(string)
            switch name {
            case "COLUMN_NAME":
                col.Name = strings.Trim(content, "` ")
            case "IS_NULLABLE":
                if "YES" == content {
                    col.Nullable = true
                }
            case "COLUMN_DEFAULT":
                // add ''
                col.Default = content
                if col.Default == "" {
                    col.DefaultIsEmpty = true
                }
            case "COLUMN_TYPE":
                cts := strings.Split(content, "(")
                var len1, len2 int
                if len(cts) == 2 {
                    idx := strings.Index(cts[1], ")")
                    lens := strings.Split(cts[1][0:idx], ",")
                    len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
                    if err != nil {
                        //return nil, nil, err
                        continue
                    }
                    if len(lens) == 2 {
                        len2, err = strconv.Atoi(lens[1])
                        if err != nil {
                            //return nil, nil, err
                            continue
                        }
                    }
                }

                col.Type = strings.ToLower(cts[0])
                col.Length = len1
                col.Length2 = len2
                //if _, ok := sqlTypes[colType]; ok {
                //  col.SQLType = SQLType{colType, len1, len2}
                //} else {
                //  return nil, nil, errors.New(fmt.Sprintf("unkonw colType %v", colType))
                //}
            case "COLUMN_KEY":
                key := content
                if key == "PRI" {
                    col.IsPrimaryKey = true
                }
                if key == "UNI" {
                    //col.is
                }
            case "EXTRA":
                extra := content
                if extra == "auto_increment" {
                    col.IsAutoIncrement = true
                }
            }

            /*
               if col.SQLType.IsText() {
               if col.Default != "" {
                   col.Default = "'" + col.Default + "'"
               } else {
                   if col.DefaultIsEmpty {
                       col.Default = "''"
                   }
               }
               }
            */
        }

        cols[col.Name] = col
    }

    return cols, nil
}

func (dc *mysqlDialect) SchemaIndexes(dbName, tableName string) (map[string]*schema.Index, error) {

    s := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    rows, err := dc.Base().Conn.Query(s, dbName, tableName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    indexes := map[string]*schema.Index{}

    for rows.Next() {

        var indexType int
        var indexName, colName, nonUnique string
        err = rows.Scan(&indexName, &nonUnique, &colName)
        if err != nil {
            return indexes, err
        }

        if indexName == "PRIMARY" {
            continue
        }

        if "YES" == nonUnique || nonUnique == "1" {
            indexType = schema.IndexType
        } else {
            indexType = schema.UniqueType
        }

        //fmt.Println("AA", indexType, indexName, colName, nonUnique)
        colName = strings.Trim(colName, "` ")
        if strings.HasPrefix(indexName, "IDX_"+tableName) ||
            strings.HasPrefix(indexName, "UQE_"+tableName) {
            indexName = indexName[5+len(tableName) : len(indexName)]
        }

        var index *schema.Index
        var ok bool
        if index, ok = indexes[indexName]; !ok {
            index = &schema.Index{
                Name: indexName,
                Type: indexType,
            }
            indexes[indexName] = index
        }

        index.AddColumn(colName)
    }

    return indexes, nil
}
