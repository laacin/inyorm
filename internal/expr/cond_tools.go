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
