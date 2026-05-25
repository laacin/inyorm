package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/dml"
	"github.com/laacin/inyorm/internal/writer"
)

type SelectQuery struct {
	Ref     string
	Dial    Dialect
	builder *core.Builder

	dml.QuerySelect
}

// start

func (q *SelectQuery) Start(dial Dialect, ref string) (*SelectQuery, *builder.ExprBuilder) {
	b := builder.New()

	q.builder = b
	q.Dial = dial
	q.Ref = ref

	e := &builder.ExprBuilder{}
	return q, e.Start(b, ref)
}

// --- Build

func (*SelectQuery) Kind() QueryKind {
	return QuerySelect
}

func (q *SelectQuery) Build() (*QueryResult, error) {
	w := writer.New(q.Dial, q.ClauseJoin.IsDeclared())

	q.QuerySelect.Build(q.builder)
	q.QuerySelect.Render(w, q.Dial)

	vals, err := q.builder.Param.Values()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}
