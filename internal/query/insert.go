package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type InsertQuery struct {
	Ref  string
	Dial Dialect

	dml.QueryInsert
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
	params := &writer.ParamStore{}

	w := &writer.WriterImpl{
		Syntax: q.Dial,
		Params: params,
	}

	q.QueryInsert.Render(w, q.Dial)
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: params.Values(),
	}, nil
}
