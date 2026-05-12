package dialect

import (
	"github.com/laacin/inyorm/internal/entity/core"
	"github.com/laacin/inyorm/internal/entity/dml"
)

func (dial *StdDialect) WritePlaceholder(w core.Writer, count int) {
	w.Char('?')
}

func (dial *StdDialect) WriteCondition(w core.Writer, cond *dml.Condition, mode core.WritingMode) {
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
		case dml.PredIsNull:

		case dml.PredBetween:
			w.Char(' ')
			w.Value(pred.Values[0], mode)
			w.Write(" AND ")
			w.Value(pred.Values[1], mode)

		case dml.PredIn:
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

func (dial *StdDialect) WriteConcat(w core.Writer, con *dml.Concat) {
	w.Write("CONCAT")
	w.Char('(')
	for i, val := range con.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, core.WriteExpr)
	}
	w.Char(')')
}

func (dial *StdDialect) WriteCaseSwitch(w core.Writer, cas *dml.CaseSwitch, mode core.WritingMode) {
	w.Write("CASE")
	w.Char(' ')
	w.Value(cas.Cond, core.WriteExpr)

	for _, when := range cas.Whens {
		w.Write(" WHEN ")
		w.Value(when.When, core.WriteExpr)

		w.Write(" THEN ")
		w.Value(when.Then, core.WriteExpr)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, core.WriteExpr)
		w.Char(' ')
	}

	w.Write("END")
}

func (dial *StdDialect) WriteCaseSearch(w core.Writer, cas *dml.CaseSearch, mode core.WritingMode) {
	w.Write("CASE WHEN")
	w.Char(' ')

	for _, when := range cas.Whens {
		w.Value(when.When, core.WriteExpr)
		w.Write(" THEN ")
		w.Value(when.Then, core.WriteExpr)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, core.WriteExpr)
		w.Char(' ')
	}

	w.Write("END")
}

// --- Helpers

func getOp(op dml.PredOperator, negated bool) string {
	if negated {
		return negatedMap[op]
	}
	return operatorsMap[op]
}

var operatorsMap = map[dml.PredOperator]string{
	dml.PredEqual:   "=",
	dml.PredLike:    "LIKE",
	dml.PredIn:      "IN",
	dml.PredBetween: "BETWEEN",
	dml.PredGreater: ">",
	dml.PredLess:    "<",
	dml.PredIsNull:  "IS NULL",
}

var negatedMap = map[dml.PredOperator]string{
	dml.PredEqual:   "<>",
	dml.PredLike:    "NOT LIKE",
	dml.PredIn:      "NOT IN",
	dml.PredBetween: "NOT BETWEEN",
	dml.PredGreater: ">=",
	dml.PredLess:    "<=",
	dml.PredIsNull:  "IS NOT NULL",
}

var connectorMap = map[dml.PredConnector]string{
	dml.PredAnd: "AND",
	dml.PredOr:  "OR",
}
