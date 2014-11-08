package hissdb

import (
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
	return r.Int64() == 1
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
