package base

import (
    "database/sql"
    "errors"
    "fmt"
    "reflect"
    "strings"
    "time"
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

func (dc *Base) QueryRaw(sql string, params ...interface{}) (rs []map[string]interface{}, err error) {

    //fmt.Println("sql", sql, params)
    stmt, err := dc.Conn.Prepare(sql)
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

func (dc *Base) Query(q *QuerySet) (rs []map[string]interface{}, err error) {

    sql, params := q.Parse()
    if len(params) == 0 {
        return rs, errors.New("Error in query syntax")
    }

    return dc.QueryRaw(sql, params...)
}

func (dc *Base) ExecRaw(query string, args ...interface{}) (Result, error) {
    return dc.Conn.Exec(query, args...)
}
