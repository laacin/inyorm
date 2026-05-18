package expr

type Table struct{ Value string }

type Column struct {
	Name  string // Column base name
	Ref   string // Table reference
	Alias string // Explicit alias
	Value string // Column expression
}

// --- Tools for building
type (
	ColAggrKind   int
	ColScalarKind int
	ColArithKind  int
	ColWrapKind   int
)

const (
	// Aggregation
	ColAggrCount ColAggrKind = iota
	ColAggrSum
	ColAggrMin
	ColAggrMax
	ColAggrAvg

	// Scalar
	ColScalarLower ColScalarKind = iota
	ColScalarUpper
	ColScalarTrim
	ColScalarRound
	ColScalarAbs

	// Arith
	ColArithAdd ColArithKind = iota
	ColArithSub
	ColArithMul
	ColArithDiv
	ColArithMod
)

type ColExpr struct{ Kind any }

type ColAggr struct {
	Kind     ColAggrKind
	Distinct bool
}

type ColScalar struct {
	Kind ColScalarKind
}

type ColArith struct {
	Kind  ColArithKind
	Value any
}

type ColWrap struct{}

func (c *ColExpr) IsScalar() (*ColScalar, bool) {
	v, ok := c.Kind.(*ColScalar)
	return v, ok
}
func (c *ColExpr) IsArith() (*ColArith, bool) {
	v, ok := c.Kind.(*ColArith)
	return v, ok
}
func (c *ColExpr) IsWrap() bool {
	_, ok := c.Kind.(*ColWrap)
	return ok
}
