package statement

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type UpdateQueryImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	clause.UpdateImpl
	clause.WhereImpl
}

func (s *UpdateQueryImpl) Start(dial ir.Dialect, ref string) api.UpdateQuery {
	s.Dialect = dial
	s.DefaultRef = ref
	return s
}

func (s *UpdateQueryImpl) Kind() dml.QueryKind {
	return dml.QueryUpdate
}

func (s *UpdateQueryImpl) Build() (string, []any, error) {
	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.UpdateImpl,
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
	for _, ord := range s.Dialect.UpdateOrder() {
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
