package statement

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/statement/writer"
)

type InsertStmtImpl struct {
	DefaultRef string
	Dialect    entity.Dialect

	clause.InsertIntoImpl
}

func (s *InsertStmtImpl) Kind() entity.StatementKind {
	return entity.StatementInsert
}

func (s *InsertStmtImpl) Build() *entity.Query {
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
		errs       []error
	)

	// --- Write the statement

	w := &writer.WriterImpl{
		Syntax: s.Dialect,
		Params: parameters,
	}

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
		errs = append(errs, err)
	}

	return &entity.Query{
		Statement: w.ToString(),
		Values:    parameters.Values(),
		Errs:      errs,
	}
}
