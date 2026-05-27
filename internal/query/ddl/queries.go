package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

type QueryCreateTable struct {
	rend Renderer

	Name   string
	Cols   []*ColDecl
	Fks    []*ForeignKey
	Checks []*Check
}

func NewCreateTable(rend Renderer) *QueryCreateTable {
	return &QueryCreateTable{rend: rend}
}

// --- PUB API
func (q *QueryCreateTable) TableName(name string) {
	q.Name = name
}

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
func (q *QueryCreateTable) Build(tools *query.Tools) error {
	tools.Aliases.SetMain(q.Name)
	return nil
}

func (q *QueryCreateTable) Render(w core.InternalWriter) error {
	q.rend.WriteQueryCreateTable(w, q)
	return nil
}
