package std_expr

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
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
	if col.Value == "" {
		s.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (s *ExprSyntax) WriteColAlias(w core.Writer, col *expr.Column) {
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
// func (s *ExprSyntax) buildFirst(w core.Writer, col *expr.Column) {
// 	if col.From != nil {
// 		if col.From.Kind() == expr.ValueWildcard {
// 			if ref, ok := w.GetRef(col.Ref); ok {
// 				w.Char(ref)
// 				w.Char('.')
// 			}
// 		}
// 		w.Value(col.From, core.WriteExpr)
// 		return
// 	}
//
// 	if col.Value != "" {
// 		w.Write(col.Value)
// 		return
// 	}
//
// 	s.WriteColBase(w, col)
// }
//
// // FIX: illegible
// func (s *ExprSyntax) buildCol(w core.Writer, col *expr.Column) {
// 	if (col == nil) || (col.Exprs == nil && col.Aggr == nil && col.From == nil) {
// 		return
// 	}
//
// 	s.buildFirst(w, col)
// 	if col.Exprs != nil {
// 		for _, e := range col.Exprs {
// 			if scalar, ok := scalarMap[e.Kind]; ok {
// 				wScalar(w, scalar)
// 				continue
// 			}
//
// 			if arith, ok := arithMap[e.Kind]; ok {
// 				wArith(w, arith, e.Value)
// 				continue
// 			}
//
// 			if e.Kind == expr.ColArithWrap {
// 				wWrap(w)
// 				continue
// 			}
// 		}
// 		col.Exprs = nil
// 	}
//
// 	if col.Aggr != nil {
// 		if aggr, ok := aggrMap[col.Aggr.Kind]; ok {
// 			wAggr(w, col.Aggr.Value, aggr)
// 		}
// 		col.Aggr = nil
// 	}
//
// 	col.Value = w.ToString()
// }

// maps
var aggrMap = map[expr.ColAggrKind]string{
	expr.ColAggrCount: "COUNT",
	expr.ColAggrSum:   "SUM",
	expr.ColAggrMin:   "MIN",
	expr.ColAggrMax:   "MAX",
	expr.ColAggrAvg:   "AVG",
}

var scalarMap = map[expr.ColScalarKind]string{
	expr.ColScalarLower: "LOWER",
	expr.ColScalarUpper: "UPPER",
	expr.ColScalarTrim:  "TRIM",
	expr.ColScalarRound: "ROUND",
	expr.ColScalarAbs:   "ABS",
}

var arithMap = map[expr.ColArithKind]byte{
	expr.ColArithAdd: '+',
	expr.ColArithSub: '-',
	expr.ColArithMul: '*',
	expr.ColArithDiv: '/',
	expr.ColArithMod: '%',
}

func (*ExprSyntax) WriteColArith(w core.Writer, arith *expr.ColArith) {
	w.Char(' ')
	w.Char(arithMap[arith.Kind])
	w.Char(' ')
	w.Value(arith.Value, core.WriteExpr)
}

func (*ExprSyntax) WriteColScalar(w core.Writer, scalar *expr.ColScalar) {
	w.Wrap(func(current string, w core.Writer) {
		w.Write(scalarMap[scalar.Kind])
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*ExprSyntax) WriteColWrap(w core.Writer) {
	w.Wrap(func(current string, w core.Writer) {
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*ExprSyntax) WriteColAggr(w core.Writer, aggr *expr.ColAggr) {
	w.Wrap(func(current string, w core.Writer) {
		w.Write(aggrMap[aggr.Kind])
		w.Char('(')
		if aggr.Distinct {
			w.Write("DISTINCT")
			w.Char(' ')
		}
		w.Write(current)
		w.Char(')')
	})
}
