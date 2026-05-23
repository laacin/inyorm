package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Cond struct {
	Predicates []Predicate
	Connectors []PredConnector
}

// --- Builder
type CondBuilder struct {
	emb     Cond
	current Predicate
}

// start
func (c *CondBuilder) Start(ident any) *CondBuilder {
	c.current = Predicate{Identifier: ident}
	return c
}

// --- PUB API

func (c *CondBuilder) Not() api.Cond {
	c.current.Negated = !c.current.Negated
	return c
}

func (c *CondBuilder) Equal(value any) api.CondNext {
	c.push(PredEqual, []any{value})
	return c
}

func (c *CondBuilder) Like(value any) api.CondNext {
	c.push(PredLike, []any{value})
	return c
}

func (c *CondBuilder) In(values []any) api.CondNext {
	c.push(PredIn, values)
	return c
}

func (c *CondBuilder) Between(val1, val2 any) api.CondNext {
	c.push(PredBetween, []any{val1, val2})
	return c
}

func (c *CondBuilder) Greater(value any) api.CondNext {
	c.push(PredGreater, []any{value})
	return c
}

func (c *CondBuilder) Less(value any) api.CondNext {
	c.push(PredLess, []any{value})
	return c
}

func (c *CondBuilder) IsNull() api.CondNext {
	c.push(PredIsNull, nil)
	return c
}

func (c *CondBuilder) And(ident any) api.Cond {
	c.emb.Connectors = append(c.emb.Connectors, PredAnd)
	return c.Start(ident)
}

func (c *CondBuilder) Or(ident any) api.Cond {
	c.emb.Connectors = append(c.emb.Connectors, PredOr)
	return c.Start(ident)
}

// --- Build
func (c *CondBuilder) Kind() ExprKind {
	return ExprCond
}

func (c *CondBuilder) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteCond(w, &c.emb, mode)
}

// --- Helpers
func (c *CondBuilder) push(op PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.emb.Predicates = append(c.emb.Predicates, c.current)
}
