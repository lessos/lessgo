package mysql

import (
	"errors"
	"fmt"
	"github.com/lessos/lessgo/data/rdo/base"
	"sort"
	"strconv"
	"strings"
)

// "SET NAMES 'utf8'"
// "SET CHARACTER_SET_CLIENT=utf8"
// "SET CHARACTER_SET_RESULTS=utf8"

func (dc *mysqlDialect) SchemaIndexAdd(dbName, tableName string, index *base.Index) error {

	action := ""
	switch index.Type {
	case base.IndexTypeIndex:
		action = "INDEX"
	case base.IndexTypeUnique:
		action = "UNIQUE"
	default:
		// PRIMARY KEY can be modified, can not be added
		return errors.New("Invalid Index Type")
	}

	sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` ADD %s `%s` (%s)",
		dbName, tableName, action, index.Name,
		dc.QuoteStr(strings.Join(index.Cols, dc.QuoteStr(","))))

	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaIndexDel(dbName, tableName string, index *base.Index) error {

	// PRIMARY KEY can be modified, can not be deleted
	if index.Type == base.IndexTypePrimaryKey {
		return nil
	}

	sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` DROP INDEX `%s`",
		dbName, tableName, index.Name)

	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaIndexSet(dbName, tableName string, index *base.Index) error {

	dropAction, addAction := "", ""

	switch index.Type {
	case base.IndexTypePrimaryKey:
		dropAction, addAction = "PRIMARY KEY", "PRIMARY KEY"
	case base.IndexTypeIndex:
		dropAction, addAction = "INDEX `"+index.Name+"`", "INDEX `"+index.Name+"`"
	case base.IndexTypeUnique:
		dropAction, addAction = "INDEX `"+index.Name+"`", "UNIQUE `"+index.Name+"`"
	default:
		return errors.New("Invalid Index Type")
	}

	sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` DROP %s, ADD %s (%s)",
		dbName, tableName, dropAction, addAction,
		dc.QuoteStr(strings.Join(index.Cols, dc.QuoteStr(","))))
	//fmt.Println(sql)
	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaIndexQuery(dbName, tableName string) ([]*base.Index, error) {

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

		exist := false
		for i, v := range indexes {

			if v.Name == indexName {
				indexes[i].AddColumn(colName)
				exist = true
			}
		}

		if !exist {
			indexes = append(indexes, base.NewIndex(indexName, indexType).AddColumn(colName))
		}
	}

	return indexes, nil
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
		lens := strings.Split(col.Length, ",")
		if lens[0] == "" {
			lens[0] = "10"
		}
		if len(lens) < 2 {
			lens = append(lens, "2")
		}
		sql = fmt.Sprintf(sql, lens[0], lens[1])
	}

	if col.NullAble {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if col.IncrAble {
		sql += " AUTO_INCREMENT"
	}

	if col.Default != "" {
		if col.Default == "null" {
			sql += " DEFAULT NULL"
		} else {
			sql += " DEFAULT '" + col.Default + "'"
		}
	}

	return dc.QuoteStr(col.Name) + " " + sql
}

func (dc *mysqlDialect) SchemaColumnAdd(dbName, tableName string, col *base.Column) error {

	sql := fmt.Sprintf("ALTER TABLE `%v`.`%v` ADD %v",
		dbName, tableName, dc.SchemaColumnTypeSql(col))

	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaColumnDel(dbName, tableName string, col *base.Column) error {

	sql := fmt.Sprintf("ALTER TABLE `%v`.`%v` DROP `%v`", dbName, tableName, col.Name)

	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaColumnSet(dbName, tableName string, col *base.Column) error {

	sql := fmt.Sprintf("ALTER TABLE `%v`.`%v` CHANGE `%v` %v",
		dbName, tableName, col.Name, dc.SchemaColumnTypeSql(col))

	_, err := dc.Base().ExecRaw(sql)

	return err
}

func (dc *mysqlDialect) SchemaColumnQuery(dbName, tableName string) ([]*base.Column, error) {

	cols := []*base.Column{}

	q := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`, `EXTRA` "
	q += "FROM `INFORMATION_SCHEMA`.`COLUMNS` "
	q += "WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

	rs, err := dc.Base().QueryRaw(q, dbName, tableName)
	if err != nil {
		return cols, err
	}

	for _, entry := range rs {

		col := &base.Column{}

		for name, v := range entry.Fields {
			content := v.String()
			switch name {
			case "COLUMN_NAME":
				col.Name = strings.Trim(content, "` ")
			case "IS_NULLABLE":
				if "YES" == content {
					col.NullAble = true
				}
			case "COLUMN_DEFAULT":
				col.Default = content
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
				if len1 > 0 {
					col.Length = fmt.Sprintf("%v", len1)
					if len2 > 0 {
						col.Length += fmt.Sprintf(",%v", len2)
					}
				}

				typepre := ""
				if strings.Contains(content, "unsigned") {
					typepre = "u"
				}
				switch col.Type {
				case "bigint":
					col.Type = typepre + "int64"
				case "int":
					col.Type = typepre + "int32"
				case "smallint":
					col.Type = typepre + "int16"
				case "tinyint":
					col.Type = typepre + "int8"
				case "varchar":
					col.Type = "string"
				case "longtext":
					col.Type = "string-text"
				}

			case "EXTRA":
				if content == "auto_increment" {
					col.IncrAble = true
				}
			}
		}

		cols = append(cols, col)
	}

	return cols, nil
}

func (dc *mysqlDialect) SchemaTableAdd(table *base.Table) error {

	if len(table.Columns) == 0 {
		return errors.New("No Columns Found")
	}

	sql := "CREATE TABLE IF NOT EXISTS " + dc.QuoteStr(table.Name) + " (\n"

	for _, col := range table.Columns {
		sql += " " + dc.SchemaColumnTypeSql(col) + ",\n"
	}

	pks := []string{}
	for _, idx := range table.Indexes {

		if len(idx.Cols) == 0 {
			continue
		}

		switch idx.Type {
		case base.IndexTypePrimaryKey:
			pks = idx.Cols
			continue
		case base.IndexTypeIndex:
			sql += " KEY "
		case base.IndexTypeUnique:
			sql += " UNIQUE KEY "
		default:
			continue
		}

		sql += dc.QuoteStr(idx.Name) + " ("
		sql += dc.QuoteStr(strings.Join(idx.Cols, dc.QuoteStr(",")))
		sql += "),\n"
	}

	if len(pks) == 0 {
		return errors.New("No PRIMARY KEY Found")
	}
	sql += " PRIMARY KEY ("
	sql += dc.QuoteStr(strings.Join(pks, dc.QuoteStr(",")))
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
	} else {
		sql += " DEFAULT CHARSET=utf8"
	}

	if table.Comment != "" {
		sql += " COMMENT='" + table.Comment + "'"
	}

	sql += ";"

	_, err := dc.Base().ExecRaw(sql)

	return err
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
					newcol.NullAble != curcol.NullAble ||
					newcol.IncrAble != curcol.IncrAble ||
					newcol.Default != curcol.Default {
					colChange = true
					break
				}
			}

			if !colExist {

				if err := dc.SchemaColumnAdd(dbName, newTable.Name, newcol); err != nil {
					return err
				}
			}

			if colChange {

				if err := dc.SchemaColumnSet(dbName, newTable.Name, newcol); err != nil {
					return err
				}
			}
		}

		// Delete Unused Indexes
		for _, curidx := range curTable.Indexes {

			curExist := false

			for _, newidx := range newTable.Indexes {

				if newidx.Name == curidx.Name {
					curExist = true
					break
				}
			}

			if !curExist {
				//fmt.Println("index del", curidx.Name)
				if err := dc.SchemaIndexDel(dbName, newTable.Name, curidx); err != nil {
					return err
				}
			}
		}

		// Delete Unused Columns
		for _, curcol := range curTable.Columns {

			colExist := false

			for _, newcol := range newTable.Columns {

				if newcol.Name == curcol.Name {
					colExist = true
					break
				}
			}

			if !colExist {
				if err := dc.SchemaColumnDel(dbName, newTable.Name, curcol); err != nil {
					return err
				}
			}
		}

		// Add New, or Update Changed Indexes
		for _, newidx := range newTable.Indexes {

			newIdxExist := false
			newIdxChange := false

			for _, curidx := range curTable.Indexes {

				if newidx.Name != curidx.Name {
					continue
				}

				newIdxExist = true

				sort.Strings(newidx.Cols)
				sort.Strings(curidx.Cols)

				if newidx.Type != curidx.Type ||
					strings.Join(newidx.Cols, ",") != strings.Join(curidx.Cols, ",") {

					newIdxChange = true
				}

				break
			}

			if newIdxChange {
				//fmt.Println("index set", newidx.Name)
				if err := dc.SchemaIndexSet(dbName, newTable.Name, newidx); err != nil {
					return err
				}

			} else if !newIdxExist {
				//fmt.Println("index add", newidx.Name)
				if err := dc.SchemaIndexAdd(dbName, newTable.Name, newidx); err != nil {
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

	ds.Tables, err = dc.SchemaTableQuery(dbName)

	return ds, err
}

func (dc *mysqlDialect) SchemaTableQuery(dbName string) ([]*base.Table, error) {

	tables := []*base.Table{}

	q := "SELECT `TABLE_NAME`, `ENGINE`, `TABLE_COLLATION`, `TABLE_COMMENT` "
	q += "FROM `INFORMATION_SCHEMA`.`TABLES` "
	q += "WHERE `TABLE_SCHEMA` = ?"

	rows, err := dc.Base().Conn.Query(q, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var name, engine, charset, comment string
		if err = rows.Scan(&name, &engine, &charset, &comment); err != nil {
			return nil, err
		}

		if i := strings.Index(charset, "_"); i > 0 {
			charset = charset[0:i]
		}

		idxs, _ := dc.SchemaIndexQuery(dbName, name)

		cols, _ := dc.SchemaColumnQuery(dbName, name)

		tables = append(tables, &base.Table{
			Name:    name,
			Engine:  engine,
			Charset: charset,
			Columns: cols,
			Indexes: idxs,
			Comment: comment,
		})
	}

	return tables, nil
}
