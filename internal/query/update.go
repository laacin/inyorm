package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type UpdateQuery struct {
	Ref  string
	Dial Dialect

	dml.UpdateQuery
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
	params := &writer.ParamStore{}
	w := &writer.WriterImpl{
		Syntax: q.Dial,
		Params: params,
	}

	q.UpdateQuery.Build(w, q.Dial)
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: params.Values(),
	}, nil
}
