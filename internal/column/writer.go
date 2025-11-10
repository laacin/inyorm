package column

import "github.com/laacin/inyorm/internal/core"

func wBase(w core.Writer, col *Column) {
	if col.table != "" {
		w.ColRef(col.table)
		w.Char('.')
	}
	w.Write(col.base)
}

func wExpr(w core.Writer, col *Column) {
	if !col.custom {
		wBase(w, col)
		return
	}
	w.Write(col.expr)
}

func wAlias(w core.Writer, col *Column) {
	if !col.custom {
		wBase(w, col)
		return
	}
	wExpr(w, col)
}

func wDef(w core.Writer, col *Column) {
	if !col.custom {
		wBase(w, col)
		return
	}
	w.Write(col.expr)
	if col.alias != "" {
		w.Write(" AS ")
		w.Write(col.alias)
	}
}

// -- internal writers

func wPrev(w core.Writer, col *Column) {
	if !col.custom {
		wBase(w, col)
		col.custom = true
	}
	w.Write(col.expr)
}

func wOp(col *Column, arg byte, value any) {
	w := col.writer
	w.Reset()

	wPrev(w, col)
	w.Char(' ')
	w.Char(arg)
	w.Char(' ')
	w.Value(value, core.WriterOpts{ColType: core.ColTypExpr})

	col.expr = w.ToString()
}

func wFunc(col *Column, arg string) {
	w := col.writer
	w.Reset()

	w.Write(arg)
	w.Char('(')
	wPrev(w, col)
	w.Char(')')

	col.expr = w.ToString()
}

func wWrap(col *Column) {
	w := col.writer
	w.Reset()

	w.Char('(')
	wPrev(w, col)
	w.Char(')')

	col.expr = w.ToString()
}

func wAggr(col *Column, distinct bool, aggr string) {
	w := col.writer
	w.Reset()

	w.Write(aggr)
	w.Char('(')
	if distinct {
		w.Write("DISTINCT ")
	}
	wPrev(w, col)
	w.Char(')')

	col.expr = w.ToString()
}

func wAs(col *Column, value string) {
	col.custom = true
	col.alias = value
}
