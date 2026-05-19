package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/exec"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type DeleteStmtImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	clause.DeleteImpl
	clause.FromImpl
	clause.WhereImpl

	*exec.Executor
}

func (s *DeleteStmtImpl) Start(ctx context.Context, eng *ir.Engine, ref string) api.DeleteStmt {
	s.DefaultRef = ref
	s.Dialect = eng.Dialect
	s.Executor = &exec.Executor{Ctx: ctx, Statement: s, Driver: eng.Driver}
	return s
}

func (s *DeleteStmtImpl) Kind() dml.StatementKind {
	return dml.StatementDelete
}

func (s *DeleteStmtImpl) Build() (*dml.Statement, error) {
	// Auto-FROM
	if !s.FromImpl.IsDeclared() && s.DefaultRef != "" {
		s.FromImpl.From((&exprimpl.TableImpl{}).Start(s.DefaultRef))
	}

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.DeleteImpl,
		&s.FromImpl,
		&s.WhereImpl,
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
	)

	w := &writer.WriterImpl{
		Syntax: s.Dialect,
		Params: parameters,
	}

	// --- Write the statement

	first := true
	for _, ord := range s.Dialect.DeleteOrder() {
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
