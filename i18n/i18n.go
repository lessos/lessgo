package i18n

import (
    "../utils"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
)

var (
    i18n   map[string]string = map[string]string{}
    Locale string            = "en_US"
)

type config struct {
    Locale string `json:"locale"`
    Data   []configItem
}
type configItem struct {
    Key string `json:"key"`
    Val string `json:"val"`
}

func Config(path string) {

    var cfg config

    str, err := fsFileGetRead(path)
    if err != nil {
        return
    }

    if err := utils.JsonDecode(str, &cfg); err != nil {
        fmt.Println("Format Error:", err)
        return
    }

    if Locale == cfg.Locale {
        Locale = cfg.Locale
    }

    for _, v := range cfg.Data {

        key := strings.ToLower(cfg.Locale + "." + v.Key)

        if v2, ok := i18n[key]; !ok || v2 != v.Val {
            i18n[key] = v.Val
        }
    }

    //fmt.Println(i18n)
}

func T(key string) string {

    key = strings.ToLower(Locale + "." + key)

    if v, ok := i18n[key]; ok {
        return v
    }

    return key
}

func fsFileGetRead(path string) (string, error) {

    reg, _ := regexp.Compile("/+")
    path = "/" + strings.Trim(reg.ReplaceAllString(path, "/"), "/")

    st, err := os.Stat(path)
    if err != nil || os.IsNotExist(err) {
        return "", errors.New("File Not Found")
    }

    if st.Size() > (10 * 1024 * 1024) {
        return "", errors.New("File size is too large")
    }

    fp, err := os.OpenFile(path, os.O_RDWR, 0754)
    if err != nil {
        return "", errors.New("File Can Not Open")
    }
    defer fp.Close()

    ctn, err := ioutil.ReadAll(fp)
    if err != nil {
        return "", errors.New("File Can Not Readable")
    }

    return string(ctn), nil
}
