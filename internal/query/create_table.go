package query

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/ddl"
)

type QueryCreateTable struct {
	Name   string
	Cols   []*ddl.ColDecl
	Fks    []*ddl.ForeignKey
	Checks []*ddl.Check
}

// --- PUB API
func (q *QueryCreateTable) TableName(name string) {
	q.Name = name
}

func (q *QueryCreateTable) String(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ddl.ColKindString)
}
func (q *QueryCreateTable) Int(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ddl.ColKindInt)
}
func (q *QueryCreateTable) Float(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ddl.ColKindFloat)
}
func (q *QueryCreateTable) Bool(name string) api.ColDeclNext {
	col := &ddl.ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ddl.ColKindBool)
}

func (q *QueryCreateTable) ForeignKey(on string) api.ForeignKey {
	fk := &ddl.ForeignKey{}
	q.Fks = append(q.Fks, fk)
	return fk.Start(on)
}
func (q *QueryCreateTable) Check(ident any) api.Cond {
	check := &ddl.Check{}
	q.Checks = append(q.Checks, check)
	return check.Start(ident)
}

// --- Render
func (q *QueryCreateTable) Build(b *core.Builder) error {
	b.Attach.MainRef = q.Name
	return nil
}

func (q *QueryCreateTable) Render(w core.InternalWriter, dial Dialect) error {
	dial.WriteQueryCreateTable(w, q)
	return nil
}
