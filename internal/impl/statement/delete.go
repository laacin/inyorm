package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/driver"
	"github.com/laacin/inyorm/internal/execution"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type DeleteStmtImpl struct {
	DefaultRef string
	Dialect    dml.Dialect

	clause.DeleteImpl
	clause.FromImpl
	clause.WhereImpl

	*execution.Executor
}

func NewDeleteStatement(ctx context.Context, dial dml.Dialect, driver driver.Driver, ref string) *DeleteStmtImpl {
	stmt := &DeleteStmtImpl{Dialect: dial, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: driver}
	stmt.Executor = exec
	return stmt
}

func (s *DeleteStmtImpl) Kind() dml.StatementKind {
	return dml.StatementDelete
}

func (s *DeleteStmtImpl) Build() (*dml.Statement, error) {
	// Auto-FROM
	if !s.FromImpl.IsDeclared() && s.DefaultRef != "" {
		s.FromImpl.From(&dml.Table{Value: s.DefaultRef})
	}

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.DeleteImpl,
		&s.FromImpl,
		&s.WhereImpl,
	}

	clauseMap := make(map[dml.ClauseKind]dml.Clause)
	for _, cls := range clauses {
		if cls.IsDeclared() {
			builded, err := cls.Build()
			if err != nil {
				return nil, err
			}

			clauseMap[cls.Kind()] = builded
		}
	}

	// --- Declarate writers

	var (
		parameters = &writer.ParamStore{}
	)

	w := &writer.WriterImpl{
		Syntax: s.Dialect,
		Params: parameters,
	}

	// --- Write the statement

	first := true
	for _, ord := range s.Dialect.DeleteOrder() {
		if clause, ok := clauseMap[ord]; ok {
			if !first {
				w.Char(' ')
			}
			first = false
			clause.Write(w, s.Dialect)
		}
	}

	// --- Validate values

	if err := parameters.Validate(); err != nil {
		return nil, err
	}

	return &dml.Statement{
		Query:  w.ToString(),
		Values: parameters.Values(),
	}, nil
}
