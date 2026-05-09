package dml

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/execution"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type InsertStmtImpl struct {
	DefaultRef string
	Dialect    entity.Dialect

	clause.InsertIntoImpl

	*execution.Executor
}

func NewInsertStatement(ctx context.Context, dial entity.Dialect, driver entity.Driver, ref string) *InsertStmtImpl {
	stmt := &InsertStmtImpl{Dialect: dial, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: driver}
	stmt.Executor = exec
	return stmt
}

func (s *InsertStmtImpl) Kind() entity.StatementKind {
	return entity.StatementInsert
}

func (s *InsertStmtImpl) Build() (*entity.Statement, error) {
	// Auto-FROM
	s.InsertIntoImpl.Table(&entity.Table{Value: s.DefaultRef})

	// --- Load clauses
	clauses := []entity.ClauseBuilder{
		&s.InsertIntoImpl,
	}

	clauseMap := make(map[entity.ClauseKind]entity.Clause)
	for _, cls := range clauses {
		if cls.IsDeclared() {
			clauseMap[cls.Kind()] = cls.Build()
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
	for _, ord := range s.Dialect.InsertOrder() {
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
