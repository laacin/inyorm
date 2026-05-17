package psql

import (
	"strconv"

	"github.com/laacin/inyorm/engine/std"
	"github.com/laacin/inyorm/internal/core"
)

type PsqlDialect struct {
	std.Dialect
}

func (dial *PsqlDialect) WritePlaceholder(w core.Writer) {
	w.Char('$')
	w.Write(strconv.Itoa(w.ValueCount()))
}
