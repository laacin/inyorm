package statement

import (
	"errors"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type SelectQueryImpl struct {
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
}

func (s *SelectQueryImpl) Start(dial ir.Dialect, ref string) api.SelectQuery {
	s.Dialect = dial
	s.DefaultRef = ref
	return s
}

func (s *SelectQueryImpl) Kind() dml.QueryKind {
	return dml.QuerySelect
}

func (s *SelectQueryImpl) Build() (string, []any, error) {
	// --- Guards

	if !s.SelectImpl.IsDeclared() {
		return "", nil, errors.New("clause 'SELECT' must be declared")
	}
	if !s.FromImpl.IsDeclared() {
		return "", nil, errors.New("clause 'FROM' must be declared")
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
				return "", nil, err
			}
		}
	}

	// --- Validate values

	if err := parameters.Validate(); err != nil {
		return "", nil, err
	}

	return w.ToString(), parameters.Values(), nil
}
