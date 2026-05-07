package standard

import "github.com/laacin/inyorm/internal/entity"

func (dial *DialectStandard) WritePlaceholder(w entity.Writer) {
	w.Char('?')
}

func (dial *DialectStandard) WriteCondition(w entity.Writer, cond *entity.Condition, mode entity.WritingMode) {
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
		case entity.PredIsNull:

		case entity.PredBetween:
			w.Char(' ')
			w.Value(pred.Values[0], mode)
			w.Write(" AND ")
			w.Value(pred.Values[1], mode)

		case entity.PredIn:
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

func (dial *DialectStandard) WriteConcat(w entity.Writer, con *entity.Concat) {
	w.Write("CONCAT")
	w.Char('(')
	for i, val := range con.Values {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(val, entity.WriteExpr)
	}
	w.Char(')')
}

func (dial *DialectStandard) WriteCaseSwitch(w entity.Writer, cas *entity.CaseSwitch, mode entity.WritingMode) {
	w.Write("CASE")
	w.Char(' ')
	w.Value(cas.Cond, entity.WriteExpr)

	for _, when := range cas.Whens {
		w.Write(" WHEN ")
		w.Value(when.When, entity.WriteExpr)

		w.Write(" THEN ")
		w.Value(when.Then, entity.WriteExpr)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, entity.WriteExpr)
		w.Char(' ')
	}

	w.Write("END")
}

func (dial *DialectStandard) WriteCaseSearch(w entity.Writer, cas *entity.CaseSearch, mode entity.WritingMode) {
	w.Write("CASE WHEN")
	w.Char(' ')

	for _, when := range cas.Whens {
		w.Value(when.When, entity.WriteExpr)
		w.Write(" THEN ")
		w.Value(when.Then, entity.WriteExpr)
		w.Char(' ')
	}

	if cas.Els != nil {
		w.Write("ELSE")
		w.Char(' ')
		w.Value(cas.Els, entity.WriteExpr)
		w.Char(' ')
	}

	w.Write("END")
}

// --- Helpers

func getOp(op entity.PredOperator, negated bool) string {
	if negated {
		return negatedMap[op]
	}
	return operatorsMap[op]
}

var operatorsMap = map[entity.PredOperator]string{
	entity.PredEqual:   "=",
	entity.PredLike:    "LIKE",
	entity.PredIn:      "IN",
	entity.PredBetween: "BETWEEN",
	entity.PredGreater: ">",
	entity.PredLess:    "<",
	entity.PredIsNull:  "IS NULL",
}

var negatedMap = map[entity.PredOperator]string{
	entity.PredEqual:   "<>",
	entity.PredLike:    "NOT LIKE",
	entity.PredIn:      "NOT IN",
	entity.PredBetween: "NOT BETWEEN",
	entity.PredGreater: ">=",
	entity.PredLess:    "<=",
	entity.PredIsNull:  "IS NOT NULL",
}

var connectorMap = map[entity.PredConnector]string{
	entity.PredAnd: "AND",
	entity.PredOr:  "OR",
}
