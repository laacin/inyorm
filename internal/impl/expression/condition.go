package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
)

type ConditionImpl struct {
	dml.Condition
	current dml.Predicate
}

func (c *ConditionImpl) Start(ident any) api.Condition {
	c.current = dml.Predicate{Identifier: ident}
	return c
}

func (c *ConditionImpl) Not() api.Condition {
	c.current.Negated = !c.current.Negated
	return c
}

func (c *ConditionImpl) Equal(value any) api.ConditionNext {
	c.push(dml.PredEqual, []any{value})
	return c
}

func (c *ConditionImpl) Like(value any) api.ConditionNext {
	c.push(dml.PredLike, []any{value})
	return c
}

func (c *ConditionImpl) In(values []any) api.ConditionNext {
	c.push(dml.PredIn, values)
	return c
}

func (c *ConditionImpl) Between(val1, val2 any) api.ConditionNext {
	c.push(dml.PredBetween, []any{val1, val2})
	return c
}

func (c *ConditionImpl) Greater(value any) api.ConditionNext {
	c.push(dml.PredGreater, []any{value})
	return c
}

func (c *ConditionImpl) Less(value any) api.ConditionNext {
	c.push(dml.PredLess, []any{value})
	return c
}

func (c *ConditionImpl) IsNull() api.ConditionNext {
	c.push(dml.PredIsNull, nil)
	return c
}

// --- Next
func (c *ConditionImpl) And(ident any) api.Condition {
	c.Connectors = append(c.Connectors, dml.PredAnd)
	return c.Start(ident)
}

func (c *ConditionImpl) Or(ident any) api.Condition {
	c.Connectors = append(c.Connectors, dml.PredOr)
	return c.Start(ident)
}

// --- Build
func (c *ConditionImpl) Build() dml.Value {
	return &c.Condition
}

// --- Helpers
func (c *ConditionImpl) push(op dml.PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.Predicates = append(c.Predicates, c.current)
}
