package internal

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/mapper"
	"github.com/laacin/inyorm/internal/query"
)

type Statement struct {
	driver core.Driver

	query query.QueryBuilder
	bind  any
}

func (s *Statement) Start(driver core.Driver, q query.QueryBuilder) api.Statement {
	s.driver = driver
	s.query = q
	return s
}

// --- Binder
func (s *Statement) Bind(binder ...any) api.Statement {
	if len(binder) > 0 {
		s.bind = binder[0]
	}
	return s
}
func (s *Statement) BindPrep(binder ...any) api.PrepStatement {
	if len(binder) > 0 {
		s.bind = binder[0]
	}
	return s
}

// --- Runner
func (s *Statement) Raw() (string, []any, error) {
	result, err := s.query.Build()
	if err != nil {
		return "", nil, err
	}

	return result.Query, result.Values, nil
}

func (s *Statement) Run(context ...context.Context) error {
	result, err := s.query.Build()
	if err != nil {
		return err
	}

	ctx := getCtx(context)
	if s.bind == nil {
		return s.driver.Exec(ctx, result.Query, result.Values...)
	}

	rows, err := s.driver.Query(ctx, result.Query, result.Values...)
	if err != nil {
		return err
	}

	return mapper.Scan(rows, s.bind)
}

func (s *Statement) RunTx(ctx context.Context, tx core.Transaction) error {
	result, err := s.query.Build()
	if err != nil {
		return err
	}

	if s.bind == nil {
		return tx.Exec(ctx, result.Query, result.Values...)
	}

	rows, err := tx.Query(ctx, result.Query, result.Values...)
	if err != nil {
		return err
	}

	return mapper.Scan(rows, s.bind)
}

// --- Prepare
func (s *Statement) Prepare() api.PrepStatement {
	panic("TODO")
}

func (s *Statement) Values(values ...any) api.Runner {
	panic("TODO")
}

// --- Helpers

func getCtx(candidate []context.Context) context.Context {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return context.Background()
}
