package rdc

import (
    "errors"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "reflect"
    "strings"
    "time"
)

func (cn *Conn) Insert(tblname string, item map[string]interface{}) (Result, error) {

    var res Result

    cols, vars, vals := []string{}, []string{}, []interface{}{}
    for key, val := range item {
        cols = append(cols, key)
        vars = append(vars, "?")
        vals = append(vals, val)
    }

    sql := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s)",
        tblname,
        strings.Join(cols, "`,`"),
        strings.Join(vars, ","))

    stmt, err := cn.db.Prepare(sql)
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

func (cn *Conn) Delete(tblname string, fr Filter) (Result, error) {

    var res Result

    frsql, params := fr.Parse()
    if len(params) == 0 {
        return res, errors.New("Error in query syntax")
    }

    sql := fmt.Sprintf("DELETE FROM `%s` WHERE %s", tblname, frsql)

    stmt, err := cn.db.Prepare(sql)
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

func (cn *Conn) Update(tblname string, item map[string]interface{}, fr Filter) (Result, error) {

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
        tblname,
        strings.Join(cols, ","),
        frsql)

    stmt, err := cn.db.Prepare(sql)
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

func (cn *Conn) Count(tblname string, fr Filter) (num int64, err error) {

    frsql, params := fr.Parse()
    hasWhere := "WHERE"
    if len(params) == 0 {
        hasWhere = ""
    }

    sql := fmt.Sprintf("SELECT COUNT(*) FROM `%s` %s %s", tblname, hasWhere, frsql)

    stmt, err := cn.db.Prepare(sql)
    if err != nil {
        return
    }
    defer stmt.Close()

    row := stmt.QueryRow(params...)
    err = row.Scan(&num)

    return
}

func (cn *Conn) QueryRaw(sql string, params ...interface{}) (rs []map[string]interface{}, err error) {

    //fmt.Println("sql", sql, params)
    stmt, err := cn.db.Prepare(sql)
    if err != nil {
        return
    }
    defer stmt.Close()

    rows, err2 := stmt.Query(params...)
    if err2 != nil {
        return
    }
    defer rows.Close()

    fields, err3 := rows.Columns()
    if err3 != nil {
        return
    }

    for rows.Next() {

        ret := map[string]interface{}{}

        var retvals []interface{}
        for i := 0; i < len(fields); i++ {
            var val interface{}
            retvals = append(retvals, &val)
        }

        if err := rows.Scan(retvals...); err != nil {
            continue
        }

        for ii, key := range fields {

            rawValue := reflect.Indirect(reflect.ValueOf(retvals[ii]))

            if rawValue.Interface() == nil {
                continue
            }

            aa := reflect.TypeOf(rawValue.Interface())
            vv := reflect.ValueOf(rawValue.Interface())

            var vi interface{}
            switch aa.Kind() {
            case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                vi = vv.Int()
            case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                vi = vv.Uint()
            case reflect.Float32, reflect.Float64:
                vi = vv.Float()
            case reflect.Slice:
                if aa.Elem().Kind() == reflect.Uint8 {
                    vi = string(rawValue.Interface().([]byte))
                }
            case reflect.String:
                vi = vv.String()
            case reflect.Struct:
                if aa.String() == "time.Time" {
                    vi = rawValue.Interface().(time.Time).In(TimeZone)
                }
            }

            ret[key] = vi
        }

        rs = append(rs, ret)

    }

    return
}

func (cn *Conn) Query(q *QuerySet) (rs []map[string]interface{}, err error) {

    sql, params := q.Parse()
    if len(params) == 0 {
        return rs, errors.New("Error in query syntax")
    }

    return cn.QueryRaw(sql, params...)
}

func (cn *Conn) ExecRaw(query string, args ...interface{}) (Result, error) {
    return cn.db.Exec(query, args...)
}
