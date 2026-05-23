package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type DeleteQuery struct {
	Ref  string
	Dial Dialect

	dml.DeleteBuilder
	dml.FromBuilder
	dml.WhereBuilder
}

// start

func (q *DeleteQuery) Start(dial Dialect, ref string) *DeleteQuery {
	q.Dial = dial
	q.Ref = ref
	return q
}

// --- Build

func (q *DeleteQuery) Kind() QueryKind {
	return QueryDelete
}

func (q *DeleteQuery) Build() (*QueryResult, error) {
	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&q.DeleteBuilder,
		&q.FromBuilder,
		&q.WhereBuilder,
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
	for _, ord := range q.Dial.DeleteOrder() {
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
