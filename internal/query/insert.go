package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type InsertQuery struct {
	Ref  string
	Dial Dialect

	dml.ClauseInsert
}

// start

func (q *InsertQuery) Start(dial Dialect, ref string) *InsertQuery {
	q.Dial = dial
	q.Ref = ref
	return q
}

// --- Build

func (*InsertQuery) Kind() QueryKind {
	return QueryInsert
}

func (q *InsertQuery) Build() (*QueryResult, error) {
	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&q.ClauseInsert,
	}

	clauseMap := make(map[dml.ClauseKind]dml.ClauseBuilder)
	for _, cls := range clauses {
		if cls.IsDeclared() {
			clauseMap[cls.Kind()] = cls
		}
	}

	// --- Declarate writers

	var (
		params = &writer.ParamStore{}
	)

	w := &writer.WriterImpl{
		Syntax: q.Dial,
		Params: params,
	}

	// --- Write the statement

	first := true
	for _, ord := range q.Dial.InsertOrder() {
		if cls, ok := clauseMap[ord]; ok {
			if !first {
				w.Char(' ')
			}
			first = false

			if err := cls.Build(w, q.Dial); err != nil {
				return nil, err
			}
		}
	}

	// --- Validate values

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: params.Values(),
	}, nil
}
