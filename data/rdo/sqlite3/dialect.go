package sqlite3

import (
	"github.com/lessos/lessgo/data/rdo/base"
)

type sqlite3Dialect struct {
	base *base.Base
}

func (dc *sqlite3Dialect) Init(base *base.Base) error {
	dc.base = base
	return nil
}

func (dc *sqlite3Dialect) Base() *base.Base {
	return dc.base
}
