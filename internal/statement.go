package internal

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
)

type Statement struct {
	rend   expr.Renderer
	query  *query.Compiler
	driver core.Driver

	bind any
}

func (s *Statement) Start(driver core.Driver, rend expr.Renderer, qc *query.Compiler) api.Statement {
	s.rend = rend
	s.driver = driver
	s.query = qc
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
	result, err := s.query.Compile(s.rend)
	if err != nil {
		return "", nil, err
	}

	vals, err := result.Params.Values()
	if err != nil {
		return "", nil, err
	}

	return result.QueryString, vals, nil
}

func (s *Statement) Run(context ...context.Context) error {
	query, values, err := s.Raw()
	if err != nil {
		return err
	}

	ctx := getCtx(context)
	if s.bind == nil {
		return s.driver.Exec(ctx, query, values...)
	}

	rows, err := s.driver.Query(ctx, query, values...)
	if err != nil {
		return err
	}

	return mapper.New().Scan(rows, s.bind)
}

func (s *Statement) RunTx(ctx context.Context, tx core.Transaction) error {
	query, values, err := s.Raw()
	if err != nil {
		return err
	}

	if s.bind == nil {
		return tx.Exec(ctx, query, values...)
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return err
	}

	return mapper.New().Scan(rows, s.bind)
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
	if len(candidate) > 0 && candidate[0] != nil {
		return candidate[0]
	}
	return context.Background()
}
