package clause

import "github.com/laacin/inyorm/internal/entity"

type SelectImpl[Next any] struct {
	declared bool
	emb      entity.Select
}

func (c *SelectImpl[Next]) Distinct() Next {
	c.declared = true
	c.emb.Dist = true
	return any(c).(Next)
}

func (c *SelectImpl[Next]) Select(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Defer

func (c *SelectImpl[Next]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *SelectImpl[Next]) Defer() entity.Clause {
	return &c.emb
}
