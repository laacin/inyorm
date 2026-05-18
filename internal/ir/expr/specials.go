package expr

type (
	Parameter struct {
		Value any
		Store bool
	}
	Condition struct {
		Predicates []Predicate
		Connectors []PredConnector
	}
	Concat     struct{ Values []any }
	CaseSwitch struct {
		Cond  any
		Whens []CaseWhen
		Els   any
	}
	CaseSearch struct {
		Whens []CaseWhen
		Els   any
	}
)

// --- Condition utilities
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

// --- Case utilities
type CaseWhen struct {
	When any
	Then any
}
