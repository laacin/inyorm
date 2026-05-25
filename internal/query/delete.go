package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type DeleteQuery struct {
	Ref  string
	Dial Dialect

	dml.QueryDelete
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
	params := &writer.ParamStore{}
	w := &writer.WriterImpl{
		Syntax: q.Dial,
		Params: params,
	}

	q.QueryDelete.Render(w, q.Dial)
	vals, err := params.GetValues()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}
