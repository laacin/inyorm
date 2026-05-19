package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type ConsForeignKeyImpl struct{ emb ddl.ConsForeignKey }

func (fk *ConsForeignKeyImpl) To(col, table string) api.ForeignKeyNext {
	fk.emb.ToTable = table
	fk.emb.ToColumn = col
	return fk
}

func (fk *ConsForeignKeyImpl) OnDel(key string) api.ForeignKeyNext {
	fk.emb.OnDelete = ddl.SetOnAct(key)
	return fk
}
func (fk *ConsForeignKeyImpl) OnUpd(key string) api.ForeignKeyNext {
	fk.emb.OnUpdate = ddl.SetOnAct(key)
	return fk
}
