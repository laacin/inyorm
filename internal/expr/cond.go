package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Cond struct {
	Predicates   []Predicate
	Connectors   []PredConnector
	current      Predicate
	identWrapper func(any) any
}

func NewCond(ident any, identWrapper ...func(any) any) *Cond {
	if len(identWrapper) > 0 {
		return &Cond{
			current:      Predicate{Identifier: identWrapper[0](ident)},
			identWrapper: identWrapper[0],
		}
	}
	return &Cond{current: Predicate{Identifier: ident}}
}

// --- PUB API

func (c *Cond) Not() api.Cond {
	c.current.Negated = !c.current.Negated
	return c
}

func (c *Cond) Equal(value any) api.CondNext {
	c.push(PredEqual, []any{value})
	return c
}

func (c *Cond) Like(value any) api.CondNext {
	c.push(PredLike, []any{value})
	return c
}

func (c *Cond) In(values ...any) api.CondNext {
	c.push(PredIn, values)
	return c
}

func (c *Cond) Between(val1, val2 any) api.CondNext {
	c.push(PredBetween, []any{val1, val2})
	return c
}

func (c *Cond) Greater(value any) api.CondNext {
	c.push(PredGreater, []any{value})
	return c
}

func (c *Cond) Less(value any) api.CondNext {
	c.push(PredLess, []any{value})
	return c
}

func (c *Cond) IsNull() api.CondNext {
	c.push(PredIsNull, nil)
	return c
}

func (c *Cond) And(ident any) api.Cond {
	c.Connectors = append(c.Connectors, PredAnd)
	return c.start(ident)
}

func (c *Cond) Or(ident any) api.Cond {
	c.Connectors = append(c.Connectors, PredOr)
	return c.start(ident)
}

// --- Build
func (*Cond) Kind() Kind { return KindCond }

func (c *Cond) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteExprCond(w, c, mode)
}

// --- Helpers
func (c *Cond) push(op PredOperator, values []any) {
	c.current.Values = values
	c.current.Operator = op
	c.current.Closed = true
	c.Predicates = append(c.Predicates, c.current)
}

func (c *Cond) start(ident any) api.Cond {
	if c.identWrapper != nil {
		c.current = Predicate{Identifier: c.identWrapper(ident)}
		return c
	}

	c.current = Predicate{Identifier: ident}
	return c
}
