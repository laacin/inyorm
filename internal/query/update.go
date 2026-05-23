package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type UpdateQuery struct {
	Ref  string
	Dial Dialect

	dml.ClauseUpdate
	dml.ClauseWhere
}

// start

func (q *UpdateQuery) Start(dial Dialect, ref string) *UpdateQuery {
	q.Dial = dial
	q.Ref = ref
	return q
}

// --- Build
func (q *UpdateQuery) Kind() QueryKind {
	return QueryUpdate
}

func (q *UpdateQuery) Build() (*QueryResult, error) {
	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&q.ClauseUpdate,
		&q.ClauseWhere,
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
	for _, ord := range q.Dial.UpdateOrder() {
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
