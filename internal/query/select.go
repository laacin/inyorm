package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type SelectQuery struct {
	Ref  string
	Dial Dialect

	dml.QuerySelect
}

// start

func (q *SelectQuery) Start(dial Dialect, ref string) *SelectQuery {
	q.Dial = dial
	q.Ref = ref
	return q
}

// --- Build

func (*SelectQuery) Kind() QueryKind {
	return QuerySelect
}

func (q *SelectQuery) Build() (*QueryResult, error) {
	params := &writer.ParamStore{}
	w := &writer.WriterImpl{
		Syntax: q.Dial,
		Params: params,
	}

	if q.ClauseJoin.IsDeclared() {
		w.Aliases = &writer.AliasStore{}
	}

	q.QuerySelect.Render(w, q.Dial)
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: params.Values(),
	}, nil
}
