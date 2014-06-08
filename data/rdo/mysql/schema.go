package mysql

import (
    "../base"
    "errors"
    "fmt"
    "strconv"
    "strings"
)

func (dc *mysqlDialect) SchemaColumnTypeSql(col *base.Column) string {

    sql, ok := columnTypes[col.Type]
    if !ok {
        return dc.QuoteStr(col.Name) + col.Type
        //, errors.New("Unsupported column type `" + col.Type + "`")
    }

    switch col.Type {
    case "string":
        sql = fmt.Sprintf(sql, col.Length)
    case "float64-decimal":
        sql = fmt.Sprintf(sql, col.Length, col.Length2)
    }

    if col.Nullable == "null" {
        sql += " NULL"
    } else if col.Nullable == "not-null" {
        sql += " NOT NULL"
    }

    if col.IndexType == base.IndexTypePrimaryKeyIncr {
        sql += " AUTO_INCREMENT"
    }

    if col.Default != "" {
        if col.Default == "null" {
            sql += " DEFAULT NULL"
        } else {
            sql += " DEFAULT " + dc.QuoteStr(col.Default)
        }
    }

    return dc.QuoteStr(col.Name) + " " + sql
}

func (dc *mysqlDialect) SchemaTableCreateSql(table *base.Table) (string, error) {

    if len(table.PrimaryKeys) == 0 {
        return "", errors.New("No PRIMARY KEY")
    }

    if len(table.Columns) == 0 {
        return "", errors.New("No Columns")
    }

    sql := "CREATE TABLE IF NOT EXISTS " + dc.QuoteStr(table.Name) + " (\n"

    for _, col := range table.Columns {
        sql += " " + dc.SchemaColumnTypeSql(col) + ",\n"
    }

    for _, idx := range table.Indexes {
        if len(idx.Cols) == 0 {
            continue
        }
        switch idx.Type {
        case base.IndexTypeIndex:
            sql += " KEY " + dc.QuoteStr(idx.Name) + " ("
            sql += dc.QuoteStr(strings.Join(idx.Cols, dc.QuoteStr(",")))
            sql += "),\n"
        case base.IndexTypeUnique:
            sql += " UNIQUE KEY " + dc.QuoteStr(idx.Name) + " ("
            sql += dc.QuoteStr(strings.Join(idx.Cols, dc.QuoteStr(",")))
            sql += "),\n"
        }
    }

    sql += " PRIMARY KEY ("
    sql += dc.QuoteStr(strings.Join(table.PrimaryKeys, dc.QuoteStr(",")))
    sql += ")\n"

    sql += ")"

    if table.Engine != "" {
        sql += " ENGINE=" + table.Engine
    } else if dc.Base().Config.Engine != "" {
        sql += " ENGINE=" + dc.Base().Config.Engine
    }

    if table.Charset != "" {
        sql += " DEFAULT CHARSET=" + table.Charset
    } else if dc.Base().Config.Charset != "" {
        sql += " DEFAULT CHARSET=" + dc.Base().Config.Charset
    }

    if table.Comment != "" {
        sql += " COMMENT='" + table.Comment + "'"
    }

    sql += ";"

    return sql, nil
}

func (dc *mysqlDialect) SchemaTableExist(dbName, tableName string) bool {

    q := "SELECT `TABLE_NAME` from `INFORMATION_SCHEMA`.`TABLES` "
    q += "WHERE `TABLE_SCHEMA` = ? and `TABLE_NAME` = ?"

    rows, err := dc.Base().QueryRaw(q, dbName, tableName)
    if err != nil {
        return false
    }

    return len(rows) > 0
}

func (dc *mysqlDialect) SchemaTables(dbName string) (map[string]*base.Table, error) {

    tables := map[string]*base.Table{}

    q := "SELECT `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, `AUTO_INCREMENT` "
    q += "FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA` = ?"

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

        pks := []string{}
        cols, _ := dc.SchemaColumns(dbName, name)
        for _, v := range cols {
            if v.IsPrimaryKey() {
                pks = append(pks, v.Name)
            }
        }
        tables[name] = &base.Table{
            Name:        name,
            Engine:      engine,
            AutoIncr:    autoIncr,
            TableRows:   tableRows,
            PrimaryKeys: pks,
            Columns:     cols,
        }
    }

    return tables, nil
}

func (dc *mysqlDialect) SchemaColumns(dbName, tableName string) (map[string]*base.Column, error) {

    cols := map[string]*base.Column{}

    q := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`," +
        " `COLUMN_KEY`, `EXTRA` FROM `INFORMATION_SCHEMA`.`COLUMNS` " +
        " WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    res, err := dc.Base().QueryRaw(q, dbName, tableName)
    if err != nil {
        return cols, err
    }

    for _, record := range res {

        col := &base.Column{}

        for name, v := range record {
            content := v.(string)
            switch name {
            case "COLUMN_NAME":
                col.Name = strings.Trim(content, "` ")
            case "IS_NULLABLE":
                if "YES" == content {
                    col.Nullable = "null"
                }
            case "COLUMN_DEFAULT":
                // add ''
                col.Default = content
                if col.Default == "" {
                    //col.DefaultIsEmpty = true
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
                    //col.IsPrimaryKey = true
                    if col.IndexType != base.IndexTypePrimaryKeyIncr {
                        col.IndexType = base.IndexTypePrimaryKey
                    }
                }
                if key == "UNI" {
                    col.IndexType = base.IndexTypeUnique
                }
            case "EXTRA":
                extra := content
                if extra == "auto_increment" {
                    //col.IsAutoIncrement = true
                    col.IndexType = base.IndexTypePrimaryKeyIncr
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

func (dc *mysqlDialect) SchemaIndexes(dbName, tableName string) (map[string]*base.Index, error) {

    s := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    rows, err := dc.Base().Conn.Query(s, dbName, tableName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    indexes := map[string]*base.Index{}

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
            indexType = base.IndexTypeIndex
        } else {
            indexType = base.IndexTypeUnique
        }

        //fmt.Println("AA", indexType, indexName, colName, nonUnique)
        colName = strings.Trim(colName, "` ")
        if strings.HasPrefix(indexName, "IDX_"+tableName) ||
            strings.HasPrefix(indexName, "UQE_"+tableName) {
            indexName = indexName[5+len(tableName) : len(indexName)]
        }

        var index *base.Index
        var ok bool
        if index, ok = indexes[indexName]; !ok {
            index = &base.Index{
                Name: indexName,
                Type: indexType,
            }
            indexes[indexName] = index
        }

        index.AddColumn(colName)
    }

    return indexes, nil
}
