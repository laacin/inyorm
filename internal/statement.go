package internal

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/builder/mapper"
	"github.com/laacin/inyorm/internal/core"
)

type queryBuilder interface {
	Build() (*builder.Builder, error)
}

type Statement struct {
	query  queryBuilder
	driver core.Driver

	bind any
}

func (s *Statement) Start(driver core.Driver, q queryBuilder) api.Statement {
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
	b, err := s.query.Build()
	if err != nil {
		return "", nil, err
	}

	vals, err := b.Params().Values()
	if err != nil {
		return "", nil, err
	}

	return b.Writer().ToString(), vals, nil
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

	return (&mapper.Mapper{}).Scan(rows, s.bind)
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

	return (&mapper.Mapper{}).Scan(rows, s.bind)
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
