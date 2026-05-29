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

	bind any
}

func New(driver core.Driver, qc *query.Compiler) *Statement {
	return &Statement{
		driver: driver,
		qc:     qc,
	}
}

// --- Binder

func (s *Statement) Bind(v any) api.Statement {
	s.bind = v
	return s
}

func (s *Statement) Values(v any) api.Statement {
	return s
}

// --- Runner

func (s *Statement) Raw() (string, []any, error) {
	return Raw(s)
}

func (s *Statement) Run(context ...context.Context) error {
	ctx := core.OptionalCtx(context)
	return Run(s, ctx)
}
