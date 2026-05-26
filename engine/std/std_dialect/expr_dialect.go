package std_dialect

import (
	"strconv"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// ----- LITERALS -----

func (*Dialect) WriteLitString(w core.Writer, v string) {
	w.Char(quote)
	w.Write(v)
	w.Char(quote)
}

func (*Dialect) WriteLitInt(w core.Writer, v int) {
	r := strconv.Itoa(v)
	w.Write(r)
}

func (*Dialect) WriteLitFloat(w core.Writer, v float64) {
	r := strconv.FormatFloat(float64(v), 'f', -1, 32)
	w.Write(r)
}

func (*Dialect) WriteLitBool(w core.Writer, v bool) {
	if v {
		w.Char('1')
		return
	}
	w.Char('0')
}

func (*Dialect) WriteLitNull(w core.Writer) {
	w.Write("NULL")
}

// ---- Placeholder -----

func (*Dialect) WriteExprPlaceholder(w core.Writer, p *expr.Placeholder) {
	w.Char('?')
}

// ---- Cond ----

func (*Dialect) WriteExprCond(w core.Writer, cond *expr.Cond, mode core.WritingMode) {
	w.Char('(')
	for i, pred := range cond.Predicates {
		if !pred.Closed {
			// TODO: avoid only write '()'
			continue
		}

		if i > 0 {
			w.Char(' ')
			w.Write(connectorMap[cond.Connectors[i-1]]) // TODO: could be fragile
			w.Char(' ')
		}

		w.Value(pred.Identifier, mode)
		w.Char(' ')
		w.Write(getOp(pred.Operator, pred.Negated))
		switch pred.Operator {
		case expr.PredIsNull:

		case expr.PredBetween:
			w.Char(' ')
			w.Value(pred.Values[0], mode)
			w.Write(" AND ")
			w.Value(pred.Values[1], mode)

		case expr.PredIn:
			w.Char(' ')

			w.Char('(')
			for i, v := range pred.Values {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, mode)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(pred.Values[0], mode)
		}
	}
	w.Char(')')
}

// ----- CONCAT -----

func (*Dialect) WriteExprConcat(w core.Writer, con *expr.Concat, mode core.WritingMode) {
	w.Write("CONCAT")
	w.Char('(')
	for i, val := range con.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, mode)
	}
	w.Char(')')
}

// ----- CASE -----

func (*Dialect) WriteExprCaseSwitch(w core.Writer, cas *expr.Case, mode core.WritingMode) {
	w.Write("CASE")
	w.Char(' ')
	w.Value(cas.Cond, core.WriteExpr)

	for _, when := range cas.Whens {
		w.Write(" WHEN ")
		w.Value(when.When, mode)

		w.Write(" THEN ")
		w.Value(when.Then, mode)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, mode)
		w.Char(' ')
	}

	w.Write("END")
}

func (*Dialect) WriteExprCaseSearch(w core.Writer, cas *expr.Case, mode core.WritingMode) {
	w.Write("CASE WHEN")
	w.Char(' ')

	for _, when := range cas.Whens {
		w.Value(when.When, mode)
		w.Write(" THEN ")
		w.Value(when.Then, mode)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, mode)
		w.Char(' ')
	}

	w.Write("END")
}

// ---- TABLE ----

func (*Dialect) WriteExprTable(w core.Writer, tbl *expr.Table) {
	w.Write(tbl.Name)
	if ref := tbl.Ref(); ref.Enable {
		w.Char(' ')
		w.Char(ref.Ref)
	}
}

// ---- COLUMN ----

func (*Dialect) WriteExprColBase(w core.Writer, col *expr.Col) {
	if ref := col.Ref(); ref.Enable {
		w.Char(ref.Ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (s *Dialect) WriteExprColExpr(w core.Writer, col *expr.Col) {
	if col.Value == "" {
		s.Self.WriteExprColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (s *Dialect) WriteExprColAlias(w core.Writer, col *expr.Col) {
	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		s.Self.WriteExprColExpr(w, col)
		return
	}

	s.Self.WriteExprColBase(w, col)
}

func (s *Dialect) WriteExprColDef(w core.Writer, col *expr.Col) {
	if col.Value == "" {
		s.WriteExprColBase(w, col)
		return
	}

	s.WriteExprColExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// ---- COLUMN EXPRS ----

func (*Dialect) WriteExprColArith(w core.Writer, arith *expr.ColArith) {
	w.Char(' ')
	w.Char(arithMap[arith.Kind])
	w.Char(' ')
	w.Value(arith.Value, core.WriteExpr)
}

func (*Dialect) WriteExprColScalar(w core.Writer, scalar *expr.ColScalar) {
	w.Wrap(func(current string, w core.Writer) {
		w.Write(scalarMap[scalar.Kind])
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*Dialect) WriteExprColWrap(w core.Writer) {
	w.Wrap(func(current string, w core.Writer) {
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*Dialect) WriteExprColAggr(w core.Writer, aggr *expr.ColAggr) {
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

// --- Helpers
var quote byte = "'"[0]

func getOp(op expr.PredOperator, negated bool) string {
	if negated {
		return negatedMap[op]
	}
	return operatorsMap[op]
}

var operatorsMap = map[expr.PredOperator]string{
	expr.PredEqual:   "=",
	expr.PredLike:    "LIKE",
	expr.PredIn:      "IN",
	expr.PredBetween: "BETWEEN",
	expr.PredGreater: ">",
	expr.PredLess:    "<",
	expr.PredIsNull:  "IS NULL",
}

var negatedMap = map[expr.PredOperator]string{
	expr.PredEqual:   "<>",
	expr.PredLike:    "NOT LIKE",
	expr.PredIn:      "NOT IN",
	expr.PredBetween: "NOT BETWEEN",
	expr.PredGreater: ">=",
	expr.PredLess:    "<=",
	expr.PredIsNull:  "IS NOT NULL",
}

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

var connectorMap = map[expr.PredConnector]string{
	expr.PredAnd: "AND",
	expr.PredOr:  "OR",
}
