package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type TableBuilderImpl struct {
	name string
	cols []*ColDeclImpl
	cons []*ConsDeclImpl
}

func (t *TableBuilderImpl) Start(name string) api.TableBuilder {
	t.name = name
	return t
}

func (t *TableBuilderImpl) Text(name string) api.ColDecl {
	c := &ColDeclImpl{}
	t.cols = append(t.cols, c)
	return c.Start(name, ddl.ColKindText)
}
func (t *TableBuilderImpl) Int(name string) api.ColDecl {
	c := &ColDeclImpl{}
	t.cols = append(t.cols, c)
	return c.Start(name, ddl.ColKindInt)
}
func (t *TableBuilderImpl) Float(name string) api.ColDecl {
	c := &ColDeclImpl{}
	t.cols = append(t.cols, c)
	return c.Start(name, ddl.ColKindFloat)
}
func (t *TableBuilderImpl) Bool(name string) api.ColDecl {
	c := &ColDeclImpl{}
	t.cols = append(t.cols, c)
	return c.Start(name, ddl.ColKindBool)
}

func (t *TableBuilderImpl) ForeignKey(col string) api.ForeignKey {
	c := &ConsDeclImpl{}
	t.cons = append(t.cons, c)
	return c.ForeignKey(col, t.name)
}
func (t *TableBuilderImpl) Check(ident any) api.Condition {
	c := &ConsDeclImpl{}
	t.cons = append(t.cons, c)
	return c.Check(ident)
}

// --- Build

func (t *TableBuilderImpl) Build(w core.InternalWriter, dial ddl.TableWriter) {
	cols := make([]ddl.ColDecl, len(t.cols))
	cons := make([]ddl.ConsDecl[any], len(t.cons))

	for i, c := range t.cols {
		cols[i] = c.emb
	}
	for i, c := range t.cons {
		cons[i] = c.emb
	}

	dial.WriteCreateTable(w, &ddl.TableDecl{
		Name: t.name,
		Cols: cols,
		Cons: cons,
	})
}
