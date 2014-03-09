package rdc

import (
    "database/sql"
    "errors"
    //"fmt"
    _ "github.com/mattn/go-sqlite3"
    "reflect"
    "time"
)

type Config struct {
    Driver string
    DbPath string
}

var configDrivers = map[string]bool{
    "sqlite3": true,
}

type Conn struct {
    db *sql.DB
}

func NewConfig() Config {
    return Config{}
}

func (c Config) Instance() (*Conn, error) {

    var err error

    if !configDrivers[c.Driver] {
        return nil, errors.New("Driver can not found")
    }

    var cn Conn

    cn.db, err = sql.Open(c.Driver, c.DbPath)
    if err != nil {
        return nil, err
    }

    return &cn, nil
}

func (cn *Conn) Close() {
    cn.db.Close()
}

func (cn *Conn) Insert(tblname string, item interface{}) error {
    return nil
}

func (cn *Conn) Delete(tblname string, where interface{}) error {
    return nil
}

func (cn *Conn) Update(tblname string, item interface{}, where interface{}) error {
    return nil
}

func (cn *Conn) Query(q *QuerySet) (rs []map[string]interface{}, err error) {

    sql, params := q.Parse()
    if len(params) == 0 {
        return rs, errors.New("Error in query syntax")
    }

    stmt, err := cn.db.Prepare(sql)
    if err != nil {
        return
    }
    defer stmt.Close()

    rows, err2 := stmt.Query(params...)
    if err2 != nil {
        return
    }

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
                vi = rawValue.Interface().(time.Time).Format("2006-01-02 15:04:05.000 -0700")
            }

            ret[key] = vi
        }

        rs = append(rs, ret)

    }

    return
}
