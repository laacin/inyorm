package std_expr

import (
	"github.com/laacin/inyorm/internal/entity/core"
	"github.com/laacin/inyorm/internal/entity/expr"
)

func (s *ExprSyntax) WriteTable(w core.Writer, tbl *expr.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

func (s *ExprSyntax) WriteColBase(w core.Writer, col *expr.Column) {
	if ref, ok := w.GetRef(col.Ref); ok {
		w.Char(ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (s *ExprSyntax) WriteColExpr(w core.Writer, col *expr.Column) {
	s.buildCol(w.New(), col)

	if col.Value == "" {
		s.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (s *ExprSyntax) WriteColAlias(w core.Writer, col *expr.Column) {
	s.buildCol(w.New(), col)

	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		s.WriteColExpr(w, col)
		return
	}

	s.WriteColBase(w, col)
}

func (s *ExprSyntax) WriteColDef(w core.Writer, col *expr.Column) {
	s.buildCol(w.New(), col)

	if col.Value == "" {
		s.WriteColBase(w, col)
		return
	}

	s.WriteColExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// --- Helpers
func (s *ExprSyntax) buildFirst(w core.Writer, col *expr.Column) {
	if col.From != nil {
		if col.From.Kind() == expr.ValueWildcard {
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

	s.WriteColBase(w, col)
}

// FIX: illegible
func (s *ExprSyntax) buildCol(w core.Writer, col *expr.Column) {
	if (col == nil) || (col.Exprs == nil && col.Aggr == nil && col.From == nil) {
		return
	}

	s.buildFirst(w, col)
	if col.Exprs != nil {
		for _, e := range col.Exprs {
			if scalar, ok := scalarMap[e.Kind]; ok {
				wScalar(w, scalar)
				continue
			}

			if arith, ok := arithMap[e.Kind]; ok {
				wArith(w, arith, e.Value)
				continue
			}

			if e.Kind == expr.ColArithWrap {
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
var aggrMap = map[expr.ColKindExpr]string{
	expr.ColAggrCount: "COUNT",
	expr.ColAggrSum:   "SUM",
	expr.ColAggrMin:   "MIN",
	expr.ColAggrMax:   "MAX",
	expr.ColAggrAvg:   "AVG",
}

var scalarMap = map[expr.ColKindExpr]string{
	expr.ColScalarLower: "LOWER",
	expr.ColScalarUpper: "UPPER",
	expr.ColScalarTrim:  "TRIM",
	expr.ColScalarRound: "ROUND",
	expr.ColScalarAbs:   "ABS",
}

var arithMap = map[expr.ColKindExpr]byte{
	expr.ColArithAdd: '+',
	expr.ColArithSub: '-',
	expr.ColArithMul: '*',
	expr.ColArithDiv: '/',
	expr.ColArithMod: '%',
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
