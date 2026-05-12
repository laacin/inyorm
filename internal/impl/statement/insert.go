package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/driver"
	"github.com/laacin/inyorm/internal/execution"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type InsertStmtImpl struct {
	DefaultRef string
	Dialect    dml.Dialect

	clause.InsertIntoImpl

	*execution.Executor
}

func NewInsertStatement(ctx context.Context, dial dml.Dialect, driver driver.Driver, ref string) *InsertStmtImpl {
	stmt := &InsertStmtImpl{Dialect: dial, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: driver}
	stmt.Executor = exec
	return stmt
}

func (s *InsertStmtImpl) Kind() dml.StatementKind {
	return dml.StatementInsert
}

func (s *InsertStmtImpl) Build() (*dml.Statement, error) {
	// Auto-FROM
	s.InsertIntoImpl.Table(&dml.Table{Value: s.DefaultRef})

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.InsertIntoImpl,
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

	return &dml.Statement{
		Query:  w.ToString(),
		Values: parameters.Values(),
	}, nil
}
