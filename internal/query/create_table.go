package query

import (
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/writer"
)

type CreateTable struct {
	Dial Dialect

	ddl.QueryCreateTable
}

// start

func (q *CreateTable) Start(dial Dialect, ref string) *CreateTable {
	q.QueryCreateTable.Start(ref)
	q.Dial = dial
	return q
}

// --- Build

func (q *CreateTable) Kind() QueryKind {
	return QueryCreateTable
}

func (q *CreateTable) Build() (*QueryResult, error) {
	w := writer.New(q.Dial, false)

	q.QueryCreateTable.Render(w, q.Dial)

	return &QueryResult{
		Query:  w.ToString(),
		Values: nil,
	}, nil
}
