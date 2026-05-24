package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

type QueryCreateTable struct {
	Name   string
	Cols   []*ColDecl
	Fks    []*ForeignKey
	Checks []*Check
}

func (q *QueryCreateTable) Start(tableName string) {
	q.Name = tableName
}

// --- PUB API
func (q *QueryCreateTable) String(name string) api.ColDeclNext {
	col := &ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ColKindString)
}
func (q *QueryCreateTable) Int(name string) api.ColDeclNext {
	col := &ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ColKindInt)
}
func (q *QueryCreateTable) Float(name string) api.ColDeclNext {
	col := &ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ColKindFloat)
}
func (q *QueryCreateTable) Bool(name string) api.ColDeclNext {
	col := &ColDecl{}
	q.Cols = append(q.Cols, col)
	return col.Start(name, ColKindBool)
}

func (q *QueryCreateTable) ForeignKey(on string) api.ForeignKey {
	fk := &ForeignKey{}
	q.Fks = append(q.Fks, fk)
	return fk.Start(on)
}
func (q *QueryCreateTable) Check(ident any) api.Cond {
	check := &Check{}
	q.Checks = append(q.Checks, check)
	return check.Start(ident)
}

// --- Render

func (b *QueryCreateTable) Render(w core.InternalWriter, dial TableWriter) {
	dial.WriteCreateTable(w, b)
}
