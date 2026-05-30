package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

type Statement struct {
	driver core.Driver
	qc     *query.Compiler

	Binder[api.Statement]
}

func New(driver core.Driver, qc *query.Compiler) *Statement {
	self := &Statement{driver: driver, qc: qc}
	self.Binder = NewBinder[api.Statement](qc.Params(), self)
	return self
}

// --- Runner

func (s *Statement) Raw() (string, []any, error) {
	return Raw(s)
}

func (s *Statement) Run(context ...context.Context) error {
	ctx := core.OptionalCtx(context)
	return Run(s, ctx)
}
