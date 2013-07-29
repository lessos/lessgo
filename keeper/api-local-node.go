package keeper

func (this *Keeper) LocalNodeGet(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "locget",
        "path":   path,
    }
    return this.req("POST", req)
}

func (this *Keeper) LocalNodeList(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "loclist",
        "path":   path,
    }
    return this.req("POST", req)
}

func (this *Keeper) LocalNodeListAndGet(path string) (rpl *Reply) {

    rpl = this.LocalNodeList(path)

    str, err := rpl.Str()
    if err != nil {
        rpl.Type = ReplyError
        return
    }

    var lsis []Node
    if err := JsonDecode(str, &lsis); err != nil {
        rpl.Type = ReplyError
        return
    }

    rpl.Type = ReplyMulti
    for _, v := range lsis {

        if v.T != NodeTypeFile {
            continue
        }

        rs := this.LocalNodeGet(path + "/" + v.P)
        if rs.Type == ReplyError {
            continue
        }

        rpl.Elems = append(rpl.Elems, rs)
    }

    return
}

func (this *Keeper) LocalNodeSet(path, val string) (rpl *Reply) {
    req := map[string]string{
        "method": "locset",
        "path":   path,
        "val":    val,
    }
    return this.req("POST", req)
}

func (this *Keeper) LocalNodeDel(path string) (rpl *Reply) {
    req := map[string]string{
        "method": "locdel",
        "path":   path,
    }
    return this.req("POST", req)
}
