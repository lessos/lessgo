package mysql

import (
    "../base"
    "errors"
    "fmt"
    "sort"
    "strconv"
    "strings"
)

// "SET NAMES 'utf8'"
// "SET CHARACTER_SET_CLIENT=utf8"
// "SET CHARACTER_SET_RESULTS=utf8"

func (dc *mysqlDialect) SchemaIndexAdd(dbName, tableName string, index *base.Index) error {

    sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` ADD INDEX (%s)", 
        dbName, tableName, dc.QuoteStr(strings.Join(index.Cols, dc.QuoteStr(","))))
    
    if _, err := dc.Base().ExecRaw(sql); err != nil {
        return err
    }

    return nil
}

func (dc *mysqlDialect) SchemaIndexDel(dbName, tableName, indexName string) error {

    if indexName == "PRIMARY" {
        fmt.Println("DROP PRIMARY KEY")
        return nil // TODO
    }

    sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` DROP ", dbName, tableName)
    if indexName == "PRIMARY" {
        sql += "PRIMARY KEY"
    } else {
        sql += "INDEX "+ indexName
    }

    fmt.Println("SchemaIndexDel", sql)

    if _, err := dc.Base().ExecRaw(sql); err != nil {
        return err
    }

    return nil
}

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

    if col.NullAble {
        sql += " NULL"
    } else {
        sql += " NOT NULL"
    }

    if col.IncrAble && col.IndexType == base.IndexTypePrimaryKey {
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

func (dc *mysqlDialect) SchemaColumnAddSql(dbName, tableName string, col *base.Column) (string, error) {

    sql := fmt.Sprintf("ALTER TABLE `%v`.`%v` ADD %v", dbName, tableName, dc.SchemaColumnTypeSql(col))

    return sql, nil
}

func (dc *mysqlDialect) SchemaColumnSetSql(dbName, tableName string, col *base.Column) (string, error) {

    sql := fmt.Sprintf("ALTER TABLE `%v`.`%v` CHANGE `%v` %v",
        dbName, tableName, col.Name, dc.SchemaColumnTypeSql(col))

    return sql, nil
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

func (dc *mysqlDialect) SchemaSync(dbName string, newds base.DataSet) error {

    curds, err := dc.SchemaDataSet(dbName)
    if err != nil {
        return err
    }

    for _, newTable := range newds.Tables {

        exist := false
        //updated := false
        var curTable *base.Table

        for _, curTable = range curds.Tables {

            if newTable.Name == curTable.Name {
                exist = true
                break
            }
        }

        if !exist {

            sql, err := dc.SchemaTableCreateSql(newTable)

            _, err = dc.Base().ExecRaw(sql)
            if err != nil {
                return err
            }

            continue
        }

        // Column
        for _, newcol := range newTable.Columns {

            colExist := false
            colChange := false

            for _, curcol := range curTable.Columns {

                if newcol.Name != curcol.Name {
                    continue
                }

                colExist = true

                if newcol.Type != curcol.Type ||
                    newcol.Length != curcol.Length ||
                    newcol.Length2 != curcol.Length2 ||
                    newcol.NullAble != curcol.NullAble ||
                    newcol.IncrAble != curcol.IncrAble ||
                    newcol.Default != curcol.Default {
                    colChange = true
                    break
                }
            }

            if !colExist {
                sql, err := dc.SchemaColumnAddSql(dbName, newTable.Name, newcol)
                _, err = dc.Base().ExecRaw(sql)
                if err != nil {
                    return err
                }
            }

            if colChange {

                sql, err := dc.SchemaColumnSetSql(dbName, newTable.Name, newcol)
                //fmt.Println("colChange", sql)
                _, err = dc.Base().ExecRaw(sql)
                if err != nil {
                    return err
                }
            }
        }

        // Index Del
        for _, curidx := range curTable.Indexes {

            curDel := true

            for _, newidx := range newTable.Indexes {

                if newidx.Name != curidx.Name {
                    continue
                }

                sort.Strings(newidx.Cols)
                sort.Strings(curidx.Cols)

                if newidx.Type == curidx.Type &&
                    strings.Join(newidx.Cols, ",") == strings.Join(curidx.Cols, ",") {

                    curDel = false
                }

                break
            }

            if curDel {
                if err := dc.SchemaIndexDel(dbName, curTable.Name, curidx.Name); err != nil {
                    return err
                }
            }
        }

        // Index New
        for _, newidx := range newTable.Indexes {

            exist := false

            for _, curidx := range curTable.Indexes {

                if newidx.Name == curidx.Name {
                    exist = true
                    break
                }
            }

            if !exist {
                if err := dc.SchemaIndexAdd(dbName, curTable.Name, newidx); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func (dc *mysqlDialect) SchemaDataSet(dbName string) (base.DataSet, error) {

    ds := base.DataSet{
        DbName: dbName,
    }

    q := "SELECT `DEFAULT_CHARACTER_SET_NAME` "
    q += "FROM `INFORMATION_SCHEMA`.`SCHEMATA` "
    q += "WHERE `SCHEMA_NAME` = ?"

    err := dc.Base().Conn.QueryRow(q, dbName).Scan(&ds.Charset)
    if err != nil {
        return ds, err
    }

    ds.Tables, err = dc.SchemaTables(dbName)

    return ds, err
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

            if col.IndexType == base.IndexTypePrimaryKey {
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
                    col.NullAble = true
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

                typepre := ""
                if strings.Contains(content, "unsigned") {
                    typepre = "u"
                }
                switch col.Type {
                case "bigint":
                    col.Type = typepre + "int64"
                    col.Length = 0
                case "int":
                    col.Type = typepre + "int32"
                    col.Length = 0
                case "smallint":
                    col.Type = typepre + "int16"
                    col.Length = 0
                case "tinyint":
                    col.Type = typepre + "int8"
                    col.Length = 0
                case "varchar":
                    col.Type = "string"
                case "longtext":
                    col.Type = "string-text"
                }

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
                    col.IncrAble = true
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
