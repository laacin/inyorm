package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/expr"
	"github.com/laacin/inyorm/internal/execution"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type UpdateStmtImpl struct {
	DefaultRef string
	Dialect    entity.Dialect

	clause.UpdateImpl
	clause.WhereImpl

	*execution.Executor
}

func NewUpdateStatement(ctx context.Context, eng *entity.Engine, ref string) *UpdateStmtImpl {
	stmt := &UpdateStmtImpl{Dialect: eng.Dialect, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: eng.Driver}
	stmt.Executor = exec
	return stmt
}

func (s *UpdateStmtImpl) Kind() dml.StatementKind {
	return dml.StatementSelect
}

func (s *UpdateStmtImpl) Build() (*dml.Statement, error) {
	s.UpdateImpl.Table(&expr.Table{Value: s.DefaultRef})

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.UpdateImpl,
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
	for _, ord := range s.Dialect.UpdateOrder() {
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
