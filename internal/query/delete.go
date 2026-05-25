package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/dml"
	"github.com/laacin/inyorm/internal/writer"
)

type DeleteQuery struct {
	Ref     string
	Dial    Dialect
	builder *core.Builder

	dml.QueryDelete
}

// start

func (q *DeleteQuery) Start(dial Dialect, ref string) (*DeleteQuery, *builder.ExprBuilder) {
	b := builder.New()

	q.builder = b
	q.Dial = dial
	q.Ref = ref

	e := &builder.ExprBuilder{}
	return q, e.Start(b, ref)
}

// --- Build

func (q *DeleteQuery) Kind() QueryKind {
	return QueryDelete
}

func (q *DeleteQuery) Build() (*QueryResult, error) {
	w := writer.New(q.Dial, false)

	q.QueryDelete.Build(q.builder)
	q.QueryDelete.Render(w, q.Dial)

	vals, err := q.builder.Param.Values()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}
