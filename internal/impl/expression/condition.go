package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/expr"
)

type ConditionImpl struct {
	expr.Condition
	current expr.Predicate
}

func (c *ConditionImpl) Start(ident any) api.Condition {
	c.current = expr.Predicate{Identifier: ident}
	return c
}

func (c *ConditionImpl) Not() api.Condition {
	c.current.Negated = !c.current.Negated
	return c
}

func (c *ConditionImpl) Equal(value any) api.ConditionNext {
	c.push(expr.PredEqual, []any{value})
	return c
}

func (c *ConditionImpl) Like(value any) api.ConditionNext {
	c.push(expr.PredLike, []any{value})
	return c
}

func (c *ConditionImpl) In(values []any) api.ConditionNext {
	c.push(expr.PredIn, values)
	return c
}

func (c *ConditionImpl) Between(val1, val2 any) api.ConditionNext {
	c.push(expr.PredBetween, []any{val1, val2})
	return c
}

func (c *ConditionImpl) Greater(value any) api.ConditionNext {
	c.push(expr.PredGreater, []any{value})
	return c
}

func (c *ConditionImpl) Less(value any) api.ConditionNext {
	c.push(expr.PredLess, []any{value})
	return c
}

func (c *ConditionImpl) IsNull() api.ConditionNext {
	c.push(expr.PredIsNull, nil)
	return c
}

// --- Next
func (c *ConditionImpl) And(ident any) api.Condition {
	c.Connectors = append(c.Connectors, expr.PredAnd)
	return c.Start(ident)
}

func (c *ConditionImpl) Or(ident any) api.Condition {
	c.Connectors = append(c.Connectors, expr.PredOr)
	return c.Start(ident)
}

// --- Build
func (c *ConditionImpl) Build() expr.Value {
	return &c.Condition
}

// --- Helpers
func (c *ConditionImpl) push(op expr.PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.Predicates = append(c.Predicates, c.current)
}
