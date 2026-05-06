package entity

type (
	Parameter struct{ Value any }
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

// Kinds
func (*Parameter) Kind() ValueKind  { return ValueParameter }
func (*Condition) Kind() ValueKind  { return ValueCondition }
func (*Concat) Kind() ValueKind     { return ValueConcat }
func (*CaseSwitch) Kind() ValueKind { return ValueCaseSwitch }
func (*CaseSearch) Kind() ValueKind { return ValueCaseSearch }

// Writers
func (v *Parameter) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WritePlaceholder(w)
}
func (v *Condition) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WriteCondition(w, v, mode)
}
func (v *Concat) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WriteConcat(w, v)
}
func (v *CaseSwitch) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WriteCaseSwitch(w, v, mode)
}
func (v *CaseSearch) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WriteCaseSearch(w, v, mode)
}

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
