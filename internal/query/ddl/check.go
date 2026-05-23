package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

type Check struct{ Cond expr.ExprBuilder }

// start

func (b *Check) Start(ident any) api.Cond {
	cond := &expr.Cond{}
	b.Cond = cond
	return cond.Start(ident)
}

// --- Build

func (b *Check) Build(w core.InternalWriter, dial TableWriter) {
	dial.WriteConsCheck(w, b)
}
