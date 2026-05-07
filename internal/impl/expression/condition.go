package expression

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type ConditionImpl struct {
	entity.Condition
	current entity.Predicate
}

func (c *ConditionImpl) Start(ident any) api.Condition {
	c.current = entity.Predicate{Identifier: ident}
	return c
}

func (c *ConditionImpl) Not() api.Condition {
	c.current.Negated = !c.current.Negated
	return c
}

func (c *ConditionImpl) Equal(value any) api.ConditionNext {
	c.push(entity.PredEqual, []any{value})
	return c
}

func (c *ConditionImpl) Like(value any) api.ConditionNext {
	c.push(entity.PredLike, []any{value})
	return c
}

func (c *ConditionImpl) In(values []any) api.ConditionNext {
	c.push(entity.PredIn, values)
	return c
}

func (c *ConditionImpl) Between(val1, val2 any) api.ConditionNext {
	c.push(entity.PredBetween, []any{val1, val2})
	return c
}

func (c *ConditionImpl) Greater(value any) api.ConditionNext {
	c.push(entity.PredGreater, []any{value})
	return c
}

func (c *ConditionImpl) Less(value any) api.ConditionNext {
	c.push(entity.PredLess, []any{value})
	return c
}

func (c *ConditionImpl) IsNull() api.ConditionNext {
	c.push(entity.PredIsNull, nil)
	return c
}

// --- Next
func (c *ConditionImpl) And(ident any) api.Condition {
	c.Connectors = append(c.Connectors, entity.PredAnd)
	return c
}

func (c *ConditionImpl) Or(ident any) api.Condition {
	c.Connectors = append(c.Connectors, entity.PredOr)
	return c
}

// --- Build
func (c *ConditionImpl) Build() entity.Value {
	return &c.Condition
}

// --- Helpers
func (c *ConditionImpl) push(op entity.PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.Predicates = append(c.Predicates, c.current)
}
