package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/dml"
	"github.com/laacin/inyorm/internal/writer"
)

type UpdateQuery struct {
	Ref     string
	Dial    Dialect
	builder *core.Builder

	dml.QueryUpdate
}

// start

func (q *UpdateQuery) Start(dial Dialect, ref string) (*UpdateQuery, *builder.ExprBuilder) {
	b := builder.New()

	q.builder = b
	q.Dial = dial
	q.Ref = ref

	e := &builder.ExprBuilder{}
	return q, e.Start(b, ref)
}

// --- Build
func (q *UpdateQuery) Kind() QueryKind {
	return QueryUpdate
}

func (q *UpdateQuery) Build() (*QueryResult, error) {
	w := writer.New(q.Dial, false)

	q.QueryUpdate.Build(q.builder)
	q.QueryUpdate.Render(w, q.Dial)

	vals, err := q.builder.Param.Values()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}
