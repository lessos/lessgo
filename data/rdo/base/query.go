package base

import (
	"fmt"
)

type QuerySet struct {
	cols   string
	table  string
	order  string
	limit  int64
	offset int64
	Where  Filter
}

func NewQuerySet() *QuerySet {
	return &QuerySet{
		cols:   "*",
		limit:  1,
		offset: 0,
		Where:  NewFilter(),
	}
}

func (q *QuerySet) Select(s string) *QuerySet {
	q.cols = s
	return q
}

func (q *QuerySet) From(s string) *QuerySet {
	q.table = s
	return q
}

func (q *QuerySet) Order(s string) *QuerySet {
	q.order = s
	return q
}

func (q *QuerySet) Limit(num int64) *QuerySet {
	q.limit = num
	return q
}

func (q *QuerySet) Offset(num int64) *QuerySet {
	q.offset = num
	return q
}

func (q *QuerySet) Parse() (sql string, params []interface{}) {

	if len(q.table) == 0 {
		return
	}

	sql = fmt.Sprintf("SELECT %s FROM %s ", q.cols, q.table)

	frsql, ps := q.Where.Parse()
	if len(ps) > 0 {
		sql += "WHERE " + frsql + " "
		params = ps
	}

	if len(q.order) > 0 {
		sql += "ORDER BY " + q.order + " "
	}

	if q.offset > 0 {
		sql += "LIMIT ?,?"
		params = append(params, q.offset, q.limit)
	} else {
		sql += "LIMIT ?"
		params = append(params, q.limit)
	}

	return
}
