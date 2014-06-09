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

    if col.IsAutoIncrement && col.IndexType == base.IndexTypePrimaryKey {
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

func (dc *mysqlDialect) SchemaTables(dbName string) ([]*base.Table, error) {

    tables := []*base.Table{}

    q := "SELECT `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, `AUTO_INCREMENT`, "
    q += "`TABLE_COLLATION`, `TABLE_COMMENT` "
    q += "FROM `INFORMATION_SCHEMA`.`TABLES` "
    q += "WHERE `TABLE_SCHEMA` = ?"

    rows, err := dc.Base().Conn.Query(q, dbName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {

        var name, engine, tableRows, autoIncr, charset, comment string
        if err = rows.Scan(&name, &engine, &tableRows, &autoIncr, &charset, &comment); err != nil {
            return nil, err
        }

        if charset == "utf8_general_ci" {
            charset = "utf8"
        }

        idxs, _ := dc.SchemaIndexes(dbName, name)

        pks := []string{}
        cols, _ := dc.SchemaColumns(dbName, name)
        for i, col := range cols {

            if col.IsPrimaryKey() {
                pks = append(pks, col.Name)
            }

            for _, idx := range idxs {

                for _, idxcol := range idx.Cols {

                    if col.Name == idxcol {

                        if len(idx.Cols) > 1 {

                            cols[i].IndexType = base.IndexTypeMultiple

                        } else {

                            cols[i].IndexType = idx.Type
                        }
                    }
                }
            }
        }

        tables = append(tables, &base.Table{
            Name:        name,
            Engine:      engine,
            Charset:     charset,
            PrimaryKeys: pks,
            Columns:     cols,
            Indexes:     idxs,
            Comment:     comment,
            //AutoIncr:    autoIncr,
            //TableRows:   tableRows,

        })
    }

    return tables, nil
}

func (dc *mysqlDialect) SchemaColumns(dbName, tableName string) ([]*base.Column, error) {

    cols := []*base.Column{}

    q := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`, "
    q += "`COLUMN_KEY`, `EXTRA` "
    q += "FROM `INFORMATION_SCHEMA`.`COLUMNS` "
    q += "WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    rs, err := dc.Base().QueryRaw(q, dbName, tableName)
    if err != nil {
        return cols, err
    }

    for _, record := range rs {

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
                // if col.Default == "" {
                //     col.DefaultIsEmpty = true
                // }
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
                // if _, ok := sqlTypes[colType]; ok {
                //  col.SQLType = SQLType{colType, len1, len2}
                // } else {
                //  return nil, nil, errors.New(fmt.Sprintf("unkonw colType %v", colType))
                // }
            case "COLUMN_KEY":
                switch content {
                case "PRI":
                    col.IndexType = base.IndexTypePrimaryKey
                case "UNI":
                    col.IndexType = base.IndexTypeUnique
                case "MUL":
                    col.IndexType = base.IndexTypeMultiple
                }
            case "EXTRA":
                if content == "auto_increment" {
                    col.IsAutoIncrement = true
                }
            }

            // if col.SQLType.IsText() {
            //     if col.Default != "" {
            //         col.Default = "'" + col.Default + "'"
            //     } else {
            //         if col.DefaultIsEmpty {
            //             col.Default = "''"
            //         }
            //     }
            // }
        }

        cols = append(cols, col)
    }

    return cols, nil
}

func (dc *mysqlDialect) SchemaIndexes(dbName, tableName string) ([]*base.Index, error) {

    indexes := []*base.Index{}

    s := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` "
    s += "FROM `INFORMATION_SCHEMA`.`STATISTICS` "
    s += "WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

    rows, err := dc.Base().Conn.Query(s, dbName, tableName)
    if err != nil {
        return indexes, err
    }
    defer rows.Close()

    for rows.Next() {

        var indexType int
        var indexName, colName, nonUnique string

        if err = rows.Scan(&indexName, &nonUnique, &colName); err != nil {
            return indexes, err
        }

        if indexName == "PRIMARY" {
            indexType = base.IndexTypePrimaryKey
        } else if "YES" == nonUnique || nonUnique == "1" {
            indexType = base.IndexTypeIndex
        } else {
            indexType = base.IndexTypeUnique
        }

        // colName = strings.Trim(colName, "` ")
        // if strings.HasPrefix(indexName, "IDX_"+tableName) ||
        //     strings.HasPrefix(indexName, "UQE_"+tableName) {
        //     indexName = indexName[5+len(tableName) : len(indexName)]
        // }

        exist := false
        for i, v := range indexes {

            if v.Name == indexName {
                indexes[i].AddColumn(colName)
                exist = true
            }
        }

        if !exist {
            //idx := base.NewIndex(indexName, indexType)
            //idx.AddColumn(colName)
            indexes = append(indexes, base.NewIndex(indexName, indexType).AddColumn(colName))
        }

        // indexes[indexName].AddColumn(colName)
    }

    return indexes, nil
}
