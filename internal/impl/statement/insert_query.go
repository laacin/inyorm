package statement

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type InsertQueryImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	clause.InsertIntoImpl
}

func (s *InsertQueryImpl) Start(dial ir.Dialect, ref string) api.InsertQuery {
	s.Dialect = dial
	s.DefaultRef = ref
	return s
}

func (s *InsertQueryImpl) Kind() dml.QueryKind {
	return dml.QueryInsert
}

func (s *InsertQueryImpl) Build() (string, []any, error) {
	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&s.InsertIntoImpl,
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
	for _, ord := range s.Dialect.InsertOrder() {
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
