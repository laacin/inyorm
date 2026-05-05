package expression

import "github.com/laacin/inyorm/internal/entity"

type ConditionImpl[Self, Next any] struct {
	entity.Condition
	Current *entity.Predicate
}

func (c *ConditionImpl[Self, Next]) Start(ident any) Self {
	c.Current = &entity.Predicate{Identifier: ident}
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Not() Self {
	c.Current.Negated = !c.Current.Negated
	return any(c).(Self)
}

func (c *ConditionImpl[Self, Next]) Equal(value any) Next {
	c.Current.Values = []any{value}
	c.Current.Operator = entity.PredEqual
	c.Current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) In(values []any) Next {
	c.Current.Values = values
	c.Current.Operator = entity.PredIn
	c.Current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Between(val1, val2 any) Next {
	c.Current.Values = []any{val1, val2}
	c.Current.Operator = entity.PredBetween
	c.Current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Greater(value any) Next {
	c.Current.Values = []any{value}
	c.Current.Operator = entity.PredGreater
	c.Current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) Less(value any) Next {
	c.Current.Values = []any{value}
	c.Current.Operator = entity.PredLess
	c.Current.Closed = true
	return any(c).(Next)
}

func (c *ConditionImpl[Self, Next]) IsNull() Next {
	c.Current.Values = nil
	c.Current.Operator = entity.PredIsNull
	c.Current.Closed = true
	return any(c).(Next)
}

// --- Next
func (c *ConditionImpl[Self, Next]) Any(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredAnd)
	c.Predicates = append(c.Predicates, *c.Current)
	return c.Start(ident)
}

func (c *ConditionImpl[Self, Next]) Or(ident any) Self {
	c.Connectors = append(c.Connectors, entity.PredOr)
	c.Predicates = append(c.Predicates, *c.Current)
	return c.Start(ident)
}

// --- Deref
func (c *ConditionImpl[Self, Next]) Deref() entity.Value {
	return &c.Condition
}
