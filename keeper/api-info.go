package keeper

type KprInfo struct {
    Leader  string          `json:"leader"`
    Vote    uint64          `json:"vote"`
    Members []KprInfoMember `json:"members"`
    Local   KprInfoLocal    `json:"local"`
}
type KprInfoMember struct {
    Id     string `json:"id"`
    Seat   int    `json:"seat"`
    Addr   string `json:"addr"`
    Port   string `json:"port"`
    Status int    `json:"status"`
}
type KprInfoLocal struct {
    Id         string `json:"id"`
    Addr       string `json:"addr"`
    KeeperPort string `json:"keeperport"`
    AgentPort  string `json:"agentport"`
    Status     int    `json:"status"`
}

func (this *Keeper) SysInfo() (rpl *Reply) {
    req := map[string]string{
        "method": "info",
    }
    return this.req("POST", req)
}
func (this *Keeper) SysInfoWithObj() (rpl *KprInfo, err error) {

    rs := this.SysInfo()

    if err := JsonDecode(rs.Body, &rpl); err != nil {
        return nil, err
    }

    return rpl, nil
}
