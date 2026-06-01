package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
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
	col := NewColDecl(name, ColKindString)
	q.Cols = append(q.Cols, col)
	return col
}
func (q *QueryCreateTable) Int(name string) api.ColDeclNext {
	col := NewColDecl(name, ColKindInt)
	q.Cols = append(q.Cols, col)
	return col
}
func (q *QueryCreateTable) Float(name string) api.ColDeclNext {
	col := NewColDecl(name, ColKindFloat)
	q.Cols = append(q.Cols, col)
	return col
}
func (q *QueryCreateTable) Bool(name string) api.ColDeclNext {
	col := NewColDecl(name, ColKindBool)
	q.Cols = append(q.Cols, col)
	return col
}

func (q *QueryCreateTable) ForeignKey(on string) api.ForeignKey {
	fk := NewForeignKey(on)
	q.Fks = append(q.Fks, fk)
	return fk
}
func (q *QueryCreateTable) Check(ident any) api.Cond {
	cond := expr.NewCond(ident, func(a any) any {
		if str, ok := a.(string); ok {
			return expr.NewCol(str, nil)
		}
		return a
	})
	q.Checks = append(q.Checks, NewCheck(cond))
	return cond
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
