package statement

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
	"github.com/laacin/inyorm/internal/query"
)

var errExecNoDriver = errors.New("missing driver")

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

func (s *Statement) Prepare(context ...context.Context) (api.Prepared, error) {
	return NewPrepared(core.OptionalCtx(context), s.driver, s.qc)
}

// --- Runner

func (s *Statement) Raw() (string, []any, error) {
	result, err := s.qc.Compile()
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
	if s.driver == nil {
		return errExecNoDriver
	}

	query, vals, err := s.Raw()
	if err != nil {
		return err
	}

	ctx := core.OptionalCtx(context)
	if s.bind == nil {
		return s.driver.Exec(ctx, query, vals...)
	}

	rows, err := s.driver.Query(ctx, query, vals...)
	if err != nil {
		return err
	}

	return mapper.New().Bind(rows, s.bind)
}
