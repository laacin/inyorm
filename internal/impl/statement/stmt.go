package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/mapper"
)

type StatementBuilder interface{ Build() (string, []any, error) }

type StatementImpl struct {
	Driver core.Driver

	stmt StatementBuilder
	bind any

	Query string
	Vals  []any
	Err   error
}

func (s *StatementImpl) Start(driver core.Driver, builder StatementBuilder) api.Statement {
	q, vals, err := builder.Build()
	if err != nil {
		s.Err = err
		return s
	}

	s.Driver = driver
	s.Query = q
	s.Vals = vals
	return s
}

// --- Binder
func (s *StatementImpl) Bind(binder ...any) api.Statement {
	if len(binder) > 0 {
		s.bind = binder[0]
	}
	return s
}
func (s *StatementImpl) BindPrep(binder ...any) api.PrepStatement {
	if len(binder) > 0 {
		s.bind = binder[0]
	}
	return s
}

// --- Runner
func (s *StatementImpl) Raw() (string, []any, error) {
	return s.Query, s.Vals, s.Err
}

func (s *StatementImpl) Run(context ...context.Context) error {
	if s.Err != nil {
		return s.Err
	}

	ctx := getCtx(context)
	if s.bind == nil {
		return s.Driver.Exec(ctx, s.Query, s.Vals...)
	}

	rows, err := s.Driver.Query(ctx, s.Query, s.Vals...)
	if err != nil {
		return err
	}

	return mapper.Scan(rows, s.bind)
}

// --- Prepare
func (s *StatementImpl) Prepare() api.PrepStatement {
	return s
}

func (s *StatementImpl) Values(values ...any) api.Runner {
	return s
}

// --- Helpers

func getCtx(candidate []context.Context) context.Context {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return context.Background()
}
