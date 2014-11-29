package mysql

import (
	"github.com/lessos/lessgo/data/rdo/base"
)

type mysqlDialect struct {
	base *base.Base
}

func (dc *mysqlDialect) Init(base *base.Base) error {
	dc.base = base
	return nil
}

func (dc *mysqlDialect) Base() *base.Base {
	return dc.base
}
