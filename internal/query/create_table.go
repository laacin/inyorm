package query

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/ddl"
)

type CreateTableQuery struct {
	Dial Dialect

	table ddl.CreateTable
}

// start

func (q *CreateTableQuery) Start(dial Dialect, ref string) *CreateTableQuery {
	q.table.Name = ref
	q.Dial = dial
	return q
}

// --- PUB API

func (q *CreateTableQuery) String(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.table.Cols = append(q.table.Cols, col)
	return col.Start(name, ddl.ColKindString)
}
func (q *CreateTableQuery) Int(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.table.Cols = append(q.table.Cols, col)
	return col.Start(name, ddl.ColKindInt)
}
func (q *CreateTableQuery) Float(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.table.Cols = append(q.table.Cols, col)
	return col.Start(name, ddl.ColKindFloat)
}
func (q *CreateTableQuery) Bool(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.table.Cols = append(q.table.Cols, col)
	return col.Start(name, ddl.ColKindBool)
}

func (q *CreateTableQuery) ForeignKey(on string) api.ForeignKey {
	fk := &ddl.ForeignKey{}
	q.table.Fks = append(q.table.Fks, fk)
	return fk.Start(on)
}
func (q *CreateTableQuery) Check(ident any) api.Cond {
	check := &ddl.Check{}
	q.table.Checks = append(q.table.Checks, check)
	return check.Start(ident)
}

// --- Build

func (q *CreateTableQuery) Kind() QueryKind {
	return QueryCreateTable
}

func (q *CreateTableQuery) Build() (*QueryResult, error) {
	w := &writer.WriterImpl{Syntax: q.Dial}
	q.table.Build(w, q.Dial)

	return &QueryResult{
		Query:  w.ToString(),
		Values: nil,
	}, nil
}
