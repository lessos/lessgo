package sqlite3

import (
    "../base"
    "errors"
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

// http://stackoverflow.com/questions/11983924/optimize-sql-query-on-sqlite3-by-using-indexes

var (
    // CREATE INDEX {index-name} ON {table-name} (`column-a`, `column-b`)
    regIndexSql = regexp.MustCompile(`INDEX(.*)ON(.*)\((.*)\)`)
    // `id` integer NULL PRIMARY KEY AUTOINCREMENT
    regIncrSql = regexp.MustCompile("\\`(.*)\\`(.*)PRIMARY(.*)AUTOINCREMENT")
)

func (dc *sqlite3Dialect) SchemaIndexAdd(dbName, tableName string, index *base.Index) error {

    action := ""
    switch index.Type {
    case base.IndexTypeIndex:
        action = "INDEX"
    case base.IndexTypeUnique:
        action = "UNIQUE INDEX"
    default:
        // PRIMARY KEY can be modified, can not be added
        // return errors.New("Invalid Index Type")
        return nil
    }

    sql := fmt.Sprintf("CREATE %s IF NOT EXISTS %s_%s_idx ON %s (%s)",
        action, tableName, index.Name, tableName,
        dc.QuoteStr(strings.Join(index.Cols, dc.QuoteStr(","))))
    //fmt.Println(sql)
    _, err := dc.Base().ExecRaw(sql)

    return err
}

func (dc *sqlite3Dialect) SchemaIndexDel(dbName, tableName string, index *base.Index) error {

    // PRIMARY KEY can be modified, can not be deleted
    if index.Type == base.IndexTypePrimaryKey {
        return nil
    }

    sql := fmt.Sprintf("DROP INDEX IF EXISTS `%s_%s_idx`", tableName, index.Name)

    _, err := dc.Base().ExecRaw(sql)

    return err
}

func (dc *sqlite3Dialect) SchemaIndexSet(dbName, tableName string, index *base.Index) error {

    if err := dc.SchemaIndexDel(dbName, tableName, index); err != nil {
        return err
    }

    return dc.SchemaIndexSet(dbName, tableName, index)
}

func (dc *sqlite3Dialect) SchemaIndexQuery(dbName, tableName string) ([]*base.Index, error) {

    indexes := []*base.Index{}

    q := "SELECT sql FROM sqlite_master WHERE type='index' and tbl_name = ?"

    rows, err := dc.Base().Conn.Query(q, tableName)
    if err != nil {
        return indexes, err
    }
    defer rows.Close()

    for rows.Next() {

        var sql string
        if err = rows.Scan(&sql); err != nil {
            return indexes, err
        }
        if sql == "" {
            continue
        }

        mat := regIndexSql.FindStringSubmatch(sql)
        if len(mat) != 4 {
            continue
        }

        //
        indexName := strings.TrimSpace(mat[1])
        if strings.HasPrefix(indexName, tableName+"_") {
            indexName = indexName[len(tableName)+1:]
        }
        if strings.HasSuffix(indexName, "_idx") {
            indexName = indexName[0 : len(indexName)-4]
        }

        //
        indexType := base.IndexTypeIndex
        if strings.HasPrefix(sql, "CREATE UNIQUE INDEX") {
            indexType = base.IndexTypeUnique
        }

        //
        indexCols := strings.Split(strings.Replace(mat[3], "`", "", -1), ",")

        //
        indexes = append(indexes, base.NewIndex(indexName, indexType).AddColumn(indexCols...))
    }

    return indexes, nil
}

func (dc *sqlite3Dialect) SchemaColumnTypeSql(col *base.Column) string {

    sql, ok := columnTypes[col.Type]
    if !ok {
        return dc.QuoteStr(col.Name) + col.Type
        //, errors.New("Unsupported column type `" + col.Type + "`")
    }

    switch col.Type {
    case "string":
        sql = fmt.Sprintf(sql, col.Length)
    case "float64-decimal":
        lens := strings.Split(col.Length, ",")
        if lens[0] == "" {
            lens[0] = "10"
        }
        if len(lens) < 2 {
            lens = append(lens, "2")
        }
        sql = fmt.Sprintf(sql, lens[0], lens[1])
    }

    if base.ArrayContain("notnull", col.Extra) {
        sql += " NOT NULL"
    } else {
        sql += " NULL"
    }

    if base.ArrayContain("PRIMARY", col.Extra) {

        if col.IncrAble {
            sql += " PRIMARY KEY AUTOINCREMENT"
        } else {
            sql += " PRIMARY KEY"
        }
    }

    // INTEGER

    if col.Default != "" {
        if col.Default == "null" {
            sql += " DEFAULT NULL"
        } else {
            sql += " DEFAULT '" + col.Default + "'"
        }
    }

    return dc.QuoteStr(col.Name) + " " + sql
}

func (dc *sqlite3Dialect) SchemaColumnAdd(dbName, tableName string, col *base.Column) error {

    return errors.New("NotImplemented")
}

func (dc *sqlite3Dialect) SchemaColumnDel(dbName, tableName string, col *base.Column) error {

    return errors.New("NotImplemented")
}

func (dc *sqlite3Dialect) SchemaColumnSet(dbName, tableName string, col *base.Column) error {

    return errors.New("NotImplemented")
}

func (dc *sqlite3Dialect) SchemaColumnQuery(dbName, tableName string) ([]*base.Column, error) {

    cols := []*base.Column{}
    incrCol := ""

    //
    q := "SELECT sql FROM sqlite_master WHERE type='table' and name = ?"
    rs, err := dc.Base().QueryRaw(q, tableName)
    if err != nil || len(rs) == 0 {
        return cols, errors.New("No Data Found")
    }
    if mat := regIncrSql.FindStringSubmatch(rs[0].Field("sql").String()); len(mat) == 4 {
        incrCol = mat[1]
    }

    //
    rs, err = dc.Base().QueryRaw("PRAGMA table_info (`" + tableName + "`)")
    if err != nil {
        return cols, err
    }
    for _, entry := range rs {

        col := &base.Column{}

        for name, v := range entry.Fields {

            switch name {
            case "name":
                col.Name = strings.Trim(v.String(), "` ")
                if col.Name == incrCol {
                    col.IncrAble = true
                }
            case "notnull":
                if "1" == v.String() {
                    col.Extra = append(col.Extra, "notnull")
                } else {
                    col.NullAble = true
                }
            case "dflt_value":
                col.Default = v.String()
            case "type":

                cts := strings.Split(v.String(), "(")
                var len1, len2 int

                if len(cts) == 2 {
                    idx := strings.Index(cts[1], ")")
                    lens := strings.Split(cts[1][0:idx], ",")
                    len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
                    if err != nil {
                        continue
                    }
                    if len(lens) == 2 {
                        len2, err = strconv.Atoi(strings.TrimSpace(lens[1]))
                        if err != nil {
                            continue
                        }
                    }
                }

                col.Type = strings.ToLower(cts[0])
                if len1 > 0 {
                    col.Length = fmt.Sprintf("%v", len1)
                    if len2 > 0 {
                        col.Length += fmt.Sprintf(",%v", len2)
                    }
                }

                typepre := ""
                if strings.Contains(v.String(), "unsigned") {
                    ucts := strings.Split(v.String(), " ")
                    if len(ucts) != 2 {
                        continue
                    }
                    col.Type = ucts[0]
                    typepre = "u"
                }
                // TODO fmt.Println("col.Type", col.Name, col.Type)

                switch col.Type {
                case "bigint":
                    col.Type = typepre + "int64"
                case "integer":
                    col.Type = typepre + "int32"
                case "smallint":
                    col.Type = typepre + "int16"
                case "tinyint":
                    col.Type = typepre + "int8"
                case "numeric":
                    col.Type = "float64-decimal" // TODO
                case "double precision":
                    col.Type = "float64"
                case "varchar":
                    col.Type = "string"
                case "longtext":
                    col.Type = "string-text"
                }

            case "pk":
                if v.String() == "1" {
                    col.Extra = append(col.Extra, "PRIMARY")
                }
            }
        }

        cols = append(cols, col)
    }

    return cols, nil
}

func (dc *sqlite3Dialect) SchemaTableAdd(table *base.Table) error {

    if len(table.Columns) == 0 {
        return errors.New("No Columns Found")
    }

    sql := "CREATE TABLE IF NOT EXISTS " + dc.QuoteStr(table.Name) + " ("
    for _, col := range table.Columns {
        sql += "\n " + dc.SchemaColumnTypeSql(col) + ","
    }
    sql = sql[0:len(sql)-1] + "\n);"

    //fmt.Println(sql)
    _, err := dc.Base().ExecRaw(sql)
    if err != nil {
        return err
    }

    for _, idx := range table.Indexes {

        if idx.Type != base.IndexTypeIndex && idx.Type != base.IndexTypeUnique {
            continue
        }

        if len(idx.Cols) == 0 {
            continue
        }

        err = dc.SchemaIndexAdd("", table.Name, idx)
        if err != nil {
            return err
        }
    }

    return err
}

func (dc *sqlite3Dialect) SchemaTableExist(dbName, tableName string) bool {

    q := "SELECT name FROM sqlite_master WHERE type='table' and name = ?"

    rows, err := dc.Base().QueryRaw(q, tableName)
    if err != nil {
        return false
    }

    return len(rows) > 0
}

func (dc *sqlite3Dialect) SchemaSync(dbName string, newds base.DataSet) error {

    curds, err := dc.SchemaDataSet(dbName)
    if err != nil {
        return err
    }

    for i, table := range newds.Tables {

        priCols := []string{}

        for _, index := range table.Indexes {

            if index.Type == base.IndexTypePrimaryKey {

                priCols = index.Cols
                break
            }
        }

        if len(priCols) == 0 {
            continue
        }

        for j, col := range table.Columns {

            for _, priCol := range priCols {

                if priCol == col.Name {
                    newds.Tables[i].Columns[j].Extra = append(newds.Tables[i].Columns[j].Extra, "PRIMARY")

                    // AUTOINCREMENT is only allowed on an INTEGER PRIMARY KEY
                    if col.IncrAble && col.Type != "int32" {
                        col.Type = "int32"
                    }
                }
            }
        }
    }

    for _, newTable := range newds.Tables {

        exist, sql := false, ""
        colsTrans := []string{}

        var curTable *base.Table

        for _, curTable = range curds.Tables {

            if newTable.Name == curTable.Name {
                exist = true
                break
            }
        }

        if !exist {

            if err := dc.SchemaTableAdd(newTable); err != nil {
                return err
            }

            continue
        }

        // Column
        if len(newTable.Columns) != len(curTable.Columns) {
            goto upgrade_columns
        }

        for _, newcol := range newTable.Columns {

            for _, curcol := range curTable.Columns {

                if newcol.Name != curcol.Name {
                    continue
                }

                if newcol.Type != curcol.Type ||
                    newcol.Length != curcol.Length ||
                    newcol.IncrAble != curcol.IncrAble ||
                    newcol.Default != curcol.Default ||
                    !base.ArrayEqual(newcol.Extra, curcol.Extra) {

                    goto upgrade_columns
                }
            }
        }

        goto upgrade_indexes

    upgrade_columns:

        //fmt.Println("\tupgrade_columns")

        sql = "ALTER TABLE `" + newTable.Name + "` RENAME TO `" + newTable.Name + "_temp`"
        _, err = dc.Base().ExecRaw(sql)
        if err != nil {
            return err
        }

        if err := dc.SchemaTableAdd(newTable); err != nil {
            return err
        }

        for _, curcol := range curTable.Columns {

            for _, newcol := range newTable.Columns {

                if curcol.Name == newcol.Name {
                    colsTrans = append(colsTrans, curcol.Name)
                }
            }
        }

        sql = "INSERT INTO `" + newTable.Name + "` (`" + strings.Join(colsTrans, "`,`") + "`) " +
            "SELECT `" + strings.Join(colsTrans, "`,`") + "` FROM `" + newTable.Name + "_temp`"
        //fmt.Println(sql)
        _, err = dc.Base().ExecRaw(sql)

        if err != nil {

            sql = "DROP TABLE `" + newTable.Name + "`"
            dc.Base().ExecRaw(sql)

            sql = "ALTER TABLE `" + newTable.Name + "_temp` RENAME TO `" + newTable.Name + "`"
            dc.Base().ExecRaw(sql)

            return err
        } else {
            sql = "DROP TABLE `" + newTable.Name + "_temp`"
            dc.Base().ExecRaw(sql)
        }

    upgrade_indexes:

        //fmt.Println("upgrade_indexes")

        curIndexes, err := dc.SchemaIndexQuery(dbName, newTable.Name)
        if err != nil {
            return err
        }

        //fmt.Println("\tcuridx count", len(curIndexes))

        // Delete Old Indexes
        for _, curidx := range curIndexes {

            exist := false

            //fmt.Println("\tcuridx", curidx.Name)
            for _, newidx := range newTable.Indexes {

                if curidx.Name != newidx.Name {
                    continue
                }

                if newidx.Type == base.IndexTypePrimaryKey {
                    continue
                }

                exist = true
                break
            }

            if !exist {
                //fmt.Println("\tSchemaIndexDel", newTable.Name, curidx.Name)
                dc.SchemaIndexDel(dbName, newTable.Name, curidx)
            }
        }

        // Add New Indexes
        for _, newidx := range newTable.Indexes {

            exist := false

            for _, curidx := range curIndexes {

                if curidx.Name != newidx.Name {
                    continue
                }

                //fmt.Println("\texist", curidx.Name)

                exist = true

                if newidx.Type != curidx.Type ||
                    !base.ArrayEqual(newidx.Cols, curidx.Cols) {

                    //fmt.Println("\tSchemaIndexSet", newTable.Name, newidx.Name)
                    if err := dc.SchemaIndexSet(dbName, newTable.Name, newidx); err != nil {
                        return err
                    }
                }

                break
            }

            if !exist {
                //fmt.Println("\tSchemaIndexAdd", newTable.Name, newidx.Name)
                err := dc.SchemaIndexAdd(dbName, newTable.Name, newidx)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func (dc *sqlite3Dialect) SchemaDataSet(dbName string) (base.DataSet, error) {

    var err error
    ds := base.DataSet{
        DbName: dbName,
    }

    ds.Tables, err = dc.SchemaTableQuery(dbName)

    return ds, err
}

func (dc *sqlite3Dialect) SchemaTableQuery(dbName string) ([]*base.Table, error) {

    tables := []*base.Table{}

    q := "SELECT name FROM sqlite_master WHERE type='table'"

    rows, err := dc.Base().Conn.Query(q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {

        var name string
        if err = rows.Scan(&name); err != nil {
            return nil, err
        }

        if name == "sqlite_sequence" {
            continue
        }

        idxs, _ := dc.SchemaIndexQuery(dbName, name)

        cols, _ := dc.SchemaColumnQuery(dbName, name)
        // Patch for sqlite3
        for _, v := range cols {
            if base.ArrayContain("PRIMARY", v.Extra) {
                idxs = append(idxs, base.NewIndex("PRIMARY", base.IndexTypePrimaryKey).AddColumn(v.Name))
                break
            }
        }

        tables = append(tables, &base.Table{
            Name:    name,
            Columns: cols,
            Indexes: idxs,
        })
    }

    return tables, nil
}
