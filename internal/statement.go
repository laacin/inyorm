package internal

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/writer"
)

type Statement struct {
	dialect  entity.Dialect
	kind     entity.StatementKind
	clauses  map[entity.ClauseKind]entity.ClauseBuilder
	preBuild func(*writer.WriterImpl)
}

func NewStatement(kind entity.StatementKind, dial entity.Dialect) *Statement {
	return &Statement{
		dialect: dial,
		kind:    kind,
	}
}

func (s *Statement) LoadClauses(clauses []entity.ClauseBuilder) {
	s.clauses = make(map[entity.ClauseKind]entity.ClauseBuilder, len(clauses))

	for _, cls := range clauses {
		s.clauses[cls.Kind()] = cls
	}
}

func (s *Statement) PreBuild(fn func(*writer.WriterImpl)) {
	s.preBuild = fn
}

func (s *Statement) Build() Result {
	var ord []entity.ClauseKind
	switch s.kind {
	case entity.StatementSelect:
		ord = s.dialect.SelectOrder()
	case entity.StatementInsert:
		ord = s.dialect.InsertOrder()
	case entity.StatementUpdate:
		ord = s.dialect.UpdateOrder()
	case entity.StatementDelete:
		ord = s.dialect.DeleteOrder()
	}

	aliases := false
	if cls, ok := s.clauses[entity.ClauseJoin]; ok && cls.IsDeclared() {
		aliases = true
	}

	w := writer.New(s.dialect, aliases)

	s.preBuild(w)
	for _, knd := range ord {
		if cls, ok := s.clauses[knd]; ok && cls.IsDeclared() {
			cls.Build().Write(w, s.dialect)
			w.Char(' ')
		}
	}

	values, err := w.GetValues()
	return Result{
		Query:  w.Result(),
		Values: values,
		Err:    err,
	}
}

type Result struct {
	Query  string
	Values []any
	Err    error
}
