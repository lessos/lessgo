package utils

import (
    "encoding/json"
    "errors"
)

func JsonDecode(src, rs interface{}) (err error) {

    defer func() {
        if r := recover(); r != nil {
            err = errors.New("json: invalid format")
        }
    }()

    var bf []byte
    switch src.(type) {
    case string:
        bf = []byte(src.(string))
    case []byte:
        bf = src.([]byte)
    default:
        panic("invalid format")
    }

    if err = json.Unmarshal(bf, &rs); err != nil {
        return err
    }

    return nil
}

func JsonEncode(rs interface{}) (str string, err error) {

    rb, err := json.Marshal(rs)

    if err == nil {
        str = string(rb)
    }

    return
}
