package base

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Conn *sql.DB

type Result sql.Result

type Base struct {
	Conn     *sql.DB
	Config   Config
	BaseStmt map[string]string
}

var baseStmt = map[string]string{
	"insertIgnore": "INSERT OR IGNORE INTO `%s` (`%s`) VALUES (%s)",
}

func BaseInit(conf Config, conn *sql.DB) (*Base, error) {
	return &Base{
		Conn:     conn,
		Config:   conf,
		BaseStmt: baseStmt,
	}, nil
}

func (dc *Base) Insert(tableName string, item map[string]interface{}) (Result, error) {

	var res Result

	cols, vars, vals := []string{}, []string{}, []interface{}{}
	for key, val := range item {
		cols = append(cols, key)
		vars = append(vars, "?")
		vals = append(vals, val)
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s)",
		tableName,
		strings.Join(cols, "`,`"),
		strings.Join(vars, ","))

	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	res, err = stmt.Exec(vals...)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (dc *Base) Delete(tableName string, fr Filter) (Result, error) {

	var res Result

	frsql, params := fr.Parse()
	if len(params) == 0 {
		return res, errors.New("Error in query syntax")
	}

	sql := fmt.Sprintf("DELETE FROM `%s` WHERE %s", tableName, frsql)

	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (dc *Base) Update(tableName string, item map[string]interface{}, fr Filter) (Result, error) {

	var res Result

	frsql, params := fr.Parse()
	if len(params) == 0 {
		return res, errors.New("Error in query syntax")
	}

	cols, vals := []string{}, []interface{}{}
	for key, val := range item {
		cols = append(cols, "`"+key+"` = ?")
		vals = append(vals, val)
	}

	vals = append(vals, params...)

	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
		tableName,
		strings.Join(cols, ","),
		frsql)

	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (dc *Base) Count(tableName string, fr Filter) (num int64, err error) {

	frsql, params := fr.Parse()
	hasWhere := "WHERE"
	if len(params) == 0 {
		hasWhere = ""
	}

	sql := fmt.Sprintf("SELECT COUNT(*) FROM `%s` %s %s", tableName, hasWhere, frsql)

	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(params...)
	err = row.Scan(&num)

	return
}

func (dc *Base) InsertIgnore(tableName string, item map[string]interface{}) (Result, error) {

	var res Result

	sqlstmt, ok := dc.BaseStmt["insertIgnore"]
	if !ok {
		return res, errors.New("CurdStmt:insertIgnore missing")
	}

	cols, vars, vals := []string{}, []string{}, []interface{}{}
	for key, val := range item {
		cols = append(cols, key)
		vars = append(vars, "?")
		vals = append(vals, val)
	}

	sql := fmt.Sprintf(sqlstmt, tableName, strings.Join(cols, "`,`"), strings.Join(vars, ","))
	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return res, err
	}
	defer stmt.Close()

	res, err = stmt.Exec(vals...)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (dc *Base) QueryRaw(sql string, params ...interface{}) (rs []Entry, err error) {

	stmt, err := dc.Conn.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	for rows.Next() {

		entry := Entry{Fields: map[string]*Field{}}

		var retvals []interface{}
		for i := 0; i < len(cols); i++ {
			var val interface{}
			retvals = append(retvals, &val)
		}

		if err := rows.Scan(retvals...); err != nil {
			continue
		}

		for ii, key := range cols {

			rawValue := reflect.Indirect(reflect.ValueOf(retvals[ii]))
			if rawValue.Interface() == nil {
				continue
			}

			entry.Fields[key] = &Field{
				valueType: reflect.TypeOf(rawValue.Interface()),
				value:     rawValue,
			}
		}

		rs = append(rs, entry)
	}

	return
}

func (dc *Base) Query(q *QuerySet) (rs []Entry, err error) {

	sql, params := q.Parse()
	if len(params) == 0 {
		return rs, errors.New("Error in query syntax")
	}

	return dc.QueryRaw(sql, params...)
}

func (dc *Base) Fetch(q *QuerySet) (Entry, error) {

	q.Limit(1)

	entry := Entry{Fields: map[string]*Field{}}

	sql, params := q.Parse()
	if len(params) == 0 {
		return entry, errors.New("Error in query syntax")
	}

	rs, err := dc.QueryRaw(sql, params...)
	if err != nil {
		return entry, err
	}

	if len(rs) > 0 {
		return rs[0], nil
	}

	return entry, errors.New("Entry Not Found")
}

func (dc *Base) ExecRaw(query string, args ...interface{}) (Result, error) {
	return dc.Conn.Exec(query, args...)
}
