package keeper

const (
    NodeTypeNil  = uint8(0)
    NodeTypeDir  = uint8(1)
    NodeTypeFile = uint8(2)
)

type Node struct {
    P   string
    C   string
    R   uint64
    T   uint8
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

func (this *Keeper) NodeListAndGet(path string) (rpl *Reply) {

    rpl = this.NodeList(path)

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

        rs := this.NodeGet(path + "/" + v.P)
        if rs.Type == ReplyError {
            continue
        }

        rpl.Elems = append(rpl.Elems, rs)
    }

    return
}

func (this *Keeper) NodeSet(path, val string) (rpl *Reply) {
    req := map[string]string{
        "method": "set",
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
