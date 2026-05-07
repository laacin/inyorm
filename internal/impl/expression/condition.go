package expression

import "github.com/laacin/inyorm/internal/entity"

type ConditionImpl[Self, Next any] struct {
	entity.Condition
	current entity.Predicate
}

func (c *ConditionImpl[Self, Next]) Start(ident any) Self {
	c.current = entity.Predicate{Identifier: ident}
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Not() Self {
	c.current.Negated = !c.current.Negated
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Equal(value any) Next {
	c.push(entity.PredEqual, []any{value})
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Like(value any) Next {
	c.push(entity.PredLike, []any{value})
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) In(values []any) Next {
	c.push(entity.PredIn, values)
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Between(val1, val2 any) Next {
	c.push(entity.PredBetween, []any{val1, val2})
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Greater(value any) Next {
	c.push(entity.PredGreater, []any{value})
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Less(value any) Next {
	c.push(entity.PredLess, []any{value})
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) IsNull() Next {
	c.push(entity.PredIsNull, nil)
	return any(c).(Next)
}

// --- Next
func (c *ConditionImpl[Self, Next]) And(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredAnd)
	return c.Start(ident)
}

func (c *ConditionImpl[Self, Next]) Or(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredOr)
	return c.Start(ident)
}

// --- Build
func (c *ConditionImpl[Self, Next]) Build() entity.Value {
	return &c.Condition
}

// --- Helpers
func (c *ConditionImpl[Self, Next]) push(op entity.PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.Predicates = append(c.Predicates, c.current)
}
