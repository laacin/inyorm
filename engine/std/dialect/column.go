package dialect

import (
	"github.com/laacin/inyorm/internal/entity/core"
	"github.com/laacin/inyorm/internal/entity/dml"
)

func (dial *StdDialect) WriteTable(w core.Writer, tbl *dml.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

func (dial *StdDialect) WriteColBase(w core.Writer, col *dml.Column) {
	if ref, ok := w.GetRef(col.Ref); ok {
		w.Char(ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (dial *StdDialect) WriteColExpr(w core.Writer, col *dml.Column) {
	dial.BuildCol(w.New(), col)

	if col.Value == "" {
		dial.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (dial *StdDialect) WriteColAlias(w core.Writer, col *dml.Column) {
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

func (dial *StdDialect) WriteColDef(w core.Writer, col *dml.Column) {
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
func (dial *StdDialect) BuildFirst(w core.Writer, col *dml.Column) {
	if col.From != nil {
		if col.From.Kind() == dml.ValueWildcard {
			if ref, ok := w.GetRef(col.Ref); ok {
				w.Char(ref)
				w.Char('.')
			}
		}
		w.Value(col.From, core.WriteExpr)
		return
	}

	if col.Value != "" {
		w.Write(col.Value)
		return
	}

	dial.WriteColBase(w, col)
}

// FIX: illegible
func (dial *StdDialect) BuildCol(w core.Writer, col *dml.Column) {
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

			if expr.Kind == dml.ColArithWrap {
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
var aggrMap = map[dml.ColKindExpr]string{
	dml.ColAggrCount: "COUNT",
	dml.ColAggrSum:   "SUM",
	dml.ColAggrMin:   "MIN",
	dml.ColAggrMax:   "MAX",
	dml.ColAggrAvg:   "AVG",
}

var scalarMap = map[dml.ColKindExpr]string{
	dml.ColScalarLower: "LOWER",
	dml.ColScalarUpper: "UPPER",
	dml.ColScalarTrim:  "TRIM",
	dml.ColScalarRound: "ROUND",
	dml.ColScalarAbs:   "ABS",
}

var arithMap = map[dml.ColKindExpr]byte{
	dml.ColArithAdd: '+',
	dml.ColArithSub: '-',
	dml.ColArithMul: '*',
	dml.ColArithDiv: '/',
	dml.ColArithMod: '%',
}

func wArith(w core.Writer, arg byte, value any) {
	w.Char(' ')
	w.Char(arg)
	w.Char(' ')
	w.Value(value, core.WriteExpr)
}

func wScalar(w core.Writer, arg string) {
	prev := w.ToString()
	w.Reset()

	w.Write(arg)
	w.Char('(')
	w.Write(prev)
	w.Char(')')
}

func wWrap(w core.Writer) {
	prev := w.ToString()
	w.Reset()

	w.Char('(')
	w.Write(prev)
	w.Char(')')
}

func wAggr(w core.Writer, distinct any, aggr string) {
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
