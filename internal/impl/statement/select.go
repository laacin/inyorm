package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/exec"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type SelectStmtImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	clause.SelectImpl
	clause.FromImpl
	clause.JoinImpl
	clause.WhereImpl
	clause.GroupByImpl
	clause.HavingImpl
	clause.OrderByImpl
	clause.LimitImpl
	clause.OffsetImpl

	*exec.Executor
}

func NewSelectStatement(ctx context.Context, eng *ir.Engine, ref string) *SelectStmtImpl {
	stmt := &SelectStmtImpl{Dialect: eng.Dialect, DefaultRef: ref}
	exec := &exec.Executor{Ctx: ctx, Statement: stmt, Driver: eng.Driver}
	stmt.Executor = exec
	return stmt
}

func (s *SelectStmtImpl) Kind() dml.StatementKind {
	return dml.StatementSelect
}

func (s *SelectStmtImpl) Build() (*dml.Statement, error) {
	// Auto-FROM
	if !s.FromImpl.IsDeclared() && s.DefaultRef != "" {
		s.FromImpl.From((&exprimpl.TableImpl{}).Start(s.DefaultRef))
	}

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
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

	clauseMap := make(map[dml.ClauseKind]dml.ClauseBuilder)
	for _, cls := range clauses {
		if cls.IsDeclared() {
			clauseMap[cls.Kind()] = cls
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
		if cls, ok := clauseMap[ord]; ok {
			if !first {
				w.Char(' ')
			}
			first = false

			if err := cls.Build(w, s.Dialect); err != nil {
				return nil, err
			}
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
