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
	query  string
	params core.ParamStore

	Binder[api.Statement]
	err error
}

func New(driver core.Driver, qc *query.Compiler) *Statement {
	result, err := qc.Compile()
	if err != nil {
		return &Statement{err: err}
	}

	self := &Statement{driver: driver, query: result.QueryString, params: result.Params}
	self.Binder = NewBinder[api.Statement](qc.Params(), self)
	return self
}

func (s *Statement) Prepare(context ...context.Context) (api.Prepared, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.driver == nil {
		return nil, errExecNoDriver
	}

	return NewPrepared(core.OptionalCtx(context), s)
}

// --- Runner

func (s *Statement) Raw() (string, []any, error) {
	if s.err != nil {
		return "", nil, s.err
	}

	vals, err := s.params.Values()
	if err != nil {
		return "", nil, err
	}

	return s.query, vals, nil
}

func (s *Statement) Run(context ...context.Context) error {
	if s.err != nil {
		return s.err
	}

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
