package std_expr

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

func (s *ExprSyntax) WriteTable(w core.Writer, tbl *expr.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

func (s *ExprSyntax) WriteColBase(w core.Writer, col *expr.Col) {
	if ref, ok := w.GetRef(col.Ref); ok {
		w.Char(ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (s *ExprSyntax) WriteColExpr(w core.Writer, col *expr.Col) {
	if col.Value == "" {
		s.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (s *ExprSyntax) WriteColAlias(w core.Writer, col *expr.Col) {
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

func (s *ExprSyntax) WriteColDef(w core.Writer, col *expr.Col) {
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
