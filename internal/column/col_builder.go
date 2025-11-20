package column

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type ColBuilder[Col, Cond, CondNext, Case, CaseNext any] struct {
	Table string
}

func (c *ColBuilder[Col, Cond, CondNext, Case, CaseNext]) Col(name string, table ...string) Col {
	tbl := c.Table
	if len(table) > 0 {
		tbl = table[0]
	}

	col := &Column[Col]{BaseName: name, Table: tbl}
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, Case, CaseNext]) All() Col {
	col := &Column[Col]{BaseName: "*"}
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, Case, CaseNext]) Ph() core.Builder {
	return func(w core.Writer) { w.Placeholder() }
}

func (c *ColBuilder[Col, Cond, CondNext, Case, CaseNext]) Cond(ident any) Cond {
	cond := &condition.Condition[Cond, CondNext]{}
	return cond.Start(ident)
}

func (c *ColBuilder[Col, Cond, CondNext, Case, CaseNext]) Concat(v ...any) Col {
	expr := func(w core.Writer) {
		w.Write("CONCAT(")
		for i, val := range v {
			if i > 0 {
				w.Write(", ")
			}
			inferColumn[Col](w, val)
		}
		w.Char(')')
	}

	col := &Column[Col]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseT, CaseNext]) Switch(cond any, fn func(cs CaseT)) Col {
	cs := &Case[CaseT, CaseNext]{}
	fn(any(cs).(CaseT))

	expr := func(w core.Writer) {
		w.Write("CASE ")
		inferColumn[Col](w, cond)
		for _, expr := range cs.Exprs {
			w.Write(" WHEN ")
			inferColumn[Col](w, expr.Identifier)
			w.Write(" THEN ")
			w.Value(expr.Argument, core.ClsTypUnset)
			w.Char(' ')
		}
		if any(cs.Els) != nil {
			w.Write("ELSE ")
			w.Value(cs.Els, core.ClsTypUnset)
			w.Char(' ')
		}
		w.Write("END")
	}

	col := &Column[Col]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseT, CaseNext]) Search(fn func(cs CaseT)) Col {
	cs := &Case[CaseT, CaseNext]{}
	fn(any(cs).(CaseT))

	expr := func(w core.Writer) {
		w.Write("CASE WHEN")
		for _, arg := range cs.Exprs {
			w.Char(' ')
			any(arg.Identifier).(*condition.Condition[Cond, CondNext]).Build(w, core.ClsTypUnset)
			w.Write(" THEN ")
			w.Value(arg.Argument, core.ClsTypUnset)
			w.Char(' ')
		}
		if any(cs.Els) != nil {
			w.Write("ELSE ")
			w.Value(cs.Els, core.ClsTypUnset)
			w.Char(' ')
		}
		w.Write("END")
	}

	col := &Column[Col]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}
