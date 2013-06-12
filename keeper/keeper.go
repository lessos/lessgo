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


type Node struct {
    P   string
    C   string
    R   uint64
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

func (this *Keeper) NodeGet(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "get",
        "path":   path,
    }
    return this.req("POST", req)
}

func (this *Keeper) NodeList(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "list",
        "path":   path,
    }
    return this.req("POST", req)
}

func (this *Keeper) NodeSet(path, val string) (rpl *Reply) {
    req := map[string]string{
        "method": "list",
        "path":   path,
        "val":    val,
    }
    return this.req("POST", req)
}

func (this *Keeper) NodeDel(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "del",
        "path":   path,
    }
    return this.req("POST", req)
}
