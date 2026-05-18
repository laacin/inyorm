package std_expr

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

func (*ExprSyntax) WritePlaceholder(w core.Writer) {
	w.Char('?')
}

func (*ExprSyntax) WriteCondition(w core.Writer, cond *expr.Condition, mode core.WritingMode) {
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

func (*ExprSyntax) WriteConcat(w core.Writer, con *expr.Concat, mode core.WritingMode) {
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

func (*ExprSyntax) WriteCaseSwitch(w core.Writer, cas *expr.CaseSwitch, mode core.WritingMode) {
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

func (*ExprSyntax) WriteCaseSearch(w core.Writer, cas *expr.CaseSearch, mode core.WritingMode) {
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

// --- Helpers

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

var connectorMap = map[expr.PredConnector]string{
	expr.PredAnd: "AND",
	expr.PredOr:  "OR",
}
