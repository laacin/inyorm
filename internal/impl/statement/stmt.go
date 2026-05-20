package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/driver"
)

type StatementImpl struct {
	Driver driver.Driver

	Query string
	Vals  []any
	Err   error
}

func (s *StatementImpl) Start(driver driver.Driver, builder interface{ Build() (string, []any, error) }) api.Statement {
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
	return s
}
func (s *StatementImpl) BindPrep(binder ...any) api.PrepStatement {
	return s
}

// --- Runner
func (s *StatementImpl) Raw() (string, []any, error) {
	return s.Query, s.Vals, s.Err
}

func (s *StatementImpl) Run(ctx ...context.Context) error {
	return nil
}

// --- Prepare
func (s *StatementImpl) Prepare() api.PrepStatement {
	return s
}

func (s *StatementImpl) Values(values ...any) api.Runner {
	return s
}
