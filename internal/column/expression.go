package column

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type ColExpr struct{}

func (c *ColExpr) Col(col, table string) core.Column {
	return &Column{table: table, base: col}
}

func (c *ColExpr) All() core.Column {
	return &Column{base: "*"}
}

func (c *ColExpr) Concat(v []any) core.Column {
	expr := func(w core.Writer) {
		w.Write("CONCAT(")
		for i, val := range v {
			if i > 0 {
				w.Write(", ")
			}
			w.Value(val, core.ColumnIdentWriteOpt)
		}
		w.Char(')')
	}

	col := &Column{}
	col.builder.wExpr(expr)
	return col
}

func (c *ColExpr) Switch(cond any, cs core.Case) core.Column {
	val := cs.(*Case)

	expr := func(w core.Writer) {
		w.Write("CASE ")
		w.Value(cond, core.ColumnIdentWriteOpt)
		for _, arg := range val.exprs {
			w.Write(" WHEN ")
			w.Value(arg.when, core.ColumnIdentWriteOpt)
			w.Write(" THEN ")
			w.Value(arg.do, core.ColumnIdentWriteOpt)
			w.Char(' ')
		}
		if val.els != nil {
			w.Write("ELSE ")
			w.Value(val.els, core.ColumnIdentWriteOpt)
			w.Char(' ')
		}
		w.Write("END")
	}

	col := &Column{}
	col.builder.wExpr(expr)
	return col
}

func (c *ColExpr) Condition(identifier any) core.Condition {
	cond := &condition.Condition{}
	cond.Start(identifier)
	return cond
}

func (c *ColExpr) Search(cs core.Case) core.Column {
	val := cs.(*Case)

	expr := func(w core.Writer) {
		w.Write("CASE WHEN")
		for _, arg := range val.exprs {
			w.Char(' ')
			arg.when.(*condition.Condition).Build(
				w,
				core.ColumnIdentWriteOpt,
				core.ColumnValueWriteOpt,
			)
			w.Write(" THEN ")
			w.Value(arg.do, core.ColumnIdentWriteOpt)
			w.Char(' ')
		}
		w.Write("ELSE ")
		w.Value(val.els, core.ColumnIdentWriteOpt)
		w.Write(" END")
	}

	col := &Column{}
	col.builder.wExpr(expr)
	return col
}
