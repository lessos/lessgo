package keeper

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
)

type Keeper struct {
    client *http.Client
    apiuri string
}

func NewKeeper(dsn string) (Keeper, error) {
    var kpr Keeper

    tr := &http.Transport{
        DisableKeepAlives: true,
        //MaxIdleConnsPerHost: 100,
    }
    kpr.client = &http.Client{Transport: tr}
    kpr.apiuri = "http://" + dsn + "/h5keeper/api"

    if false {
        fmt.Println("DDDDDDDDD")
    }

    return kpr, nil
}

func (this *Keeper) req(m string, req interface{}) (rpl *Reply) {

    rqj, err := JsonEncode(req)
    if err != nil {
        rpl.Type = ReplyError
        return
    }

    body := bytes.NewBufferString(rqj)

    rq, err := http.NewRequest(m, this.apiuri, body)
    if err != nil {
        return
    }

    rsp, err := this.client.Do(rq)
    if err != nil {
        return
    }
    defer rsp.Body.Close()

    rpb, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        return
    }

    err = JsonDecode(string(rpb), &rpl)
    if err != nil {
        return
    }

    return rpl
}
