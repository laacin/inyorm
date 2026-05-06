package clause

import "github.com/laacin/inyorm/internal/entity"

type SelectImpl[Next any] struct {
	declared bool
	emb      entity.Select
}

func (c *SelectImpl[Next]) Distinct() Next {
	c.declared = true
	c.emb.Distinct = true
	return any(c).(Next)
}

func (c *SelectImpl[Next]) Select(values ...any) {
	c.declared = true
	c.emb.Values = values
}

// --- Build

func (c *SelectImpl[Next]) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *SelectImpl[Next]) Kind() entity.ClauseKind {
	return entity.ClauseSelect
}

func (c *SelectImpl[Next]) Build() entity.Clause {
	return &c.emb
}
