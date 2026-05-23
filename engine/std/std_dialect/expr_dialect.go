package std_dialect

import (
	"strconv"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// ----- LITERALS -----

func (*Dialect) WriteString(w core.Writer, v string) {
	w.Char(quote)
	w.Write(v)
	w.Char(quote)
}

func (*Dialect) WriteNumber(w core.Writer, v int) {
	r := strconv.Itoa(v)
	w.Write(r)
}

func (*Dialect) WriteFloat(w core.Writer, v float64) {
	r := strconv.FormatFloat(float64(v), 'f', -1, 32)
	w.Write(r)
}

func (*Dialect) WriteBool(w core.Writer, v bool) {
	if v {
		w.Char('1')
		return
	}
	w.Char('0')
}

func (*Dialect) WriteNull(w core.Writer) {
	w.Write("NULL")
}

// ---- Placeholder -----

func (*Dialect) WritePlaceholder(w core.Writer) {
	w.Char('?')
}

// ---- Cond ----

func (*Dialect) WriteCond(w core.Writer, cond *expr.Cond, mode core.WritingMode) {
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

func (*Dialect) WriteConcat(w core.Writer, con *expr.Concat, mode core.WritingMode) {
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

func (*Dialect) WriteCaseSwitch(w core.Writer, cas *expr.Case, mode core.WritingMode) {
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

func (*Dialect) WriteCaseSearch(w core.Writer, cas *expr.Case, mode core.WritingMode) {
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

func (*Dialect) WriteTable(w core.Writer, tbl *expr.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

// ---- COLUMN ----

func (*Dialect) WriteColBase(w core.Writer, col *expr.Col) {
	if ref, ok := w.GetRef(col.Ref); ok {
		w.Char(ref)
		w.Char('.')
	}
	w.Write(col.Name)
}

func (s *Dialect) WriteColExpr(w core.Writer, col *expr.Col) {
	if col.Value == "" {
		s.Self.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (s *Dialect) WriteColAlias(w core.Writer, col *expr.Col) {
	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		s.Self.WriteColExpr(w, col)
		return
	}

	s.Self.WriteColBase(w, col)
}

func (s *Dialect) WriteColDef(w core.Writer, col *expr.Col) {
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

// ---- COLUMN EXPRS ----

func (*Dialect) WriteColArith(w core.Writer, arith *expr.ColArith) {
	w.Char(' ')
	w.Char(arithMap[arith.Kind])
	w.Char(' ')
	w.Value(arith.Value, core.WriteExpr)
}

func (*Dialect) WriteColScalar(w core.Writer, scalar *expr.ColScalar) {
	w.Wrap(func(current string, w core.Writer) {
		w.Write(scalarMap[scalar.Kind])
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*Dialect) WriteColWrap(w core.Writer) {
	w.Wrap(func(current string, w core.Writer) {
		w.Char('(')
		w.Write(current)
		w.Char(')')
	})
}

func (*Dialect) WriteColAggr(w core.Writer, aggr *expr.ColAggr) {
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
