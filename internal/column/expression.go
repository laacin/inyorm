package column

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type ColExpr struct {
	Statement *writer.Statement
}

func (c *ColExpr) Col(v any, table ...string) core.Column {
	tbl := c.Statement.GetFrom()
	if len(table) > 0 {
		tbl = table[0]
	}

	switch val := v.(type) {
	case string:
		return &Column{
			Type:   normalCol,
			Writer: c.Statement.Writer(),
			Alias:  tbl,
			Value:  val,
		}

	case core.Column:
		return val

	default:
		return &Column{
			Type:   normalCol,
			Writer: c.Statement.Writer(),
			Alias:  tbl,
			Value:  core.Sqlify(v),
		}
	}
}

func (c *ColExpr) All() core.Column {
	return &Column{
		Writer: c.Statement.Writer(),
		Type:   customCol,
		Value:  "*",
	}
}

func (c *ColExpr) Concat(v ...any) core.Column {
	w := c.Statement.Writer()

	w.Write("CONCAT(")
	for i, val := range v {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, &core.ValueOpts{Definition: true})
	}
	w.Char(')')

	return &Column{
		Writer: w,
		Type:   customCol,
		Value:  w.ToString(),
	}
}

func (c *ColExpr) Switch(cond any, fn func(cs core.CaseSwitch)) core.Column {
	w := c.Statement.Writer()
	opts := &core.ValueOpts{Definition: true}

	cs := &Case[any]{}
	fn(cs)

	w.Write("CASE ")
	w.Value(cond, opts)
	for _, arg := range cs.args {
		w.Write(" WHEN ")
		w.Value(arg.when, opts)
		w.Write(" THEN ")
		w.Value(arg.do, opts)
		w.Char(' ')
	}
	w.Write("ELSE ")
	w.Value(cs.els, opts)
	w.Write(" END")

	return &Column{
		Writer: w,
		Type:   customCol,
		Value:  w.ToString(),
	}
}

func (c *ColExpr) Condition(identifier any) core.Cond {
	cond := &condition.Condition{}
	return cond.Start(identifier)
}

func (c *ColExpr) Search(fn func(cs core.CaseSearch)) core.Column {
	w := c.Statement.Writer()
	opts := &core.ValueOpts{Definition: true}

	cs := &Case[core.CondNext]{}
	fn(cs)

	w.Write("CASE WHEN")
	for _, arg := range cs.args {
		w.Char(' ')
		arg.when.(*condition.ConditionNext).Build(w, opts)
		w.Write(" THEN ")
		w.Value(arg.do, opts)
		w.Char(' ')
	}
	w.Write("ELSE ")
	w.Value(cs.els, opts)
	w.Write(" END")

	return &Column{
		Writer: w,
		Type:   customCol,
		Value:  w.ToString(),
	}
}
