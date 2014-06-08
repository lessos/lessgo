package base

import (
    "../../../utils"
    "io/ioutil"
    "os"
)

type DataSet struct {
    DbName  string  `json:"dbname"`
    Engine  string  `json:"engine"`
    Charset string  `json:"charset"`
    Tables  []Table `json:"tables"`
}

func LoadDataSetConfig(file string) (DataSet, error) {

    var ds DataSet
    var err error

    if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
        return ds, err
    }

    fp, err := os.Open(file)
    if err != nil {
        return ds, err
    }
    defer fp.Close()

    cfg, err := ioutil.ReadAll(fp)
    if err != nil {
        return ds, err
    }

    if err = utils.JsonDecode(cfg, &ds); err != nil {
        return ds, err
    }

    return ds, err
}
