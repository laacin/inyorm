package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
)

type Check struct{ Cond expr.Expr }

// start

func (b *Check) Start(ident any) api.Cond {
	cond := &expr.Cond{}
	b.Cond = cond
	return cond.Start(ident)
}

// --- Build

func (b *Check) Build() error {
	return nil
}
