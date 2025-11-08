package column

import "github.com/laacin/inyorm/internal/core"

func wDefinition(w core.Writer, col *Column) {
	switch col.Type {
	case normalCol:
		if col.Alias != "" {
			w.ColRef(col.Alias)
			w.Char('.')
		}
		w.Write(col.Value)

	case customCol:
		w.Write(col.Value)
		if col.Alias != "" {
			w.Write(" AS ")
			w.Write(col.Alias)
		}
	}
}

func wReference(w core.Writer, col *Column) {
	switch col.Type {
	case normalCol:
		if col.Alias != "" {
			w.ColRef(col.Alias)
			w.Char('.')
		}
		w.Write(col.Value)

	case customCol:
		if col.Alias == "" {
			w.Write(col.Value)
			return
		}
		w.Write(col.Alias)
	}
}

func wPrev(w core.Writer, col *Column) {
	if col.Type == normalCol {
		if col.Alias != "" {
			w.ColRef(col.Alias)
			w.Char('.')
			col.Alias = ""
		}
		col.Type = customCol
	}
	w.Write(col.Value)
}

func wOp(col *Column, arg byte, value any) {
	w := col.Writer
	w.Reset()

	wPrev(w, col)
	w.Char(' ')
	w.Char(arg)
	w.Char(' ')
	w.Value(value, &core.ValueOpts{Definition: true})

	col.Value = w.ToString()
}

func wFunc(col *Column, arg string) {
	w := col.Writer
	w.Reset()

	w.Write(arg)
	w.Char('(')
	wPrev(w, col)
	w.Char(')')

	col.Value = w.ToString()
}

func wWrap(col *Column) {
	w := col.Writer
	w.Reset()

	w.Char('(')
	wPrev(w, col)
	w.Char(')')

	col.Value = w.ToString()
}

func wAggr(col *Column, distinct bool, aggr string) {
	w := col.Writer
	w.Reset()

	w.Write(aggr)
	w.Char('(')
	if distinct {
		w.Write("DISTINCT ")
	}
	wPrev(w, col)
	w.Char(')')

	col.Value = w.ToString()
}

func wAs(col *Column, value string) {
	col.Type = customCol
	col.Alias = value
}
