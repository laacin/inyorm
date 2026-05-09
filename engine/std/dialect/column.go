package dialect

import "github.com/laacin/inyorm/internal/entity"

func (dial *StdDialect) WriteTable(w entity.Writer, tbl *entity.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

func (dial *StdDialect) WriteColBase(w entity.Writer, col *entity.Column) {
	if ref, ok := w.GetRef(col.Ref); ok {
		w.Char(ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (dial *StdDialect) WriteColExpr(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Value == "" {
		dial.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (dial *StdDialect) WriteColAlias(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		dial.WriteColExpr(w, col)
		return
	}

	dial.WriteColBase(w, col)
}

func (dial *StdDialect) WriteColDef(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Value == "" {
		dial.WriteColBase(w, col)
		return
	}

	dial.WriteColExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// --- Helpers
func (dial *StdDialect) BuildFirst(w entity.Writer, col *entity.Column) {
	if col.From != nil {
		if col.From.Kind() == entity.ValueWildcard {
			if ref, ok := w.GetRef(col.Ref); ok {
				w.Char(ref)
				w.Char('.')
			}
		}
		w.Value(col.From, entity.WriteExpr)
		return
	}

	if col.Value != "" {
		w.Write(col.Value)
		return
	}

	dial.WriteColBase(w, col)
}

// FIX: illegible
func (dial *StdDialect) BuildCol(w entity.Writer, col *entity.Column) {
	if (col == nil) || (col.Exprs == nil && col.Aggr == nil && col.From == nil) {
		return
	}

	dial.BuildFirst(w, col)
	if col.Exprs != nil {
		for _, expr := range col.Exprs {
			if scalar, ok := scalarMap[expr.Kind]; ok {
				wScalar(w, scalar)
				continue
			}

			if arith, ok := arithMap[expr.Kind]; ok {
				wArith(w, arith, expr.Value)
				continue
			}

			if expr.Kind == entity.ColArithWrap {
				wWrap(w)
				continue
			}
		}
		col.Exprs = nil
	}

	if col.Aggr != nil {
		if aggr, ok := aggrMap[col.Aggr.Kind]; ok {
			wAggr(w, col.Aggr.Value, aggr)
		}
		col.Aggr = nil
	}

	col.Value = w.ToString()
}

// maps
var aggrMap = map[entity.ColKindExpr]string{
	entity.ColAggrCount: "COUNT",
	entity.ColAggrSum:   "SUM",
	entity.ColAggrMin:   "MIN",
	entity.ColAggrMax:   "MAX",
	entity.ColAggrAvg:   "AVG",
}

var scalarMap = map[entity.ColKindExpr]string{
	entity.ColScalarLower: "LOWER",
	entity.ColScalarUpper: "UPPER",
	entity.ColScalarTrim:  "TRIM",
	entity.ColScalarRound: "ROUND",
	entity.ColScalarAbs:   "ABS",
}

var arithMap = map[entity.ColKindExpr]byte{
	entity.ColArithAdd: '+',
	entity.ColArithSub: '-',
	entity.ColArithMul: '*',
	entity.ColArithDiv: '/',
	entity.ColArithMod: '%',
}

func wArith(w entity.Writer, arg byte, value any) {
	w.Char(' ')
	w.Char(arg)
	w.Char(' ')
	w.Value(value, entity.WriteExpr)
}

func wScalar(w entity.Writer, arg string) {
	prev := w.ToString()
	w.Reset()

	w.Write(arg)
	w.Char('(')
	w.Write(prev)
	w.Char(')')
}

func wWrap(w entity.Writer) {
	prev := w.ToString()
	w.Reset()

	w.Char('(')
	w.Write(prev)
	w.Char(')')
}

func wAggr(w entity.Writer, distinct any, aggr string) {
	prev := w.ToString()
	w.Reset()

	w.Write(aggr)
	w.Char('(')
	if dist, ok := distinct.(bool); ok && dist {
		w.Write("DISTINCT")
		w.Char(' ')
	}
	w.Write(prev)
	w.Char(')')
}
