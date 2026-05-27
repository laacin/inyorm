package expr

type (
	PredOperator  int
	PredConnector int
)

const (
	PredEqual PredOperator = iota
	PredLike
	PredIn
	PredBetween
	PredGreater
	PredLess
	PredIsNull

	PredAnd PredConnector = iota
	PredOr
)

type Predicate struct {
	Negated    bool
	Identifier any
	Operator   PredOperator
	Values     []any
	Closed     bool // Guard for incomplete expressions
}

func (op PredOperator) String() string {
	switch op {
	case PredEqual:
		return "EQUAL"
	case PredLike:
		return "LIKE"
	case PredIn:
		return "IN"
	case PredBetween:
		return "BETWEEN"
	case PredGreater:
		return "GREATER"
	case PredLess:
		return "LESS"
	case PredIsNull:
		return "IS NULL"
	default:
		return "UNKNOWN"
	}
}

func (con PredConnector) String() string {
	switch con {
	case PredAnd:
		return "AND"
	case PredOr:
		return "OR"
	default:
		return "UNKNOWN"
	}
}
