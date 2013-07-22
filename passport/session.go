package passport

import (
    "../keeper"
    "../utils"
    //"fmt"
)

type Session struct {
    Kpr keeper.Keeper
}

type SessionInstance struct {
    Id    string `json:"id"`
    Uid   string `json:"uid"`
    Uname string `json:"uname"`
}

func NewSession(kpr keeper.Keeper) (Session, error) {
    var si Session

    si.Kpr = kpr

    return si, nil
}

func (this *Session) Instance(token string) SessionInstance {

    rs := this.Kpr.LocalNodeGet("/u/s/" + token)

    var rpl SessionInstance
    utils.JsonDecode(rs.Body, &rpl)

    return rpl
}

func (this *Session) IsLogin(token string) bool {

    ins := this.Instance(token)
    if ins.Uid == "0" || ins.Uid == "" {
        return false
    }

    return true
}
