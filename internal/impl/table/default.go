package table

import "github.com/laacin/inyorm/internal/ir/ddl"

type ConsDefaultImpl struct{ emb ddl.ConsDefault }

func (c *ConsDefaultImpl) Value(value any) {
	c.emb.Value = value
}
