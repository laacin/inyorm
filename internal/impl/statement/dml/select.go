package dml

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/execution"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type SelectStmtImpl struct {
	DefaultRef string
	Dialect    entity.Dialect

	clause.SelectImpl
	clause.FromImpl
	clause.JoinImpl
	clause.WhereImpl
	clause.GroupByImpl
	clause.HavingImpl
	clause.OrderByImpl
	clause.LimitImpl
	clause.OffsetImpl

	*execution.Executor
}

func NewSelectStatement(ctx context.Context, dial entity.Dialect, driver entity.Driver, ref string) *SelectStmtImpl {
	stmt := &SelectStmtImpl{Dialect: dial, DefaultRef: ref}
	exec := &execution.Executor{Ctx: ctx, Statement: stmt, Driver: driver}
	stmt.Executor = exec
	return stmt
}

func (s *SelectStmtImpl) Kind() entity.StatementKind {
	return entity.StatementSelect
}

func (s *SelectStmtImpl) Build() (*entity.Statement, error) {
	// Auto-FROM
	if !s.FromImpl.IsDeclared() && s.DefaultRef != "" {
		s.FromImpl.From(&entity.Table{Value: s.DefaultRef})
	}

	// --- Load clauses
	clauses := []entity.ClauseBuilder{
		&s.SelectImpl,
		&s.FromImpl,
		&s.JoinImpl,
		&s.WhereImpl,
		&s.GroupByImpl,
		&s.HavingImpl,
		&s.OrderByImpl,
		&s.LimitImpl,
		&s.OffsetImpl,
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
		aliases    *writer.AliasStore
	)

	// --- Set table references if Join exists

	if s.JoinImpl.IsDeclared() {
		aliases = &writer.AliasStore{}
	}

	// --- Write the statement

	w := &writer.WriterImpl{
		Syntax:  s.Dialect,
		Params:  parameters,
		Aliases: aliases,
	}

	w.SetRef(s.DefaultRef)
	first := true
	for _, ord := range s.Dialect.SelectOrder() {
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
