package expression

import "github.com/laacin/inyorm/internal/entity"

type ConditionImpl[Self, Next any] struct {
	entity.Condition
	current *entity.Predicate
}

func (c *ConditionImpl[Self, Next]) Start(ident any) Self {
	c.current = &entity.Predicate{Identifier: ident}
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Not() Self {
	c.current.Negated = !c.current.Negated
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Equal(value any) Next {
	c.current.Values = []any{value}
	c.current.Operator = entity.PredEqual
	c.current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) In(values []any) Next {
	c.current.Values = values
	c.current.Operator = entity.PredIn
	c.current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Between(val1, val2 any) Next {
	c.current.Values = []any{val1, val2}
	c.current.Operator = entity.PredBetween
	c.current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Greater(value any) Next {
	c.current.Values = []any{value}
	c.current.Operator = entity.PredGreater
	c.current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Less(value any) Next {
	c.current.Values = []any{value}
	c.current.Operator = entity.PredLess
	c.current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) IsNull() Next {
	c.current.Values = nil
	c.current.Operator = entity.PredIsNull
	c.current.Closed = true
	return any(c).(Next)
}

// --- Next
func (c *ConditionImpl[Self, Next]) Any(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredAnd)
	c.Predicates = append(c.Predicates, *c.current)
	return c.Start(ident)
}

func (c *ConditionImpl[Self, Next]) Or(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredOr)
	c.Predicates = append(c.Predicates, *c.current)
	return c.Start(ident)
}

// --- Deref
func (c *ConditionImpl[Self, Next]) Deref() entity.Value {
	return &c.Condition
}
