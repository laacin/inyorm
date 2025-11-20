package column

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value any] struct {
	Table string
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Col(name string, table ...string) Col {
	tbl := c.Table
	if len(table) > 0 {
		tbl = table[0]
	}

	col := &Column[Col, Value]{BaseName: name, Table: tbl}
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) All() Col {
	col := &Column[Col, Value]{BaseName: "*"}
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Ph() core.Builder {
	return func(w core.Writer) { w.Placeholder() }
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Cond(ident Ident) Cond {
	cond := &condition.Condition[Cond, CondNext, Ident, Value]{}
	return cond.Start(ident)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Concat(v ...Value) Col {
	expr := func(w core.Writer) {
		w.Write("CONCAT(")
		for i, val := range v {
			if i > 0 {
				w.Write(", ")
			}
			inferColumn[Col, Value](w, val)
		}
		w.Char(')')
	}

	col := &Column[Col, Value]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Switch(cond Ident, fn func(cs CaseSwitch)) Col {
	cs := &Case[CaseSwitch, CaseNext, Ident, Value]{}
	fn(any(cs).(CaseSwitch))

	expr := func(w core.Writer) {
		w.Write("CASE ")
		inferColumn[Col, Value](w, cond)
		for _, expr := range cs.Exprs {
			w.Write(" WHEN ")
			inferColumn[Col, Value](w, expr.Identifier)
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

	col := &Column[Col, Value]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}

func (c *ColBuilder[Col, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext, Ident, Value]) Search(fn func(cs CaseSearch)) Col {
	cs := &Case[CaseSearch, CaseNext, CondNext, Value]{}
	fn(any(cs).(CaseSearch))

	expr := func(w core.Writer) {
		w.Write("CASE WHEN")
		for _, arg := range cs.Exprs {
			w.Char(' ')
			any(arg.Identifier).(*condition.Condition[Cond, CondNext, Ident, Value]).Build(w, core.ClsTypUnset)
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

	col := &Column[Col, Value]{}
	col.Builder.WExpr(expr)
	return any(col).(Col)
}
