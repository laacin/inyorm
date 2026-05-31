package ddl

import "github.com/laacin/inyorm/internal/expr"

type Check struct{ Cond expr.Expr }

func NewCheck(cond *expr.Cond) *Check {
	return &Check{Cond: cond}
}

// --- Build

func (b *Check) Build() error {
	return nil
}
