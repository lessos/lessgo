package hissdb

import (
	"github.com/lessos/lessgo/utils"
	"strconv"
)

const (
	ReplyOK          string = "ok"
	ReplyNotFound    string = "not_found"
	ReplyError       string = "error"
	ReplyFail        string = "fail"
	ReplyClientError string = "client_error"
)

type Reply struct {
	State string
	Data  []string
}

type Entry struct {
	Key, Value string
}

func (r *Reply) String() string {

	if len(r.Data) > 0 {
		return r.Data[0]
	}

	return ""
}

func (r *Reply) Int() int {
	return int(r.Int64())
}

func (r *Reply) Int64() int64 {

	if len(r.Data) < 1 {
		return 0
	}

	i64, err := strconv.ParseInt(r.Data[0], 10, 64)
	if err == nil {
		return i64
	}

	return 0
}

func (r *Reply) Uint() uint {
	return uint(r.Uint64())
}

func (r *Reply) Uint64() uint64 {

	if len(r.Data) < 1 {
		return 0
	}

	i64, err := strconv.ParseUint(r.Data[0], 10, 64)
	if err == nil {
		return i64
	}

	return 0
}

func (r *Reply) Float64() float64 {

	if len(r.Data) < 1 {
		return 0
	}

	f64, err := strconv.ParseFloat(r.Data[0], 64)
	if err == nil {
		return f64
	}

	return 0
}

func (r *Reply) Bool() bool {

	if len(r.Data) < 1 {
		return false
	}

	b, err := strconv.ParseBool(r.Data[0])
	if err == nil {
		return b
	}

	return false
}

func (r *Reply) List() []string {

	if len(r.Data) < 1 {
		return []string{}
	}

	return r.Data
}

func (r *Reply) Hash() []Entry {

	hs := []Entry{}

	if len(r.Data) < 2 {
		return hs
	}

	for i := 0; i < (len(r.Data) - 1); i += 2 {
		hs = append(hs, Entry{r.Data[i], r.Data[i+1]})
	}

	return hs
}

// Json returns the map that marshals from the reply bytes as json in response .
func (r *Reply) Json(v interface{}) error {
	return utils.JsonDecode(r.String(), v)
}

func (r *Entry) Json(v interface{}) error {
	return utils.JsonDecode(r.Value, v)
}
