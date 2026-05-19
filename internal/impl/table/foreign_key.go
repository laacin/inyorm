package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type ConsForeignKeyImpl struct{ emb ddl.ConsForeignKey }

func (c *ConsForeignKeyImpl) To(col, table string) api.ForeignKeyNext {
	c.emb.ToTable = table
	c.emb.ToColumn = col
	return c
}

func (c *ConsForeignKeyImpl) OnDel(key string) api.ForeignKeyNext {
	c.emb.OnDelete = ddl.SetOnAct(key)
	return c
}
func (c *ConsForeignKeyImpl) OnUpd(key string) api.ForeignKeyNext {
	c.emb.OnUpdate = ddl.SetOnAct(key)
	return c
}
