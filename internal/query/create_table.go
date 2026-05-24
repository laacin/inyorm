package query

import (
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/ddl"
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
	w := &writer.WriterImpl{Syntax: q.Dial}
	q.QueryCreateTable.Render(w, q.Dial)

	return &QueryResult{
		Query:  w.ToString(),
		Values: nil,
	}, nil
}
