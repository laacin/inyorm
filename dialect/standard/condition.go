package standard

import "github.com/laacin/inyorm/intr/dialect"

func (dial *DialectStandard) Cond(w dialect.Writer, cond dialect.Cond, mode dialect.WritingMode) {
	w.Char('(')
	for i, expr := range cond.Exprs {
		if !expr.Closed {
			// TODO: avoid only write '()'
			continue
		}

		if i > 0 {
			w.Char(' ')
			w.Write(connectorMap[cond.Connectors[i-1]]) // TODO: could be fragile
			w.Char(' ')
		}

		w.Value(expr.Identifier, mode)
		w.Char(' ')
		w.Write(getOp(expr.Operator, expr.Negated))

		switch expr.Operator {
		case dialect.ExprIsNull:

		case dialect.ExprBetween:
			w.Char(' ')
			w.Value(expr.Values[0], mode)
			w.Write(" AND ")
			w.Value(expr.Values[1], mode)

		case dialect.ExprIn:
			w.Char(' ')

			w.Char('(')
			for i, v := range expr.Values {
				if i > 0 {
					w.Write(", ")
				}
				w.Value(v, mode)
			}
			w.Char(')')

		default:
			w.Char(' ')
			w.Value(expr.Values[0], mode)
		}
		w.Char(')')
	}
}

// --- Maps

func getOp(op dialect.ExprOperator, negated bool) string {
	if negated {
		return negatedMap[op]
	}
	return operatorsMap[op]
}

var operatorsMap = map[dialect.ExprOperator]string{
	dialect.ExprEqual:   "=",
	dialect.ExprLike:    "LIKE",
	dialect.ExprIn:      "IN",
	dialect.ExprBetween: "BETWEEN",
	dialect.ExprGreater: ">",
	dialect.ExprLess:    "<",
	dialect.ExprIsNull:  "IS NULL",
}

var negatedMap = map[dialect.ExprOperator]string{
	dialect.ExprEqual:   "<>",
	dialect.ExprLike:    "NOT LIKE",
	dialect.ExprIn:      "NOT IN",
	dialect.ExprBetween: "NOT BETWEEN",
	dialect.ExprGreater: ">=",
	dialect.ExprLess:    "<=",
	dialect.ExprIsNull:  "IS NOT NULL",
}

var connectorMap = map[dialect.ExprConnector]string{
	dialect.ExprAnd: "AND",
	dialect.ExprOr:  "OR",
}
