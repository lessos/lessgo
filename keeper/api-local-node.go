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
