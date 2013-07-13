package utils

import (
    "encoding/json"
    "errors"
)

func JsonDecode(str string, rs interface{}) (err error) {

    defer func() {
        if r := recover(); r != nil {
            err = errors.New("json: invalid format")
        }
    }()

    if err = json.Unmarshal([]byte(str), &rs); err != nil {
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
