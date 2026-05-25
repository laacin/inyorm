package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/dml"
	"github.com/laacin/inyorm/internal/writer"
)

type InsertQuery struct {
	Ref     string
	Dial    Dialect
	builder *core.Builder

	dml.QueryInsert
}

// start

func (q *InsertQuery) Start(dial Dialect, ref string) (*InsertQuery, *builder.ExprBuilder) {
	b := builder.New()

	q.builder = b
	q.Dial = dial
	q.Ref = ref

	e := &builder.ExprBuilder{}
	return q, e.Start(b, ref)
}

// --- Build

func (*InsertQuery) Kind() QueryKind {
	return QueryInsert
}

func (q *InsertQuery) Build() (*QueryResult, error) {
	w := writer.New(q.Dial, false)

	q.QueryInsert.Build(q.builder)
	q.QueryInsert.Render(w, q.Dial)

	vals, err := q.builder.Param.Values()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}
