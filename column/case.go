package column

import (
	"github.com/laacin/inyorm/internal/stmt"
	"strings"
)

// ---- CASE
func (c *Column) Switch(target any, fn func(iface.Case[string])) {
	cs := &Case[string]{
		target: stmt.Stringify(target),
	}

	fn(cs)
	buildSwitch(c, cs)
}

func (c *Column) Search() iface.Case[func(iface.Expression)] {
	return &Case[func(iface.Expression)]{}
}

type Case[T any] struct {
	target string
	thens  []*CaseThen[T]
	els    string
}

func (c *Case[T]) When(value T) iface.Then[T] {
	thn := &CaseThen[T]{
		ctx:       c,
		condition: value,
	}
	c.thens = append(c.thens, thn)
	return thn
}

func (s *Case[T]) Else(do any) {
	s.els = stmt.Stringify(do)
}

type CaseThen[T any] struct {
	ctx       *Case[T]
	condition T
	do        string
}

func (c *CaseThen[T]) Then(do any) iface.Case[T] {
	c.do = stmt.Stringify(do)
	return c.ctx
}

// -- Build
func buildSwitch(c *Column, cs *Case[string]) {
	var sb strings.Builder

	sb.WriteString("CASE ")
	sb.WriteString(cs.target)
	for _, thn := range cs.thens {
		sb.WriteString(" WHEN ")
		sb.WriteString(stmt.Stringify(thn.condition))
		sb.WriteString(" THEN ")
		sb.WriteString(thn.do)
	}

	if cs.els != "" {
		sb.WriteString(" ELSE ")
		sb.WriteString(cs.els)
	}
	sb.WriteString(" END")
	*c = Column(sb.String())
}
