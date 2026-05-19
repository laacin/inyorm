package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type ConsCheckImpl struct{ emb ddl.ConsCheck }

func (c *ConsCheckImpl) Start(ident any) api.Condition {
	cond := &exprimpl.ConditionImpl{}
	c.emb.Cond = cond
	return cond.Start(ident)
}
