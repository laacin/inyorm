package dml

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
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

func NewUpdateStatement(ctx context.Context, dial entity.Dialect, driver entity.Driver, ref string) *UpdateStmtImpl {
	stmt := &UpdateStmtImpl{Dialect: dial, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: driver}
	stmt.Executor = exec
	return stmt
}

func (s *UpdateStmtImpl) Kind() entity.StatementKind {
	return entity.StatementSelect
}

func (s *UpdateStmtImpl) Build() (*entity.Statement, error) {
	s.UpdateImpl.Table(&entity.Table{Value: s.DefaultRef})

	// --- Load clauses
	clauses := []entity.ClauseBuilder{
		&s.UpdateImpl,
		&s.WhereImpl,
	}

	clauseMap := make(map[entity.ClauseKind]entity.Clause)
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

	return &entity.Statement{
		Query:  w.ToString(),
		Values: parameters.Values(),
	}, nil
}
