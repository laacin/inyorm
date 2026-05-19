package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type ConsDeclImpl struct{ emb ddl.ConsDecl[any] }

func (c *ConsDeclImpl) Start(defaultTable string) api.ConsDecl {
	c.emb.Table = defaultTable
	return c
}

func (c *ConsDeclImpl) Index(col string) {
	c.emb.Kind = ddl.ConsKindIndex
	c.emb.Value = &ddl.ConsIndex{}
	c.emb.Column = col
}

func (c *ConsDeclImpl) ForeignKey(col string) api.ForeignKey {
	c.emb.Kind = ddl.ConsKindForeignKey
	c.emb.Column = col
	fk := &ConsForeignKeyImpl{}
	c.emb.Value = &fk.emb
	return fk
}

func (c *ConsDeclImpl) Check(ident any) api.Condition {
	c.emb.Kind = ddl.ConsKindCheck
	check := &ConsCheckImpl{}
	c.emb.Value = &check.emb
	return check.Start(ident)
}

func (c *ConsDeclImpl) Default(col string) api.Default {
	c.emb.Kind = ddl.ConsKindDefault
	c.emb.Column = col
	dflt := &ConsDefaultImpl{}
	c.emb.Value = &dflt.emb
	return dflt
}
